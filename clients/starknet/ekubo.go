package starknet_client

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"

	"github.com/NethermindEth/juno/core/felt"
	"github.com/NethermindEth/starknet.go/rpc"
	"github.com/NethermindEth/starknet.go/types"
	db "github.com/ulerdogan/pickaxe/db/sqlc"
)

type ekubo struct{}

type EkuboData struct {
	TickSpacing   string `json:"tick_spacing"`
	SqrtPriceLow  string `json:"sqrt_price_low,omitempty"`
	SqrtPriceHigh string `json:"sqrt_price_high,omitempty"`
	CurrentTick   string `json:"current_tick,omitempty"`
	TickSign      string `json:"tick_sign,omitempty"`
	Liqudity      string `json:"liqudity,omitempty"`
	KeyExtension  string `json:"key_extension"`
}

func newEkubo() Dex {
	return &ekubo{}
}

func (d *ekubo) SyncPoolFromFn(pool PoolInfo, store db.Store, client Client) error {
	pl, err := store.GetPoolByAddressExtra(context.Background(), db.GetPoolByAddressExtraParams{
		Address:   pool.Address,
		ExtraData: sql.NullString{String: pool.ExtraData, Valid: true},
	})
	if err != nil {
		return err
	}

	paHash := GetAddressFelt(pl.Address)
	var data EkuboData
	json.Unmarshal([]byte(pl.GeneralExtraData.String), &data)

	calldata := []*felt.Felt{GetAddressFelt(pl.TokenA), GetAddressFelt(pl.TokenB), GetAddressFelt(pl.Fee), GetAddressFelt(data.TickSpacing), GetAddressFelt(data.KeyExtension)}

	call, err := client.Call(rpc.FunctionCall{
		ContractAddress:    paHash,
		EntryPointSelector: types.GetSelectorFromNameFelt("get_pool_price"),
		Calldata:           calldata,
	})
	if err != nil {
		return errors.New("starknet query error")
	}

	data.SqrtPriceLow, data.SqrtPriceHigh = call[0].String(), call[1].String()
	data.CurrentTick, data.TickSign = call[2].String(), call[3].String()

	call, err = client.Call(rpc.FunctionCall{
		ContractAddress:    paHash,
		EntryPointSelector: types.GetSelectorFromNameFelt("get_pool_liquidity"),
		Calldata:           calldata,
	})
	if err != nil {
		return errors.New("starknet query error")
	}

	data.Liqudity = call[0].String()

	jsonBytes, _ := json.Marshal(data)

	_, err = store.UpdatePoolGeneralExtraData(context.Background(), db.UpdatePoolGeneralExtraDataParams{
		PoolID:           pl.PoolID,
		GeneralExtraData: sql.NullString{String: string(jsonBytes), Valid: true},
		LastBlock:        pool.Block.Int64(),
	})
	if err != nil {
		return err
	}

	return nil
}

func (d *ekubo) SyncPoolFromEvent(pool PoolInfo, store db.Store) error {
	tokenA := GetAdressFormatFromFelt(pool.Event.Data[1])
	tokenB := GetAdressFormatFromFelt(pool.Event.Data[2])
	fee := pool.Event.Data[3].String()
	tickSpacing := pool.Event.Data[4].String()

	ekuboHash := GetUniqueEkuboHash(tokenA, tokenB, fee, tickSpacing)

	pl, err := store.GetPoolByAddressExtra(context.Background(), db.GetPoolByAddressExtraParams{
		Address:   pool.Address,
		ExtraData: sql.NullString{String: ekuboHash, Valid: true},
	})
	if err != nil {
		return err
	}

	if pl.LastBlock > pool.Block.Int64() {
		return nil
	}

	var data EkuboData
	json.Unmarshal([]byte(pl.GeneralExtraData.String), &data)

	data.SqrtPriceLow, data.SqrtPriceHigh = pool.Event.Data[16].String(), pool.Event.Data[17].String()
	data.CurrentTick, data.TickSign = pool.Event.Data[18].String(), pool.Event.Data[19].String()
	data.Liqudity = pool.Event.Data[20].String()

	jsonBytes, _ := json.Marshal(data)

	_, err = store.UpdatePoolGeneralExtraData(context.Background(), db.UpdatePoolGeneralExtraDataParams{
		PoolID:           pl.PoolID,
		GeneralExtraData: sql.NullString{String: string(jsonBytes), Valid: true},
		LastBlock:        pool.Block.Int64(),
	})
	if err != nil {
		return err
	}

	return nil
}

func (d *ekubo) SyncFee(pool PoolInfo, store db.Store, client Client) error {
	pl, err := store.GetPoolByAddressExtra(context.Background(), db.GetPoolByAddressExtraParams{
		Address:   pool.Address,
		ExtraData: sql.NullString{String: pool.ExtraData, Valid: true},
	})
	if err != nil {
		return err
	}

	_, err = store.UpdatePoolFee(context.Background(), db.UpdatePoolFeeParams{
		PoolID: pl.PoolID,
		Fee:    pool.Fee,
	})

	return err
}
