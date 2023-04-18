package api

func (r *ginServer) MapUrls() {
	r.router.GET("/ping", r.Ping)
	r.router.GET("/api/status", r.GetIndexerStatus)

	r.router.GET("/api/tokens", r.GetAllTokens)
	r.router.GET("/api/pools", r.GetAllPools)
	r.router.GET("/api/amm/add", r.GetAllAmms)

	r.router.POST("/api/tokens/add", r.AddToken)
	r.router.POST("/api/pools/add", r.AddPool)
	r.router.POST("/api/amm/add", r.AddAmm)
}
