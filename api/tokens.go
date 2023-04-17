package api

import (
	"context"
	"database/sql"
	"net/http"

	"github.com/gin-gonic/gin"
)

type TokenResponse struct {
	Address  string `json:"address"`
	Name     string `json:"name"`
	Symbol   string `json:"symbol"`
	Decimals int32  `json:"decimals"`
	Price    string `json:"price,omitempty"`
}

func (r *ginServer) GetAllTokens(ctx *gin.Context) {
	tokens, err := r.store.GetAllTokens(context.Background())
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, err.Error())
			return
		}
		ctx.JSON(http.StatusInternalServerError, err.Error())
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
