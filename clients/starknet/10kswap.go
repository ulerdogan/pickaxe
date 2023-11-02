package starknet_client

import (
	"context"
	"errors"

	"github.com/NethermindEth/starknet.go/rpc"
	"github.com/NethermindEth/starknet.go/types"
	db "github.com/ulerdogan/pickaxe/db/sqlc"
	logger "github.com/ulerdogan/pickaxe/utils/logger"
	utils "github.com/ulerdogan/pickaxe/utils/starknet"
)

type swap10k struct{}

func newSwap10k() Dex {
	return &swap10k{}
}

func (d *swap10k) SyncPoolFromFn(pool PoolInfo, store db.Store, client Client) error {
	pl, err := store.GetPoolByAddress(context.Background(), pool.Address)
	if err != nil {
		return err
	}

	paHash := GetAddressFelt(pool.Address)

	call, err := client.Call(rpc.FunctionCall{
		ContractAddress:    paHash,
		EntryPointSelector: types.GetSelectorFromNameFelt("getReserves"),
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

	rsA := utils.GetDecimal(call[0], int(tA.Decimals))
	rsB := utils.GetDecimal(call[1], int(tB.Decimals))

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

func (d *swap10k) SyncPoolFromEvent(pool PoolInfo, store db.Store) error {
	pl, err := store.GetPoolByAddress(context.Background(), GetAdressFormatFromStr(pool.Address))
	if err != nil {
		return err
	}

	if pl.LastBlock > pool.Block.Int64() {
		return nil
	}

	tA, _ := store.GetTokenByAddress(context.Background(), GetAdressFormatFromStr(pl.TokenA))
	tB, _ := store.GetTokenByAddress(context.Background(), GetAdressFormatFromStr(pl.TokenB))

	rsA := utils.GetDecimal(pool.Event.Data[0], int(tA.Decimals))
	rsB := utils.GetDecimal(pool.Event.Data[1], int(tB.Decimals))

	_, err = store.UpdatePoolReserves(context.Background(), db.UpdatePoolReservesParams{
		PoolID:    pl.PoolID,
		ReserveA:  rsA.String(),
		ReserveB:  rsB.String(),
		LastBlock: pool.Block.Int64(),
	})
	if err != nil {
		logger.Error(err, "cannot update pool reserves: "+pl.Address)
		return err
	}

	return nil
}

func (d *swap10k) SyncFee(pool PoolInfo, store db.Store, client Client) error {
	pl, err := store.GetPoolByAddress(context.Background(), GetAdressFormatFromStr(pool.Address))
	if err != nil {
		return err
	}

	store.UpdatePoolFee(context.Background(), db.UpdatePoolFeeParams{
		PoolID: pl.PoolID,
		Fee:    "0.3",
	})

	return nil
}
