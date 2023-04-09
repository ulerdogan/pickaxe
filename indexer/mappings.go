package indexer

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (ix *indexer) mapUrls() {
	ix.router.GET("/ping", func(ctx *gin.Context) {
		ctx.String(http.StatusOK, "pickaxe")
	})
}
