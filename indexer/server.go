package indexer

import (
	"database/sql"

	"github.com/gin-gonic/gin"
	starknet "github.com/ulerdogan/pickaxe/clients/starknet"
	"github.com/ulerdogan/pickaxe/db/migration"
	db "github.com/ulerdogan/pickaxe/db/sqlc"
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
		logger.Error(err, "cannot load config for: "+environment)
	}
	logger.Info("config loaded for: " + environment)

	// create db connection
	conn, err := sql.Open(cnfg.DBDriver, cnfg.DBSource)
	if err != nil {
		logger.Error(err, "cannot connect to the db")
	}

	logger.Info("gin server will be running at " + cnfg.ServerAddress)
	// init server
	initServer(conn, cnfg)
}

func initServer(conn *sql.DB, cnfg config.Config) {
	// run db migrations if needed
	migration.RunDBMigration(cnfg.MigrationURL, cnfg.DBSource)

	// setting the indexer
	store := db.NewStore(conn)
	client := starknet.NewStarknetClient(cnfg)
	ix := NewIndexer(store, client, cnfg)

	// setup and run jobs
	setupJobs(ix)
	go ix.scheduler.StartBlocking()

	// run gin server
	router.Run(cnfg.ServerAddress)
}
