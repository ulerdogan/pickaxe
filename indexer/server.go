package indexer

import (
	"github.com/gin-gonic/gin"
	config "github.com/ulerdogan/pickaxe/utils/config"
	logger "github.com/ulerdogan/pickaxe/utils/logger"
)

var (
	router = gin.Default()
)

func Init(environment string) {
	cnfg, err := config.LoadConfig(environment, ".")
	if err != nil {
		logger.Error(err, "couldn't load config for: "+environment)
	}
	logger.Info("config loaded for: " + environment)

	logger.Info("gin server will be running at " + cnfg.ServerAddress)
	router.Run(cnfg.ServerAddress)
}
