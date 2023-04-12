package rest_client

import (
	"github.com/shopspring/decimal"
	db "github.com/ulerdogan/pickaxe/db/sqlc"
)

type PriceAPI interface {
	GetPrice(client Client, token db.Token) (*decimal.Decimal, error)
}
