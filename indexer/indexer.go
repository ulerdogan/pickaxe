package indexer

import (
	"context"
	"database/sql"
	"fmt"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-co-op/gocron"
	starknet "github.com/ulerdogan/pickaxe/clients/starknet"
	db "github.com/ulerdogan/pickaxe/db/sqlc"
	config "github.com/ulerdogan/pickaxe/utils/config"
	logger "github.com/ulerdogan/pickaxe/utils/logger"
)

type indexer struct {
	store  db.Store
	router *gin.Engine
	client starknet.Client
	config config.Config

	isIndexing  bool
	lastQueried *uint64
	poolIndexes chan Item

	scheduler *gocron.Scheduler
	ixMutex   *sync.Mutex
}

type Item struct {
	Index string
}

func NewIndexer(str db.Store, cli starknet.Client, cnfg config.Config) *indexer {
	ix := &indexer{
		store:       str,
		router:      gin.Default(),
		client:      cli,
		config:      cnfg,
		poolIndexes: make(chan Item),

		scheduler: gocron.NewScheduler(time.UTC),
		ixMutex:   &sync.Mutex{},
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
		ix.store.InitIndexer(context.Background(), sql.NullInt64{Int64: int64(lb), Valid: true})
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

func (ix *indexer) QueryBlocks() {
	if ix.isIndexing {
		return
	}

	ix.ixMutex.Lock()
	defer ix.ixMutex.Unlock()

	ix.isIndexing = true

	lastBlock, err := ix.client.LastBlock()
	if err != nil {
		logger.Error(err, "cannot get the last block")
		ix.isIndexing = false
		return
	}

	if lastBlock > *ix.lastQueried {
		// TODO: do sth
		logger.Info("new block catched: " + fmt.Sprint(lastBlock))
		ix.lastQueried = &lastBlock
		ix.store.UpdateIndexerStatus(context.Background(), sql.NullInt64{Int64: int64(lastBlock), Valid: true})
	} else {
		// FIXME: remove the part
		logger.Info("no new block: " + fmt.Sprint(lastBlock))
	}

	ix.isIndexing = false
}
