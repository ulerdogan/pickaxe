package api

import (
	"context"
	"database/sql"
	"net/http"

	"github.com/gin-gonic/gin"
	db "github.com/ulerdogan/pickaxe/db/sqlc"
)

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

func (r *ginServer) AddPool(ctx *gin.Context) {
	var req AddPoolParams
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, err.Error())
		return
	}

	pool, err := r.store.CreatePool(context.Background(), db.CreatePoolParams{
		Address: req.Address,
		TokenA:  req.TokenA,
		TokenB:  req.TokenB,
		AmmID:   req.AmmId,
		Fee:     req.Fee.String(),
	})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	rsp := PoolResponse{
		Address:  pool.Address,
		TokenA:   pool.TokenA,
		TokenB:   pool.TokenB,
		ReserveA: pool.ReserveA,
		ReserveB: pool.ReserveB,
		Fee:      pool.Fee,
	}

	ctx.JSON(http.StatusOK, rsp)
}
