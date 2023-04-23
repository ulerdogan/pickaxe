// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.17.2

package db

import (
	"context"
)

type Querier interface {
	CreateAmm(ctx context.Context, arg CreateAmmParams) (Amm, error)
	CreatePool(ctx context.Context, arg CreatePoolParams) (PoolsV2, error)
	CreateToken(ctx context.Context, arg CreateTokenParams) (Token, error)
	DeleteAmm(ctx context.Context, ammID int64) error
	DeletePool(ctx context.Context, address string) error
	DeleteToken(ctx context.Context, address string) error
	GetAllAmms(ctx context.Context) ([]Amm, error)
	GetAllPools(ctx context.Context) ([]PoolsV2, error)
	GetAllTokens(ctx context.Context) ([]Token, error)
	GetAllTokensWithTickers(ctx context.Context) ([]Token, error)
	GetAmmByDEX(ctx context.Context, dexName string) ([]Amm, error)
	GetAmmById(ctx context.Context, ammID int64) (Amm, error)
	GetAmmKeys(ctx context.Context) ([]string, error)
	GetBaseTokens(ctx context.Context) ([]Token, error)
	GetHashedIndexerPwd(ctx context.Context) (string, error)
	GetIndexerStatus(ctx context.Context) (Indexer, error)
	GetKeys(ctx context.Context) ([]string, error)
	GetNativeTokens(ctx context.Context) ([]Token, error)
	GetPoolByAddress(ctx context.Context, address string) (PoolsV2, error)
	GetPoolByAddressExtra(ctx context.Context, arg GetPoolByAddressExtraParams) (PoolsV2, error)
	GetPoolsByAmm(ctx context.Context, ammID int64) ([]PoolsV2, error)
	GetPoolsByPair(ctx context.Context, arg GetPoolsByPairParams) ([]PoolsV2, error)
	GetPoolsByToken(ctx context.Context, tokenA string) ([]PoolsV2, error)
	GetTokenAPriceByPool(ctx context.Context, poolID int64) (string, error)
	GetTokenBPriceByPool(ctx context.Context, poolID int64) (string, error)
	GetTokenByAddress(ctx context.Context, address string) (Token, error)
	GetTokenBySymbol(ctx context.Context, symbol string) (Token, error)
	InitIndexer(ctx context.Context, arg InitIndexerParams) (Indexer, error)
	UpdateBaseNativeStatus(ctx context.Context, arg UpdateBaseNativeStatusParams) (Token, error)
	UpdateIndexerStatus(ctx context.Context, arg UpdateIndexerStatusParams) (Indexer, error)
	UpdatePoolExtraData(ctx context.Context, arg UpdatePoolExtraDataParams) (PoolsV2, error)
	UpdatePoolReserves(ctx context.Context, arg UpdatePoolReservesParams) (PoolsV2, error)
	UpdatePoolTV(ctx context.Context, arg UpdatePoolTVParams) (PoolsV2, error)
	UpdatePrice(ctx context.Context, arg UpdatePriceParams) (Token, error)
	UpdateTicker(ctx context.Context, arg UpdateTickerParams) (Token, error)
}

var _ Querier = (*Queries)(nil)
