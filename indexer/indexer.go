package indexer

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/go-co-op/gocron"
	"github.com/streadway/amqp"
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
	LastQueried *uint64
	Scheduler   *gocron.Scheduler
}

func NewIndexer(str db.Store, cli starknet.Client, rs rest.Client, cnfg config.Config, rmq *amqp.Channel) *Indexer {
	ix := &Indexer{
		Store:     str,
		Client:    cli,
		Rest:      rs,
		Config:    cnfg,
		RabbitMQ:  rmq,
		Scheduler: gocron.NewScheduler(time.UTC),
	}

	ix.syncBlockFromDB()
	return ix
}

func (ix *Indexer) syncBlockFromDB() {
	// set indexer records in db if not exists
	ixStatus, err := ix.Store.GetIndexerStatus(context.Background())
	if err == sql.ErrNoRows || ixStatus.LastQueried.Int64 == 0 {
		lb, err := ix.Client.LastBlock()
		ix.LastQueried = &lb
		ix.Store.InitIndexer(context.Background(), db.InitIndexerParams{
			HashedPassword: hasher.HashPassword(ix.Config.AuthPassword),
			LastQueried:    sql.NullInt64{Int64: int64(lb), Valid: true},
		})
		logger.Info("indexer initialized with the last block in the db: " + fmt.Sprint(lb))
		if err != nil {
			logger.Error(err, "cannot get the last block")
			return
		}
	} else {
		lq := uint64(ixStatus.LastQueried.Int64)
		ix.LastQueried = &lq
		logger.Info("indexer synced from the db: " + fmt.Sprint(lq))
	}
}
