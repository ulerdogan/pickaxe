package indexer

import (
	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func (ix *indexer) mapUrls() {
	ix.router.GET("/ping", func(ctx *gin.Context) {
		ctx.String(http.StatusOK, "pickaxe")
	})
}

func (ix *indexer) setMiddlewares() {
	ix.router.Use(cors.Default())
	ix.router.Use(gin.Recovery())
}
