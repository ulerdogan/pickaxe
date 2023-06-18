package api

import (
	"context"
	"database/sql"
	"encoding/json"
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

		var prp PoolResponse

		sf := struct {
			Fee0 string `json:"fee_0"`
			Fee1 string `json:"fee_1"`
		}{}

		err = json.Unmarshal([]byte(p.Fee), &sf)
		if err != nil {
			prp = PoolResponse{
				AmmID:    p.AmmID,
				Address:  p.Address,
				TokenA:   p.TokenA,
				TokenB:   p.TokenB,
				ReserveA: p.ReserveA,
				ReserveB: p.ReserveB,
				Fee:      p.Fee,
			}
		} else {
			prp = PoolResponse{
				AmmID:    p.AmmID,
				Address:  p.Address,
				TokenA:   p.TokenA,
				TokenB:   p.TokenB,
				ReserveA: p.ReserveA,
				ReserveB: p.ReserveB,
				Fee:      sf,
			}
		}
		prp.LastUpdated = p.LastUpdated.String()
		prp.LastBlock = p.LastBlock
		if p.TotalValue != "" {
			prp.TotalValue = p.TotalValue
		}
		if p.ExtraData.String != "" {
			prp.ExtraData = p.ExtraData.String
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

	if req.ExtraData != "" {
		r.store.UpdatePoolExtraData(context.Background(), db.UpdatePoolExtraDataParams{
			PoolID:    pool.PoolID,
			ExtraData: sql.NullString{String: req.ExtraData, Valid: true},
		})
	}

	dex, _ := r.client.NewDex(int(pool.AmmID))
	go dex.SyncFee(starknet.PoolInfo{
		Address:   pool.Address,
		ExtraData: pool.ExtraData.String,
	}, r.store, r.client)

	var rsp PoolResponse

	sf := struct {
		Fee0 string `json:"fee_0"`
		Fee1 string `json:"fee_1"`
	}{}

	err = json.Unmarshal([]byte(pool.Fee), &sf)
	if err != nil {
		rsp = PoolResponse{
			AmmID:    pool.AmmID,
			Address:  pool.Address,
			TokenA:   pool.TokenA,
			TokenB:   pool.TokenB,
			ReserveA: pool.ReserveA,
			ReserveB: pool.ReserveB,
			Fee:      pool.Fee,
		}
	} else {
		rsp = PoolResponse{
			AmmID:    pool.AmmID,
			Address:  pool.Address,
			TokenA:   pool.TokenA,
			TokenB:   pool.TokenB,
			ReserveA: pool.ReserveA,
			ReserveB: pool.ReserveB,
			Fee:      sf,
		}
	}

	ctx.JSON(http.StatusOK, rsp)
}
