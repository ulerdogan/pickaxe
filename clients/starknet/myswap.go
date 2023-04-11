package starknet_client

import db "github.com/ulerdogan/pickaxe/db/sqlc"

type myswap struct{}

func newMyswap() Dex {
	return &myswap{}
}

func (d *myswap) SyncPoolFromFn(pool PoolInfo, store db.Store) error {
	return nil
}

func (d *myswap) SyncPoolFromEvent(pool PoolInfo, store db.Store) error {
	return nil
}
