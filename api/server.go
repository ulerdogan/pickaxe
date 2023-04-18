package api

import (
	"github.com/gin-gonic/gin"
	auth "github.com/ulerdogan/pickaxe/auth"
	starknet "github.com/ulerdogan/pickaxe/clients/starknet"
	db "github.com/ulerdogan/pickaxe/db/sqlc"
	config "github.com/ulerdogan/pickaxe/utils/config"
)

type Api interface {
	Run()
	MapUrls()
}

type ginServer struct {
	router *gin.Engine
	store  db.Store
	client starknet.Client
	token  auth.Maker
	config config.Config
}

func NewRouter(store db.Store, client starknet.Client, maker auth.Maker, cnfg config.Config) Api {
	return &ginServer{
		router: gin.Default(),
		store:  store,
		client: client,
		token:  maker,
		config: cnfg,
	}
}

func (r *ginServer) Run() {
	r.router.Run(r.config.ServerAddress)
}

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}
