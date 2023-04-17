package api

import (
	"context"
	"database/sql"
	"net/http"

	"github.com/gin-gonic/gin"
)

type AmmResponse struct {
	Name    string `json:"name"`
	Address string `json:"router_address"`
}

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
