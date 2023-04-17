package api

import (
	"github.com/gin-gonic/gin"
	db "github.com/ulerdogan/pickaxe/db/sqlc"
)

type Api interface {
	Run(serverAddress string)
	MapUrls()
	SetMiddlewares()
}

type ginServer struct {
	Router *gin.Engine
	store  db.Store
}

func NewRouter(store db.Store) Api {
	return &ginServer{gin.Default(), store}
}

func (r *ginServer) Run(serverAddress string) {
	r.Router.Run(serverAddress)
}
