package api

import (
	"context"
	"database/sql"
	"net/http"

	"github.com/gin-gonic/gin"
	starknet "github.com/ulerdogan/pickaxe/clients/starknet"
	db "github.com/ulerdogan/pickaxe/db/sqlc"
)

func (r *ginServer) GetAllPools(ctx *gin.Context) {
	pools, err := r.store.GetAllPools(context.Background())
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
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
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	pool, err := r.store.CreatePool(context.Background(), db.CreatePoolParams{
		Address: req.Address,
		TokenA:  req.TokenA,
		TokenB:  req.TokenB,
		AmmID:   req.AmmId,
	})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	dex, _ := r.client.NewDex(int(pool.AmmID))
	go dex.SyncFee(starknet.PoolInfo{
		Address:   pool.Address,
		ExtraData: pool.ExtraData.String,
	}, r.store, r.client)

	rsp := PoolResponse{
		Address:  pool.Address,
		TokenA:   pool.TokenA,
		TokenB:   pool.TokenB,
		ReserveA: pool.ReserveA,
		ReserveB: pool.ReserveB,
	}

	ctx.JSON(http.StatusOK, rsp)
}
