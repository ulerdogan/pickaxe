package api

import (
	"context"
	"database/sql"
	"math/big"
	"net/http"

	"github.com/NethermindEth/juno/core/felt"
	"github.com/NethermindEth/starknet.go/rpc"
	"github.com/NethermindEth/starknet.go/types"
	"github.com/gin-gonic/gin"

	//auth "github.com/ulerdogan/pickaxe/auth"
	starknet "github.com/ulerdogan/pickaxe/clients/starknet"
	db "github.com/ulerdogan/pickaxe/db/sqlc"
)

func (r *ginServer) GetAllTokens(ctx *gin.Context) {
	tokens, err := r.store.GetAllTokens(context.Background())
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	rsp := make([]TokenResponse, len(tokens))
	for i, t := range tokens {
		rsp[i] = TokenResponse{
			Address:  t.Address,
			Name:     t.Name,
			Symbol:   t.Symbol,
			Decimals: t.Decimals,
			Price:    t.Price,
		}
	}

	ctx.JSON(http.StatusOK, rsp)
}

func (r *ginServer) AddToken(ctx *gin.Context) {
	var req AddTokenParams
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	decimal, err := getTokenDecimal(r.client, req.Address)
	if err != nil {
		d := 18
		decimal = &d
	}

	token, err := r.store.CreateToken(context.Background(), db.CreateTokenParams{
		Address:  req.Address,
		Name:     req.Name,
		Symbol:   req.Symbol,
		Ticker:   req.Ticker,
		Decimals: int32(*decimal),
	})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	rsp := TokenResponse{
		Address:  token.Address,
		Name:     token.Name,
		Symbol:   token.Symbol,
		Decimals: token.Decimals,
	}

	ctx.JSON(http.StatusOK, rsp)
}

func getTokenDecimal(c starknet.Client, address string) (*int, error) {
	paHash := starknet.GetAddressFelt(address)
	r, err := c.Call(rpc.FunctionCall{
		ContractAddress:    paHash,
		EntryPointSelector: types.GetSelectorFromNameFelt("decimals"),
		Calldata:           []*felt.Felt{},
	})

	if err != nil {
		return nil, err
	}

	decimal := int((r[0].BigInt(new(big.Int)).Int64()))
	return &decimal, nil
}
