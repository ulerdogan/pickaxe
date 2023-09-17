// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.21.0

package db

import (
	"database/sql"
	"time"
)

type Amm struct {
	AmmID         int64  `json:"amm_id"`
	DexName       string `json:"dex_name"`
	RouterAddress string `json:"router_address"`
	Key           string `json:"key"`
	AlgorithmType string `json:"algorithm_type"`
	// initialized
	CreatedAt time.Time `json:"created_at"`
}

type Indexer struct {
	ID               int32          `json:"id"`
	HashedPassword   string         `json:"hashed_password"`
	LastQueriedBlock sql.NullInt64  `json:"last_queried_block"`
	LastQueriedHash  sql.NullString `json:"last_queried_hash"`
	LastUpdated      sql.NullTime   `json:"last_updated"`
}

type Pool struct {
	PoolID           int64          `json:"pool_id"`
	Address          string         `json:"address"`
	AmmID            int64          `json:"amm_id"`
	TokenA           string         `json:"token_a"`
	TokenB           string         `json:"token_b"`
	ReserveA         string         `json:"reserve_a"`
	ReserveB         string         `json:"reserve_b"`
	Fee              string         `json:"fee"`
	TotalValue       string         `json:"total_value"`
	ExtraData        sql.NullString `json:"extra_data"`
	ExtraDataGeneral sql.NullString `json:"extra_data_general"`
	LastUpdated      time.Time      `json:"last_updated"`
	LastBlock        int64          `json:"last_block"`
}

type Token struct {
	Address   string    `json:"address"`
	Name      string    `json:"name"`
	Symbol    string    `json:"symbol"`
	Decimals  int32     `json:"decimals"`
	Base      bool      `json:"base"`
	Native    bool      `json:"native"`
	Ticker    string    `json:"ticker"`
	Price     string    `json:"price"`
	CreatedAt time.Time `json:"created_at"`
}
