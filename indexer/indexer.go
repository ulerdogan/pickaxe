package indexer

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/go-co-op/gocron"
	"github.com/streadway/amqp"
	rpc "github.com/ulerdogan/caigo-rpcv02/rpcv02"
	rest "github.com/ulerdogan/pickaxe/clients/rest"
	starknet "github.com/ulerdogan/pickaxe/clients/starknet"
	db "github.com/ulerdogan/pickaxe/db/sqlc"
	config "github.com/ulerdogan/pickaxe/utils/config"
	hasher "github.com/ulerdogan/pickaxe/utils/hasher"
	logger "github.com/ulerdogan/pickaxe/utils/logger"
)

type Indexer struct {
	Store       db.Store
	Client      starknet.Client
	Rest        rest.Client
	Config      config.Config
	RabbitMQ    *amqp.Channel
	LastQueried *rpc.BlockHashAndNumberOutput
	Scheduler   *gocron.Scheduler
}

func NewIndexer(str db.Store, cli starknet.Client, rs rest.Client, cnfg config.Config, rmq *amqp.Channel) *Indexer {
	ix := &Indexer{
		Store:       str,
		Client:      cli,
		Rest:        rs,
		Config:      cnfg,
		RabbitMQ:    rmq,
		LastQueried: nil,
		Scheduler:   gocron.NewScheduler(time.UTC),
	}

	ix.syncBlockFromDB()
	return ix
}

func (ix *Indexer) syncBlockFromDB() {
	// set indexer records in db if not exists
	ixStatus, err := ix.Store.GetIndexerStatus(context.Background())
	if err == sql.ErrNoRows || ixStatus.LastQueriedBlock.Int64 == 0 {
		lb, err := ix.Client.LastBlock()
		if err != nil {
			logger.Error(err, "cannot get the last block")
			return
		}
		ix.LastQueried = lb
		ix.Store.InitIndexer(context.Background(), db.InitIndexerParams{
			HashedPassword:   hasher.HashPassword(ix.Config.AuthPassword),
			LastQueriedBlock: sql.NullInt64{Int64: int64(lb.BlockNumber), Valid: true},
			LastQueriedHash:  sql.NullString{String: lb.BlockHash, Valid: true},
		})
		logger.Info("indexer initialized with the last block: " + fmt.Sprint(lb.BlockNumber))
		if err != nil {
			logger.Error(err, "cannot get the last block")
			return
		}
	} else {
		lq := uint64(ixStatus.LastQueriedBlock.Int64)
		ix.LastQueried = &rpc.BlockHashAndNumberOutput{BlockNumber: lq}
		logger.Info("indexer synced from the db: " + fmt.Sprint(lq))
	}
}
