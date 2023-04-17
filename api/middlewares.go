package api

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func (r *ginServer) SetMiddlewares() {
	r.Router.Use(cors.Default())
	r.Router.Use(gin.Recovery())
}
