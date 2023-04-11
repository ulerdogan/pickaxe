package indexer

import (
	"database/sql"

	starknet "github.com/ulerdogan/pickaxe/clients/starknet"
	init_db "github.com/ulerdogan/pickaxe/db/init"
	"github.com/ulerdogan/pickaxe/db/migration"
	db "github.com/ulerdogan/pickaxe/db/sqlc"
	config "github.com/ulerdogan/pickaxe/utils/config"
	logger "github.com/ulerdogan/pickaxe/utils/logger"
)

func Init(environment string) {
	// load app configs
	cnfg, err := config.LoadConfig(environment, ".")
	if err != nil {
		logger.Error(err, "cannot load config for: "+environment)
		return
	}
	logger.Info("config loaded for: " + environment)

	// create db connection
	conn, err := sql.Open(cnfg.DBDriver, cnfg.DBSource)
	if err != nil {
		logger.Error(err, "cannot connect to the db")
		return
	}

	logger.Info("gin server will be running at " + cnfg.ServerAddress)
	// init server
	initServer(conn, cnfg)
}

func initServer(conn *sql.DB, cnfg config.Config) {
	// run db migrations if needed
	ok, err := migration.RunDBMigration(cnfg.MigrationURL, cnfg.DBSource)
	if err != nil {
		return
	}

	// setting the indexer
	store := db.NewStore(conn)
	client := starknet.NewStarknetClient(cnfg)
	// adding the initial state to db
	if ok {
		init_db.Init(cnfg, store, client)
	}
	// starting the indexer
	ix := NewIndexer(store, client, cnfg)

	// setup and run jobs
	setupJobs(ix)
	go ix.scheduler.StartBlocking()

	// setup n run gin server
	ix.mapUrls()
	ix.router.Run(cnfg.ServerAddress)
}
