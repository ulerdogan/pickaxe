package api

import (
	"github.com/gin-gonic/gin"
	starknet "github.com/ulerdogan/pickaxe/clients/starknet"
	db "github.com/ulerdogan/pickaxe/db/sqlc"
)

type Api interface {
	Run(serverAddress string)
	MapUrls()
	SetMiddlewares()
}

type ginServer struct {
	router *gin.Engine
	store  db.Store
	client starknet.Client
}

func NewRouter(store db.Store, client starknet.Client) Api {
	return &ginServer{gin.Default(), store, client}
}

func (r *ginServer) Run(serverAddress string) {
	r.router.Run(serverAddress)
}
