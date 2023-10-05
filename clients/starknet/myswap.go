package starknet_client

import (
	"context"
	"database/sql"
	"errors"

	"github.com/NethermindEth/juno/core/felt"
	"github.com/NethermindEth/starknet.go/rpc"
	"github.com/NethermindEth/starknet.go/types"
	db "github.com/ulerdogan/pickaxe/db/sqlc"
	logger "github.com/ulerdogan/pickaxe/utils/logger"
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

	paHash := GetAddressFelt(pool.Address)

	call, err := client.Call(rpc.FunctionCall{
		ContractAddress:    paHash,
		EntryPointSelector: types.GetSelectorFromNameFelt("get_pool"),
		Calldata:           []*felt.Felt{getStrBigIntFelt(pool.ExtraData)},
	})
	if err != nil {
		return errors.New("starknet query error")
	}

	tA, err := store.GetTokenByAddress(context.Background(), pl.TokenA)
	if err != nil {
		logger.Error(err, "cannot get token A: "+pl.TokenA)
		return err
	}
	tB, err := store.GetTokenByAddress(context.Background(), pl.TokenB)
	if err != nil {
		logger.Error(err, "cannot get token B: "+pl.TokenB)
		return err
	}

	rsA := utils.GetDecimal(call[2], int(tA.Decimals))
	rsB := utils.GetDecimal(call[5], int(tB.Decimals))

	_, err = store.UpdatePoolReserves(context.Background(), db.UpdatePoolReservesParams{
		PoolID:    pl.PoolID,
		ReserveA:  rsA.String(),
		ReserveB:  rsB.String(),
		LastBlock: pool.Block.Int64(),
	})
	if err != nil {
		return err
	}

	return nil
}

func (d *myswap) SyncPoolFromEvent(pool PoolInfo, store db.Store) error {
	return errors.New("myswap sync from even is not implemented")
}

func (d *myswap) SyncFee(pool PoolInfo, store db.Store, client Client) error {
	pl, err := store.GetPoolByAddressExtra(context.Background(), db.GetPoolByAddressExtraParams{
		Address:   pool.Address,
		ExtraData: sql.NullString{String: pool.ExtraData, Valid: true},
	})
	if err != nil {
		return err
	}

	_, err = store.UpdatePoolFee(context.Background(), db.UpdatePoolFeeParams{
		PoolID: pl.PoolID,
		Fee:    "0.3",
	})

	return err
}
