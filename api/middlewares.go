package api

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func (r *ginServer) SetMiddlewares() {
	r.router.Use(cors.Default())
	r.router.Use(gin.Recovery())
}
