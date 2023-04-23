package indexer

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"math/big"
	"strconv"
	"time"

	rpc "github.com/ulerdogan/caigo-rpcv02/rpcv02"
	starknet "github.com/ulerdogan/pickaxe/clients/starknet"
	logger "github.com/ulerdogan/pickaxe/utils/logger"
)

func (ix *indexer) getEvents(from, to uint64) error {
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

	for _, event := range events {
		msEvent, _ := json.Marshal(event)
		ix.publishRmqMsg(msEvent)
	}

	logger.Info("new events queried for blocks (" + strconv.Itoa(len(events)) + "): " + fmt.Sprint(from) + " <-> " + fmt.Sprint(to))

	return nil
}

func (ix *indexer) processEvents() {
	msgs, err := ix.rabbitmq.Consume(
		"EventsQueue",
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		logger.Error(err, "cannot consume the events queue")

		time.Sleep(3 * time.Second)
		go ix.processEvents()
		return
	}

	for e := range msgs {
		var event rpc.EmittedEvent
		if err := json.Unmarshal(e.Body, &event); err != nil {
			logger.Error(err, "cannot unmarshal the event msg from the rabbitmq")
			continue
		}

		pool, err := ix.store.GetPoolByAddress(context.Background(), event.Event.FromAddress.String())
		if err != nil {
			if err == sql.ErrNoRows {
				continue
			}
			logger.Error(err, "cannot get the pool by address: "+event.Event.FromAddress.String())
			continue
		}

		dex, _ := ix.client.NewDex(int(pool.AmmID))
		err = dex.SyncPoolFromEvent(starknet.PoolInfo{
			Address: event.Event.FromAddress.String(),
			Event:   event.Event,
			Block:   big.NewInt(int64(event.BlockNumber)),
		}, ix.store)
		if err != nil {
			logger.Error(err, "cannot sync pool from event: "+event.Event.FromAddress.String())
			continue
		}
	}
}

func getEventsLoop(from, to uint64, keys []string, getEvents func(from uint64, to uint64, address string, c_token *string, keys []string) ([]rpc.EmittedEvent, *string, error)) ([]rpc.EmittedEvent, error) {
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
