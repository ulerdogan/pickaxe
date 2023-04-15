package starknet_client

import (
	"math/big"

	"github.com/shopspring/decimal"
	rpc "github.com/ulerdogan/caigo-rpcv02/rpcv02"
	db "github.com/ulerdogan/pickaxe/db/sqlc"
)

// TODO: add an event data input
type PoolInfo struct {
	Address   string
	ReserveA  decimal.Decimal
	ReserveB  decimal.Decimal
	ExtraData string
	Event     rpc.Event
	Block     *big.Int
}

type Dex interface {
	SyncPoolFromFn(pool PoolInfo, store db.Store, client Client) error
	SyncPoolFromEvent(pool PoolInfo, store db.Store) error
}
