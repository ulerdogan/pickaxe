package indexer

import (
	"context"
	"fmt"

	rpc "github.com/ulerdogan/caigo-rpcv02/rpcv02"
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
		if ok := processEvents(ix.Events[i]); ok {
			ix.Events = ix.Events[:i]
		}
	}

	ix.ixMutex.Unlock()

	return nil
}

func processEvents(event rpc.EmittedEvent) bool {
	_ = event
	return true
}
