package migration

import (
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func RunDBMigration(migrationURL, dbSource string) (bool, error) {
	migration, err := migrate.New(migrationURL, dbSource)
	if err != nil {
		return false, err
	}

	if err = migration.Up(); err != nil {
		if err == migrate.ErrNoChange {
			return false, nil
		}

		return false, err
	}

	return true, nil
}
