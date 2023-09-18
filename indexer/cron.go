package indexer

import (
	"context"
	"strconv"
	"sync"
	"sync/atomic"

	"github.com/shopspring/decimal"
	rest "github.com/ulerdogan/pickaxe/clients/rest"
	starknet "github.com/ulerdogan/pickaxe/clients/starknet"
	db "github.com/ulerdogan/pickaxe/db/sqlc"
	logger "github.com/ulerdogan/pickaxe/utils/logger"
)

func setupJobs(ix *Indexer) {
	ix.Scheduler.Every(5).Minutes().Do(ix.QueryPrices)
	ix.Scheduler.Every(1).Days().Do(ix.CheckFees)

	go ix.Scheduler.StartBlocking()
}

var feesToCheck = []int{4}

func (ix *Indexer) CheckFees() {
	var pools []db.Pool

	for i := range feesToCheck {
		pls, err := ix.Store.GetPoolsByAmm(context.Background(), int64(feesToCheck[i]))
		if err != nil {
			logger.Error(err, "cannot get the pools by amm: "+strconv.Itoa(feesToCheck[i]))
			return
		}
		pools = append(pools, pls...)
	}

	var scp *atomic.Uint64 = &atomic.Uint64{}
	var wg *sync.WaitGroup = &sync.WaitGroup{}
	for _, pool := range pools {
		wg.Add(1)
		go updateFees(ix.Store, ix.Client, pool, scp, wg)
	}

	wg.Wait()
	logger.Info("in " + strconv.Itoa(len(pools)) + ", fees of " + strconv.Itoa(int(scp.Load())) + " pools are updated")
}

func (ix *Indexer) QueryPrices() {
	tokens, err := ix.Store.GetAllTokensWithTickers(context.Background())
	if err != nil {
		logger.Error(err, "cannot get the token list")
		return
	}

	const numWorkers = 10
	jobs := make(chan db.Token, len(tokens))
	results := make(chan *db.Token, len(tokens))

	for w := 0; w < numWorkers; w++ {
		go func(jobs chan db.Token, results chan *db.Token) {
			getPriceConc(jobs, results, ix.Rest)
		}(jobs, results)
	}

	for _, token := range tokens {
		jobs <- token
	}
	close(jobs)

	var sct atomic.Uintptr
	for res := 0; res < len(tokens); res++ {
		token := <-results
		if token != nil {
			_, err := ix.Store.UpdatePrice(context.Background(), db.UpdatePriceParams{Address: token.Address, Price: token.Price})
			if err != nil {
				logger.Error(err, "cannot update the price of the token: "+token.Name)
				continue
			}
			sct.Add(1)
		}
	}

	logger.Info("in " + strconv.Itoa(len(tokens)) + " tokens, prices of " + strconv.Itoa(int(sct.Load())) + " is synced")

	pools, err := ix.Store.GetAllPools(context.Background())
	if err != nil {
		logger.Error(err, "cannot get the pool list")
		return
	}

	var scp *atomic.Uint64 = &atomic.Uint64{}
	var wg *sync.WaitGroup = &sync.WaitGroup{}
	for _, pool := range pools {
		wg.Add(1)
		go updateValueV2(ix.Store, pool, scp, wg)
	}

	wg.Wait()
	logger.Info("in " + strconv.Itoa(len(tokens)) + " pools, total values of " + strconv.Itoa(int(scp.Load())) + " is synced")
}

func getPriceConc(jobs <-chan db.Token, results chan<- *db.Token, rest rest.Client) {
	api := rest.NewPriceAPI()

	for token := range jobs {
		dc, err := api.GetPrice(rest, token)
		token.Price = dc.String()
		if err != nil {
			results <- nil
		}
		results <- &token
	}
}

// FIXME: only compatible with the v2 pools for now, should be improved in the future
func updateValueV2(store db.Store, pool db.Pool, scp *atomic.Uint64, wg *sync.WaitGroup) error {
	if pool.ReserveA == "0" || pool.ReserveB == "0" {
		return nil
	}

	var priceA, priceB decimal.Decimal

	if pA, err := store.GetTokenAPriceByPool(context.Background(), pool.PoolID); err != nil {
		logger.Error(err, "cannot get the token_a price")
		return err
	} else if pA == "0" {
		return nil
	} else {
		priceA, _ = decimal.NewFromString(pA)
	}

	if pB, err := store.GetTokenBPriceByPool(context.Background(), pool.PoolID); err != nil {
		logger.Error(err, "cannot get the token_b price")
		return err
	} else if pB == "0" {
		return nil
	} else {
		priceB, _ = decimal.NewFromString(pB)
	}

	vlA, _ := decimal.NewFromString(pool.ReserveA)
	vlB, _ := decimal.NewFromString(pool.ReserveA)
	tvl := vlA.Mul(priceA).Add(vlB.Mul(priceB))

	_, err := store.UpdatePoolTV(context.Background(), db.UpdatePoolTVParams{
		PoolID:     pool.PoolID,
		TotalValue: tvl.String(),
	})
	if err != nil {
		logger.Error(err, "cannot get the token_b price")
		return err
	}

	scp.Add(1)
	wg.Done()
	return nil
}

func updateFees(store db.Store, client starknet.Client, pool db.Pool, scp *atomic.Uint64, wg *sync.WaitGroup) error {
	dex, err := client.NewDex(int(pool.AmmID))
	if err != nil {
		logger.Error(err, "cannot get the dex "+strconv.Itoa(int(pool.AmmID))+" to update fees")
		return err
	}

	err = dex.SyncFee(starknet.PoolInfo{
		Address:   pool.Address,
		ExtraData: pool.ExtraData.String,
	}, store, client)
	if err != nil {
		logger.Error(err, "cannot get the fee")
		return err
	}

	scp.Add(1)
	wg.Done()

	return nil
}
