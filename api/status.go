package api

import (
	"context"
	"database/sql"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (r *ginServer) Ping(ctx *gin.Context) {
	ctx.String(http.StatusOK, "pickaxe")
}

func (r *ginServer) GetIndexerStatus(ctx *gin.Context) {
	status, err := r.store.GetIndexerStatus(context.Background())
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, err.Error())
			return
		}
		ctx.JSON(http.StatusInternalServerError, err.Error())
		return
	}
	
	ctx.JSON(http.StatusOK, status)
}
