package starknet_client

import db "github.com/ulerdogan/pickaxe/db/sqlc"

type jediswap struct{}

func newJediswap() Dex {
	return &jediswap{}
}

func (d *jediswap) SyncPoolFromFn(pool PoolInfo, store db.Store) error {
	return nil
}

func (d *jediswap) SyncPoolFromEvent(pool PoolInfo, store db.Store) error {
	return nil
}
