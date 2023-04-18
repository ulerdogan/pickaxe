package api

import "github.com/shopspring/decimal"

type TokenResponse struct {
	Address  string `json:"address"`
	Name     string `json:"name"`
	Symbol   string `json:"symbol"`
	Decimals int32  `json:"decimals"`
	Price    string `json:"price,omitempty"`
}

type AddTokenParams struct {
	Address string `json:"address"`
	Name    string `json:"name"`
	Symbol  string `json:"symbol"`
	Ticker  string `json:"ticker,omitempty"`
}

type PoolResponse struct {
	Address     string `json:"address"`
	TokenA      string `json:"token_a"`
	TokenB      string `json:"token_b"`
	ReserveA    string `json:"reserve_a"`
	ReserveB    string `json:"reserve_b"`
	Fee         string `json:"fee"`
	TotalValue  string `json:"total_value,omitempty"`
	LastUpdated string `json:"last_updated"`
	LastBlock   int64  `json:"last_block"`
}

type AddPoolParams struct {
	Address   string          `json:"address"`
	TokenA    string          `json:"token_a"`
	TokenB    string          `json:"token_b"`
	AmmId     int64           `json:"amm_id"`
	Fee       decimal.Decimal `json:"fee"`
	ExtraData string          `json:"extra_data,omitempty"`
}

type AmmResponse struct {
	Name    string `json:"name"`
	Address string `json:"router_address"`
}

type AddAmmParams struct {
	Name          string `json:"name"`
	Address       string `json:"router_address"`
	Key           string `json:"key,omitempty"`
	AlgorithmType string `json:"algorithm_type"`
}

type AuthRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type AuthResponse struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int64  `json:"expires_in"`
}
