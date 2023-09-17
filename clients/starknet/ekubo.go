package starknet_client

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"

	"github.com/dontpanicdao/caigo/types"
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

	paHash := types.HexToHash(pool.Address)

	var data EkuboData
	json.Unmarshal([]byte(pool.ExtraDataGeneral), &pl.ExtraDataGeneral)

	calldata := []string{pl.TokenA, pl.TokenB, pl.Fee, data.TickSpacing, data.KeyExtension}

	call, err := client.Call(types.FunctionCall{
		ContractAddress:    paHash,
		EntryPointSelector: "get_pool_price",
		Calldata:           calldata,
	})
	if err != nil {
		return errors.New("starknet query error")
	}

	data.SqrtPriceLow, data.SqrtPriceHigh = call[0], call[1]
	data.CurrentTick, data.TickSign = call[2], call[3]

	call, err = client.Call(types.FunctionCall{
		ContractAddress:    paHash,
		EntryPointSelector: "get_pool_liquidity",
		Calldata:           calldata,
	})
	if err != nil {
		return errors.New("starknet query error")
	}

	data.Liqudity = call[0]

	jsonBytes, _ := json.Marshal(data)

	_, err = store.UpdatePoolReservesWithExtraData(context.Background(), db.UpdatePoolReservesWithExtraDataParams{
		PoolID:    pl.PoolID,
		ExtraData: sql.NullString{String: string(jsonBytes), Valid: true},
		LastBlock: pool.Block.Int64(),
	})
	if err != nil {
		return err
	}

	return nil
}

func (d *ekubo) SyncPoolFromEvent(pool PoolInfo, store db.Store) error {
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

	store.UpdatePoolFee(context.Background(), db.UpdatePoolFeeParams{
		PoolID: pl.PoolID,
		Fee:    pl.Fee,
	})

	return nil
}
