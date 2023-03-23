package migration

import (
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	logger "github.com/ulerdogan/pickaxe/utils/logger"
)

func RunDBMigration(migrationURL, dbSource string) {
	migration, err := migrate.New(migrationURL, dbSource)
	if err != nil {
		logger.Error(err, "can not create migration")
		return
	}

	if err = migration.Up(); err != nil {
		if err == migrate.ErrNoChange {
			return
		}
		
		logger.Error(err, "can not migrate up")
		return
	}

	logger.Info("db migration is completed")
}
