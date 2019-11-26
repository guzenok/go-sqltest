package postgres

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/golang-migrate/migrate/v4/source/go_bindata"

	"github.com/guzenok/go-sqltest/sample/postgres/migrations"
)

//go:generate go run github.com/kevinburke/go-bindata/go-bindata -o ./migrations/bindata.go -pkg migrations -ignore=\\*.go ./migrations/...

// Migrate db to actualVersion.
func Migrate(db *sql.DB) (err error) {
	m, err := migration(db)
	if err != nil {
		return
	}

	version := actualVersion
	if version == 0 {
		err = m.Up()
	} else {
		err = m.Migrate(version)
	}
	if err == migrate.ErrNoChange {
		err = nil
	}
	if err != nil {
		return
	}

	version, dirty, err := m.Version()
	if err != nil {
		return
	}
	if dirty {
		return errors.New("migrations is dirty")
	}

	return
}

func migration(db *sql.DB) (m *migrate.Migrate, err error) {
	names, err := migrations.AssetDir(migrationsDir)
	if err != nil {
		return
	}

	loadByName := func(name string) ([]byte, error) {
		return migrations.Asset(
			fmt.Sprintf("%s/%s", migrationsDir, name))
	}

	assets := bindata.Resource(names, loadByName)

	src, err := bindata.WithInstance(assets)
	if err != nil {
		return
	}

	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		return
	}

	m, err = migrate.NewWithInstance("go-bindata", src, driverName, driver)
	return
}
