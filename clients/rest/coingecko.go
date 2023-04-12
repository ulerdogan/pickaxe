package rest_client

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"

	"github.com/shopspring/decimal"
	db "github.com/ulerdogan/pickaxe/db/sqlc"
)

type coingecko struct{}

func NewCoingeckoAPI() PriceAPI {
	return &coingecko{}
}

func (p *coingecko) GetPrice(client Client, token db.Token) (*decimal.Decimal, error) {
	d := fmt.Sprintf("https://api.coingecko.com/api/v3/simple/price?ids=%s&vs_currencies=usd", token.Ticker)
	res, err := client.Get(d, nil)
	if err != nil {
		return nil, errors.New("coingecko query error in the request")
	}
	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, errors.New("coingecko query error in reading the response")
	}

	var price map[string]struct {
		Usd decimal.Decimal
	}

	json.Unmarshal(body, &price)
	mp := price[token.Ticker].Usd

	return &mp, nil
}
