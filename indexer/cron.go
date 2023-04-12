package indexer

import (
	"context"
	"database/sql"
	"fmt"
	"strconv"

	"github.com/shopspring/decimal"
	rest "github.com/ulerdogan/pickaxe/clients/rest"
	db "github.com/ulerdogan/pickaxe/db/sqlc"
	logger "github.com/ulerdogan/pickaxe/utils/logger"
)

func setupJobs(ix *indexer) {
	ix.scheduler.Every(5).Seconds().Do(ix.QueryBlocks)
	ix.scheduler.Every(5).Minutes().Do(ix.QueryPrices)
}

func (ix *indexer) QueryBlocks() {
	if ix.isIndexing {
		return
	}

	ix.ixMutex.Lock()
	defer ix.ixMutex.Unlock()

	ix.isIndexing = true

	lastBlock, err := ix.client.LastBlock()
	if err != nil {
		logger.Error(err, "cannot get the last block")
		ix.isIndexing = false
		return
	}

	if lastBlock > *ix.lastQueried {
		// TODO: do sth
		logger.Info("new block catched: " + fmt.Sprint(lastBlock))
		ix.lastQueried = &lastBlock
		ix.store.UpdateIndexerStatus(context.Background(), sql.NullInt64{Int64: int64(lastBlock), Valid: true})
	} else {
		// FIXME: remove the part
		logger.Info("no new block: " + fmt.Sprint(lastBlock))
	}

	ix.isIndexing = false
}

func (ix *indexer) QueryPrices() {
	tokens, err := ix.store.GetAllTokensWithTickers(context.Background())
	if err != nil {
		logger.Error(err, "cannot get the token list")
		return
	}

	const numWorkers = 10
	jobs := make(chan db.Token, len(tokens))
	results := make(chan *db.Token, len(tokens))

	for w := 0; w < numWorkers; w++ {
		go func(jobs chan db.Token, results chan *db.Token) {
			getPriceConc(jobs, results, ix.rest)
		}(jobs, results)
	}

	for _, token := range tokens {
		jobs <- token
	}
	close(jobs)

	var s int
	for res := 0; res < len(tokens); res++ {
		token := <-results
		if token != nil {
			_, err := ix.store.UpdatePrice(context.Background(), db.UpdatePriceParams{Address: token.Address, Price: token.Price})
			if err != nil {
				logger.Error(err, "cannot update the price of the token: "+token.Name)
				continue
			}
			s++
		}
	}

	logger.Info("in " + strconv.Itoa(len(tokens)) + " tokens, prices of " + strconv.Itoa(s) + " is synced")

	pools, err := ix.store.GetAllPools(context.Background())
	if err != nil {
		logger.Error(err, "cannot get the pool list")
		return
	}

	s = 0
	for _, pool := range pools {
		if pool.ReserveA == "0" || pool.ReserveB == "0" {
			continue
		}

		var priceA, priceB decimal.Decimal

		if pA, err := ix.store.GetTokenAPriceByPool(context.Background(), pool.PoolID); err != nil {
			logger.Error(err, "cannot get the token_a price")
			continue
		} else if pA == "0" {
			continue
		} else {
			priceA, _ = decimal.NewFromString(pA)
		}

		if pB, err := ix.store.GetTokenBPriceByPool(context.Background(), pool.PoolID); err != nil {
			logger.Error(err, "cannot get the token_b price")
			continue
		} else if pB == "0" {
			continue
		} else {
			priceB, _ = decimal.NewFromString(pB)
		}

		vlA, _ := decimal.NewFromString(pool.ReserveA)
		vlB, _ := decimal.NewFromString(pool.ReserveA)
		tvl := vlA.Mul(priceA).Add(vlB.Mul(priceB))

		_, err = ix.store.UpdatePoolTV(context.Background(), db.UpdatePoolTVParams{
			PoolID:     pool.PoolID,
			TotalValue: tvl.String(),
		})
		if err != nil {
			logger.Error(err, "cannot get the token_b price")
			continue
		}
		s++
	}

	logger.Info("in " + strconv.Itoa(len(tokens)) + " pools, total values of " + strconv.Itoa(s) + " is synced")
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
