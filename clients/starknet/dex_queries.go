package starknet_client

import (
	"errors"

	"github.com/shopspring/decimal"
	db "github.com/ulerdogan/pickaxe/db/sqlc"
)

//TODO: add an event data input
type PoolInfo struct {
	Address  string
	ReserveA decimal.Decimal
	ReserveB decimal.Decimal
	ExtraData string
}

type Dex interface {
	SyncPoolFromFn(pool PoolInfo, store db.Store, client starknetClient) error
	SyncPoolFromEvent(pool PoolInfo, store db.Store) error
}

func (c *starknetClient) NewDex(amm_id int) (Dex, error) {
	switch amm_id {
	case 1:
		return newMyswap(), nil
	case 2:
		return newJediswap(), nil
	case 3:
		return  newSwap10k(), nil
	}

	return nil, errors.New("cannot find the dex")
}
