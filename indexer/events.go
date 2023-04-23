package indexer

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"math/big"
	"time"

	rpc "github.com/ulerdogan/caigo-rpcv02/rpcv02"
	starknet "github.com/ulerdogan/pickaxe/clients/starknet"
	logger "github.com/ulerdogan/pickaxe/utils/logger"
)

func (ix *Indexer) GetEvents(from, to uint64) error {
	keys, err := ix.Store.GetAmmKeys(context.Background())
	if err != nil {
		logger.Error(err, "cannot get the amm keys")
		return err
	}

	err = getEventsLoop(from, to, keys, ix)
	if err != nil {
		logger.Error(err, "cannot query the events")
		return err
	}

	logger.Info("new events queried for blocks: " + fmt.Sprint(from) + " <-> " + fmt.Sprint(to))

	return nil
}

func getEventsLoop(from, to uint64, keys []string, ix *Indexer) error {
	events, c_token, err := ix.Client.GetEvents(from, to, "", nil, keys)
	if err != nil {
		return err
	}

	done := make(chan bool)
	go func(dn chan bool) {
		for _, event := range events {
			msEvent, _ := json.Marshal(event)
			ix.PublishRmqMsg(msEvent)
		}
		dn <- true
	}(done)

	for c_token != nil {
		events, c_token, err = ix.Client.GetEvents(from, to, "", c_token, keys)
		if err != nil {
			return err
		}

		<-done
		go func(dn chan bool) {
			for _, event := range events {
				msEvent, _ := json.Marshal(event)
				ix.PublishRmqMsg(msEvent)
			}
			dn <- true
		}(done)
	}

	return nil
}

func (ix *Indexer) ProcessEvents() {
	msgs, err := ix.RabbitMQ.Consume(
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
		go ix.ProcessEvents()
		return
	}

	for e := range msgs {
		var event rpc.EmittedEvent
		if err := json.Unmarshal(e.Body, &event); err != nil {
			logger.Error(err, "cannot unmarshal the event msg from the rabbitmq")
			continue
		}

		// FIXME: pool processing may not be compatible with each type of the pools in future
		pool, err := ix.Store.GetPoolByAddress(context.Background(), event.Event.FromAddress.String())
		if err != nil {
			if err == sql.ErrNoRows {
				continue
			}
			logger.Error(err, "cannot get the pool by address: "+event.Event.FromAddress.String())
			continue
		}

		dex, _ := ix.Client.NewDex(int(pool.AmmID))
		err = dex.SyncPoolFromEvent(starknet.PoolInfo{
			Address: event.Event.FromAddress.String(),
			Event:   event.Event,
			Block:   big.NewInt(int64(event.BlockNumber)),
		}, ix.Store)
		if err != nil {
			logger.Error(err, "cannot sync pool from event: "+event.Event.FromAddress.String())
			continue
		}
	}
}
