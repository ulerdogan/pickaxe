package starknet_client

import db "github.com/ulerdogan/pickaxe/db/sqlc"

type swap10k struct{}

func newSwap10k() Dex {
	return &swap10k{}
}

func (d *swap10k) SyncPoolFromFn(pool PoolInfo, store db.Store) error {
	return nil
}

func (d *swap10k) SyncPoolFromEvent(pool PoolInfo, store db.Store) error {
	return nil
}
