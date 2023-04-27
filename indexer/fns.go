package indexer

import (
	"context"
	"math/big"

	starknet "github.com/ulerdogan/pickaxe/clients/starknet"
	db "github.com/ulerdogan/pickaxe/db/sqlc"
	logger "github.com/ulerdogan/pickaxe/utils/logger"
)

func (ix *Indexer) UpdateByFnsAll(block uint64) {
	pools, err := ix.Store.GetAllPoolsWithoutKeys(context.Background())
	if err != nil {
		logger.Error(err, "cannot get the pools")
		return
	}

	processByFns(pools, block, ix.Store, ix.Client)
}

func (ix *Indexer) UpdateByFns(block uint64) {
	pools, err := ix.Store.GetAllPoolsWithoutKeys(context.Background())
	if err != nil {
		logger.Error(err, "cannot get the pools")
		return
	}

	processByFns(pools, block, ix.Store, ix.Client)
}

func processByFns(pools []db.PoolsV2, block uint64, store db.Store, client starknet.Client) {
	const numWorkers = 10
	jobs := make(chan db.PoolsV2, len(pools))
	results := make(chan bool, len(pools))

	for w := 0; w < numWorkers; w++ {
		go func(jobs chan db.PoolsV2, results chan bool) {
			syncPoolFromFnConc(jobs, results, block, store, client)
		}(jobs, results)
	}

	for _, pool := range pools {
		jobs <- pool
	}
	close(jobs)

	var s int
	for res := 0; res < len(pools); res++ {
		if <-results {
			s++
		}
	}
}

func syncPoolFromFnConc(jobs <-chan db.PoolsV2, results chan<- bool, lastBlock uint64, store db.Store, client starknet.Client) {
	for pool := range jobs {
		dex, _ := client.NewDex(int(pool.AmmID))

		err := dex.SyncPoolFromFn(starknet.PoolInfo{
			Address:   pool.Address,
			ExtraData: pool.ExtraData.String,
			Block:     big.NewInt(int64(lastBlock)),
		}, store, client)
		if err != nil {
			logger.Error(err, "sync pool error: "+pool.Address)
			results <- false
			continue
		}

		results <- true
	}
}
