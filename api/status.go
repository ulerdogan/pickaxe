package api

import (
	"context"
	"database/sql"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (r *ginServer) Base(ctx *gin.Context) {
	ctx.String(http.StatusOK, "pickaxe")
}

func (r *ginServer) Ping(ctx *gin.Context) {
	ctx.String(http.StatusOK, "pong")
}

type IndexerStatusResponse struct {
	LastBlock   int64  `json:"last_block"`
	LastUpdated string `json:"last_updated"`
}

func (r *ginServer) GetIndexerStatus(ctx *gin.Context) {
	status, err := r.store.GetIndexerStatus(context.Background())
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	if !status.LastQueried.Valid {
		ctx.JSON(http.StatusNotFound, errorResponse(err))
		return
	}

	rsp := IndexerStatusResponse{
		LastBlock:   status.LastQueried.Int64,
		LastUpdated: status.LastUpdated.Time.String(),
	}

	ctx.JSON(http.StatusOK, rsp)
}
