package starknet_utils

import (
	"math/big"

	"github.com/NethermindEth/juno/core/felt"
	"github.com/shopspring/decimal"
)

func GetDecimal(f *felt.Felt, dc int) decimal.Decimal {
	return toDecimal(f.BigInt(new(big.Int)), dc)
}

func toDecimal(value *big.Int, decimals int) decimal.Decimal {
	mul := decimal.NewFromFloat(float64(10)).Pow(decimal.NewFromFloat(float64(decimals)))
	num, _ := decimal.NewFromString(value.String())
	result := num.Div(mul)

	return result
}
