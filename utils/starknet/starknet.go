package starknet_utils

import (
	"github.com/dontpanicdao/caigo/types"
	"github.com/shopspring/decimal"
	"math/big"
)

func GetStrFormat(s string) string {
	return types.StrToFelt(s).Big().String()
}

func GetDecimal(s string, dc int) decimal.Decimal {
	b := types.StrToFelt(s).Big()
	return toDecimal(b, dc)
}

func toDecimal(value *big.Int, decimals int) decimal.Decimal {
	mul := decimal.NewFromFloat(float64(10)).Pow(decimal.NewFromFloat(float64(decimals)))
	num, _ := decimal.NewFromString(value.String())
	result := num.Div(mul)

	return result
}
