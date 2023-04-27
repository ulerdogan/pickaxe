package api

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func (r *ginServer) MapUrls() {
	r.router.GET("/", r.Base)
	r.router.GET("/ping", r.Ping)
	r.router.GET("/api/status", r.GetIndexerStatus)

	r.router.GET("/api/tokens", r.GetAllTokens)
	r.router.GET("/api/pools", r.GetAllPools)
	r.router.GET("/api/amm", r.GetAllAmms)

	r.router.POST("/api/login", r.LoginAdmin)

	authRoutes := r.router.Group("/").Use(authMiddleware(r.token))
	authRoutes.POST("/api/tokens/add", r.AddToken)
	authRoutes.POST("/api/pools/add", r.AddPool)
	authRoutes.POST("/api/amm/add", r.AddAmm)
	authRoutes.POST("/api/status/sync", r.UpdateState)

	r.router.Use(cors.Default())
	r.router.Use(gin.Recovery())
}
