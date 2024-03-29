package api

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
	AmmID            int64       `json:"amm_id"`
	PoolID           int64       `json:"pool_id"`
	Address          string      `json:"address"`
	TokenA           string      `json:"token_a"`
	TokenB           string      `json:"token_b"`
	ReserveA         string      `json:"reserve_a"`
	ReserveB         string      `json:"reserve_b"`
	TotalValue       string      `json:"total_value,omitempty"`
	Fee              interface{} `json:"fee,omitempty"`
	LastUpdated      string      `json:"last_updated"`
	LastBlock        int64       `json:"last_block"`
	ExtraData        string      `json:"extra_data,omitempty"`
	GeneralExtraData interface{} `json:"extra_data_general,omitempty"`
}

type AddPoolParams struct {
	Address          string      `json:"address"`
	TokenA           string      `json:"token_a"`
	TokenB           string      `json:"token_b"`
	AmmId            int64       `json:"amm_id"`
	ExtraData        string      `json:"extra_data,omitempty"`
	GeneralExtraData interface{} `json:"extra_data_general,omitempty"`
}

type AddEkuboPoolParams struct {
	Fee         string `json:"fee"`
	TickSpacing string `json:"tick_spacing"`
}

type RemovePoolParams struct {
	PoolID int64 `json:"pool_id"`
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
