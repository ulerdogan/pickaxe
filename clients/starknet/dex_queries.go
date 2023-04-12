package starknet_client

import (
	"github.com/shopspring/decimal"
	db "github.com/ulerdogan/pickaxe/db/sqlc"
)

// TODO: add an event data input
type PoolInfo struct {
	Address   string
	ReserveA  decimal.Decimal
	ReserveB  decimal.Decimal
	ExtraData string
}

type Dex interface {
	SyncPoolFromFn(pool PoolInfo, store db.Store, client Client) error
	SyncPoolFromEvent(pool PoolInfo, store db.Store) error
}
