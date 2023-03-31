package indexer

import (
	"database/sql"

	"github.com/gin-gonic/gin"
	"github.com/ulerdogan/pickaxe/db/migration"
	config "github.com/ulerdogan/pickaxe/utils/config"
	logger "github.com/ulerdogan/pickaxe/utils/logger"
)

var (
	router = gin.Default()
)

func Init(environment string) {
	// load app configs
	cnfg, err := config.LoadConfig(environment, ".")
	if err != nil {
		logger.Error(err, "couldn't load config for: "+environment)
	}
	logger.Info("config loaded for: " + environment)

	// create db connection
	conn, err := sql.Open(cnfg.DBDriver, cnfg.DBSource)
	if err != nil {
		logger.Error(err, "can not connect to the db")
	}

	logger.Info("gin server will be running at " + cnfg.ServerAddress)
	// init server
	initServer(conn, cnfg)
}

func initServer(conn *sql.DB, cnfg config.Config) {
	// run db migrations if needed
	migration.RunDBMigration(cnfg.MigrationURL, cnfg.DBSource)
	// run gin server
	router.Run(cnfg.ServerAddress)
}
