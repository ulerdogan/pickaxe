package indexer

import (
	"context"
	"database/sql"
	"fmt"
	"sync"
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

type indexer struct {
	store  db.Store
	client starknet.Client
	rest   rest.Client
	config config.Config

	rabbitmq *amqp.Connection
	// FIXME: remove the array
	Events []rpc.EmittedEvent

	lastQueried *uint64

	scheduler *gocron.Scheduler
	ixMutex   *sync.Mutex
	stMutex   *sync.Mutex
}

func NewIndexer(str db.Store, cli starknet.Client, rs rest.Client, cnfg config.Config, rmq *amqp.Connection) *indexer {
	ix := &indexer{
		store:  str,
		client: cli,
		rest:   rs,
		config: cnfg,

		rabbitmq: rmq,
		// FIXME: remove the array
		Events: make([]rpc.EmittedEvent, 0),

		scheduler: gocron.NewScheduler(time.UTC),
		ixMutex:   &sync.Mutex{},
		stMutex:   &sync.Mutex{},
	}

	ix.syncBlockFromDB()
	return ix
}

func (ix *indexer) syncBlockFromDB() {
	ix.ixMutex.Lock()
	defer ix.ixMutex.Unlock()

	// set indexer records in db if not exists
	ixStatus, err := ix.store.GetIndexerStatus(context.Background())
	if err == sql.ErrNoRows || ixStatus.LastQueried.Int64 == 0 {
		lb, err := ix.client.LastBlock()
		ix.lastQueried = &lb
		ix.store.InitIndexer(context.Background(), db.InitIndexerParams{
			HashedPassword: hasher.HashPassword(ix.config.AuthPassword),
			LastQueried:    sql.NullInt64{Int64: int64(lb), Valid: true},
		})
		logger.Info("indexer initialized with the last block in the db: " + fmt.Sprint(lb))
		if err != nil {
			logger.Error(err, "cannot get the last block")
			return
		}
	} else {
		lq := uint64(ixStatus.LastQueried.Int64)
		ix.lastQueried = &lq
		logger.Info("indexer synced from the db: " + fmt.Sprint(lq))
	}
}
