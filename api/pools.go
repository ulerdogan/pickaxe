package api

import (
	"context"
	"database/sql"
	"net/http"

	"github.com/gin-gonic/gin"
)

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

func (r *ginServer) GetAllPools(ctx *gin.Context) {
	pools, err := r.store.GetAllPools(context.Background())
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, err.Error())
			return
		}
		ctx.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	rsp := make([]PoolResponse, len(pools))
	for i, p := range pools {
		prp := PoolResponse{
			Address:    p.Address,
			TokenA:     p.TokenA,
			TokenB:     p.TokenB,
			ReserveA:   p.ReserveA,
			ReserveB:   p.ReserveB,
			Fee:        p.Fee,
			TotalValue: p.TotalValue,
		}
		prp.LastUpdated = p.LastUpdated.String()
		prp.LastBlock = p.LastBlock
		if p.TotalValue != "" {
			prp.TotalValue = p.TotalValue
		}
		rsp[i] = prp
	}

	ctx.JSON(http.StatusOK, rsp)
}
