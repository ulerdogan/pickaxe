package starknet_client

import (
	"math/big"

	"github.com/NethermindEth/starknet.go/rpc"
	"github.com/shopspring/decimal"
	db "github.com/ulerdogan/pickaxe/db/sqlc"
)

type PoolInfo struct {
	Address          string
	ReserveA         decimal.Decimal
	ReserveB         decimal.Decimal
	ExtraData        string
	GeneralExtraData string
	Fee              string
	Event            rpc.Event
	Block            *big.Int
}

type Dex interface {
	SyncPoolFromFn(pool PoolInfo, store db.Store, client Client) error
	SyncPoolFromEvent(pool PoolInfo, store db.Store) error
	SyncFee(pool PoolInfo, store db.Store, client Client) error
}
