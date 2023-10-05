package starknet_client

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/NethermindEth/starknet.go/rpc"
	"github.com/NethermindEth/starknet.go/types"
	"github.com/shopspring/decimal"
	db "github.com/ulerdogan/pickaxe/db/sqlc"
	logger "github.com/ulerdogan/pickaxe/utils/logger"
	utils "github.com/ulerdogan/pickaxe/utils/starknet"
)

type sithFees struct {
	Fee0 string `json:"fee_0"`
	Fee1 string `json:"fee_1"`
}

type sithswap struct{}

func newSithswap() Dex {
	return &sithswap{}
}

func (d *sithswap) SyncPoolFromFn(pool PoolInfo, store db.Store, client Client) error {
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
	rsB := utils.GetDecimal(call[2], int(tB.Decimals))

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

func (d *sithswap) SyncPoolFromEvent(pool PoolInfo, store db.Store) error {
	pl, err := store.GetPoolByAddress(context.Background(), pool.Address)
	if err != nil {
		return err
	}

	if pl.LastBlock >= pool.Block.Int64() {
		return nil
	}

	tA, _ := store.GetTokenByAddress(context.Background(), pl.TokenA)
	tB, _ := store.GetTokenByAddress(context.Background(), pl.TokenB)

	rsA := utils.GetDecimal(pool.Event.Data[0], int(tA.Decimals))
	rsB := utils.GetDecimal(pool.Event.Data[2], int(tB.Decimals))

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

func (d *sithswap) SyncFee(pool PoolInfo, store db.Store, client Client) error {
	pl, err := store.GetPoolByAddress(context.Background(), pool.Address)
	if err != nil {
		return err
	}

	paHash := GetAddressFelt(pool.Address)

	call0, err := client.Call(rpc.FunctionCall{
		ContractAddress:    paHash,
		EntryPointSelector: types.GetSelectorFromNameFelt("getFee0"),
	})
	if err != nil {
		return errors.New("starknet query error")
	}
	call1, err := client.Call(rpc.FunctionCall{
		ContractAddress:    paHash,
		EntryPointSelector: types.GetSelectorFromNameFelt("getFee1"),
	})
	if err != nil {
		return errors.New("starknet query error")
	}

	fees := sithFees{
		Fee0: decimal.NewFromInt(types.HexToBN(call0[0].String()).Int64()).Div(decimal.NewFromInt(10000)).String(),
		Fee1: decimal.NewFromInt(types.HexToBN(call1[0].String()).Int64()).Div(decimal.NewFromInt(10000)).String(),
	}
	jsonBytes, _ := json.Marshal(fees)

	_, err = store.UpdatePoolFee(context.Background(), db.UpdatePoolFeeParams{
		PoolID: pl.PoolID,
		Fee:    string(jsonBytes),
	})

	return err
}
