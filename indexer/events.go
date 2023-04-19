package indexer

import (
	"context"
	"database/sql"
	"fmt"
	"math/big"
	"strconv"

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

	events, err := getEventsLoop(from, to, keys, ix.client.GetEvents)
	if err != nil {
		logger.Error(err, "cannot query the events")
		return err
	}

	ix.ixMutex.Lock()

	// TODO: use redis / message broker instead of in-memory array
	ix.Events = append(ix.Events, events...)

	ix.ixMutex.Unlock()

	logger.Info("new events queried for blocks: " + fmt.Sprint(from) + " - " + fmt.Sprint(to))

	ix.ixMutex.Lock()

	const numWorkers = 10
	numJobs := len(ix.Events)
	jobs := make(chan rpc.EmittedEvent, numJobs)
	results := make(chan bool, numJobs)

	for w := 0; w < numWorkers; w++ {
		go func(jobs chan rpc.EmittedEvent, results chan bool) {
			processEventsConc(jobs, results, ix.store, ix.client)
		}(jobs, results)
	}

	for i := numJobs - 1; i >= 0; i-- {
		jobs <- ix.Events[i]
		ix.Events = ix.Events[:i]
	}
	close(jobs)

	var s int
	for res := 0; res < numJobs; res++ {
		if <-results {
			s++
		}
	}
	logger.Info("in " + strconv.Itoa(numJobs) + " events, " + strconv.Itoa(s) + " is processed")

	ix.ixMutex.Unlock()

	return nil
}

func processEventsConc(jobs <-chan rpc.EmittedEvent, results chan<- bool, store db.Store, client starknet.Client) {
	for event := range jobs {
		pool, err := store.GetPoolByAddress(context.Background(), event.Event.FromAddress.String())
		if err != nil {
			if err == sql.ErrNoRows {
				results <- true
				continue
			}
			logger.Error(err, "cannot get the pool by address: "+event.Event.FromAddress.String())
			results <- false
			continue
		}

		dex, _ := client.NewDex(int(pool.AmmID))
		err = dex.SyncPoolFromEvent(starknet.PoolInfo{
			Address: event.Event.FromAddress.String(),
			Event:   event.Event,
			Block:   big.NewInt(int64(event.BlockNumber)),
		}, store)
		if err != nil {
			logger.Error(err, "cannot sync pool from event: "+event.Event.FromAddress.String())
			results <- false
			continue
		}

		results <- true
	}
}

func getEventsLoop(from, to uint64, keys []string, getEvents func (from uint64, to uint64, address string, c_token *string, keys []string) ([]rpc.EmittedEvent, *string, error)) ([]rpc.EmittedEvent, error) {
	eventsArr := make([]rpc.EmittedEvent, 0)
	
	events, c_token, err := getEvents(from, to, "", nil, keys)
	if err != nil {
		return nil, err
	}

	eventsArr = append(eventsArr, events...)

	for c_token != nil {
		events, c_token, err = getEvents(from, to, "", c_token, keys)
		if err != nil {
			return nil, err
		}

		eventsArr = append(eventsArr, events...)
	}

	return eventsArr, nil
}
