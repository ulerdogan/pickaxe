package api

import (
	"context"
	"database/sql"
	"net/http"

	"github.com/gin-gonic/gin"
	db "github.com/ulerdogan/pickaxe/db/sqlc"
)

func (r *ginServer) GetAllAmms(ctx *gin.Context) {
	amms, err := r.store.GetAllAmms(context.Background())
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, err.Error())
			return
		}
		ctx.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	rsp := make([]AmmResponse, len(amms))
	for i, a := range amms {
		rsp[i] = AmmResponse{
			Name:    a.DexName,
			Address: a.RouterAddress,
		}
	}

	ctx.JSON(http.StatusOK, rsp)
}

func (r *ginServer) AddAmm(ctx *gin.Context) {
	var req AddAmmParams
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, err.Error())
		return
	}

	amm, err := r.store.CreateAmm(context.Background(), db.CreateAmmParams{
		DexName:       req.Name,
		RouterAddress: req.Address,
		Key:           req.Key,
		AlgorithmType: req.AlgorithmType,
	})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	rsp := AmmResponse{
		Address: amm.RouterAddress,
		Name:    amm.DexName,
	}

	ctx.JSON(http.StatusOK, rsp)
}
