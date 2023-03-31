package indexer

import (
	"sync"

	db "github.com/ulerdogan/pickaxe/db/sqlc"
	config "github.com/ulerdogan/pickaxe/utils/config"
)

type indexer struct {
	Store      db.Store
	Config     config.Config
	storeMutex *sync.Mutex
	stateMutex *sync.Mutex
}

func NewIndexer(store db.Store, cnfg config.Config) *indexer {
	i := &indexer{
		Store:      store,
		Config:     cnfg,
		storeMutex: &sync.Mutex{},
		stateMutex: &sync.Mutex{},
	}

	return i
}
