package indexer

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"math/big"
	"strings"
	"time"

	"github.com/NethermindEth/starknet.go/rpc"
	starknet "github.com/ulerdogan/pickaxe/clients/starknet"
	logger "github.com/ulerdogan/pickaxe/utils/logger"
)

func (ix *Indexer) GetEvents(from uint64, to rpc.BlockHashAndNumberOutput) error {
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

	logger.Info("new events queried for blocks: " + fmt.Sprint(from) + " <-> " + fmt.Sprint(to.BlockNumber))

	return nil
}

func getEventsLoop(from uint64, to rpc.BlockHashAndNumberOutput, keys []string, ix *Indexer) error {
	var events []rpc.EmittedEvent
	var c_token string = ""
	var err error
	th := to.BlockHash

	for i := 0; i < 4; i++ {
		events, c_token, err = ix.Client.GetEventsWithID(
			rpc.BlockID{Number: &from},
			rpc.BlockID{Hash: to.BlockHash},
			"", c_token, keys)
		if err != nil {
			if strings.Compare(err.Error(), "Block not found") == 0 {
				if i == 3 {
					return err
				}
				logger.Info("requerying to wait block sync. block hash: " + to.BlockHash.String())
				time.Sleep(5 * time.Second)
			} else {
				return err
			}
		} else {
			break
		}
	}

	done := make(chan bool)
	go func(dn chan bool) {
		for _, event := range events {
			msEvent, _ := json.Marshal(event)
			ix.PublishRmqMsg(msEvent)
		}
		dn <- true
	}(done)

	for c_token != "" {
		for i := 0; i < 4; i++ {
			events, c_token, err = ix.Client.GetEventsWithID(
				rpc.BlockID{Number: &from},
				rpc.BlockID{Hash: th},
				"", c_token, keys)
			if err != nil {
				if strings.Compare(err.Error(), "Block not found") == 0 {
					if i == 3 {
						return err
					}
					logger.Info("requerying to wait block sync. block hash: " + to.BlockHash.String())
					time.Sleep(5 * time.Second)
				} else {
					return err
				}
			} else {
				break
			}
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
		pool, err := ix.Store.GetPoolByAddress(context.Background(), starknet.GetAdressFormatFromFelt(event.Event.FromAddress))
		if err != nil {
			if err == sql.ErrNoRows {
				continue
			}
			logger.Error(err, "cannot get the pool by address: "+starknet.GetAdressFormatFromFelt(event.Event.FromAddress))
			continue
		}

		dex, _ := ix.Client.NewDex(int(pool.AmmID))
		err = dex.SyncPoolFromEvent(starknet.PoolInfo{
			Address: starknet.GetAdressFormatFromFelt(event.Event.FromAddress),
			Event:   event.Event,
			Block:   big.NewInt(int64(event.BlockNumber)),
		}, ix.Store)
		if err != nil {
			logger.Error(err, "cannot sync pool from event: "+starknet.GetAdressFormatFromFelt(event.Event.FromAddress))
			continue
		}
	}
}
