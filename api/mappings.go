package api

func (r *ginServer) MapUrls() {
	r.Router.GET("/ping", r.Ping)
	r.Router.GET("/api/tokens", r.GetAllTokens)
	r.Router.GET("/api/pools", r.GetAllPools)
	r.Router.GET("/api/status", r.GetIndexerStatus)
}
