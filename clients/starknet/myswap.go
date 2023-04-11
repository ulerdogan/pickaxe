package starknet_client

import (
	"context"
	"database/sql"
	"errors"

	"github.com/dontpanicdao/caigo/types"
	db "github.com/ulerdogan/pickaxe/db/sqlc"
	utils "github.com/ulerdogan/pickaxe/utils/starknet"
)

type myswap struct{}

func newMyswap() Dex {
	return &myswap{}
}

func (d *myswap) SyncPoolFromFn(pool PoolInfo, store db.Store, client Client) error {
	pl, err := store.GetPoolByAddressExtra(context.Background(), db.GetPoolByAddressExtraParams{
		Address:   pool.Address,
		ExtraData: sql.NullString{String: pool.ExtraData, Valid: true},
	})
	if err != nil {
		return err
	}

	paHash := types.HexToHash(pool.Address)

	call, err := client.Call(types.FunctionCall{
		ContractAddress:    paHash,
		EntryPointSelector: "get_pool",
		Calldata:           []string{pool.ExtraData},
	})
	if err != nil {
		return errors.New("starknet query error")
	}

	tA, _ := store.GetTokenByAddress(context.Background(), pl.TokenA)
	tB, _ := store.GetTokenByAddress(context.Background(), pl.TokenA)

	rsA := utils.GetDecimal(call[2], int(tA.Decimals))
	rsB := utils.GetDecimal(call[5], int(tB.Decimals))
	pool.ReserveA, pool.ReserveB = rsA, rsB

	_, err = store.UpdatePoolReserves(context.Background(), db.UpdatePoolReservesParams{
		PoolID:   pl.PoolID,
		ReserveA: rsA.String(),
		ReserveB: rsB.String(),
	})
	if err != nil {
		return err
	}

	return nil
}

func (d *myswap) SyncPoolFromEvent(pool PoolInfo, store db.Store) error {
	return nil
}
