package migration

import (
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	logger "github.com/ulerdogan/pickaxe/utils/logger"
)

func RunDBMigration(migrationURL, dbSource string) (bool, error) {
	migration, err := migrate.New(migrationURL, dbSource)
	if err != nil {
		logger.Error(err, "can not create migration")
		return false, err
	}

	if err = migration.Up(); err != nil {
		if err == migrate.ErrNoChange {
			return false, nil
		}

		logger.Error(err, "can not migrate up")
		return false, err
	}

	logger.Info("db migration is completed")
	return true, nil
}
