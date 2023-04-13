package indexer

import (
	"context"
	"fmt"

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

	// process the events

	return nil
}
