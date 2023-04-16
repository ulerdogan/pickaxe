package indexer

import (
	"context"
	"database/sql"
	"fmt"
	"math/big"

	rpc "github.com/ulerdogan/caigo-rpcv02/rpcv02"
	starknet "github.com/ulerdogan/pickaxe/clients/starknet"
	db "github.com/ulerdogan/pickaxe/db/sqlc"
	logger "github.com/ulerdogan/pickaxe/utils/logger"
)

func (ix *indexer) GetEvents(from, to uint64) error {
	keys, err := ix.store.GetAmmKeys(context.Background())
	if err != nil {
		logger.Error(err, "cannot get the amm keys")
		return err
	}

	events, err := ix.client.GetEvents(from, to, "", nil, keys)
	if err != nil {
		logger.Error(err, "cannot query the events")
		return err
	}

	ix.ixMutex.Lock()

	// TODO: use redis / message broker instead of in-memory array
	ix.Events = append(ix.Events, events...)

	ix.ixMutex.Unlock()

	logger.Info("new events queried for blocks: " + fmt.Sprint(from) + " - " + fmt.Sprint(to))

	// TODO: process the events

	ix.ixMutex.Lock()

	for i := len(ix.Events) - 1; i >= 0; i-- {
		if ok := processEvents(ix.client, ix.store, ix.Events[i]); ok {
			ix.Events = ix.Events[:i]
		}
	}

	ix.ixMutex.Unlock()

	return nil
}

func processEvents(client starknet.Client, store db.Store, event rpc.EmittedEvent) bool {
	pool, err := store.GetPoolByAddress(context.Background(), event.Event.FromAddress.String())
	if err != nil {
		if err == sql.ErrNoRows {
			return true
		}
		logger.Error(err, "cannot get the pool by address: "+event.Event.FromAddress.String())
		return false
	}

	dex, _ := client.NewDex(int(pool.AmmID))
	err = dex.SyncPoolFromEvent(starknet.PoolInfo{
		Address: event.Event.FromAddress.String(),
		Event:   event.Event,
		Block: big.NewInt(int64(event.BlockNumber)),
	}, store)
	if err != nil {
		logger.Error(err, "cannot sync pool from event: "+event.Event.FromAddress.String())
		return false
	}

	return true
}
