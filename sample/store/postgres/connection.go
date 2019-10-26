package postgres

import (
	"fmt"
	"log"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	bindata "github.com/golang-migrate/migrate/v4/source/go_bindata"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq" // postgresql driver
	"github.com/pkg/errors"

	"github.com/guzenok/go-sqltest/sample/store"
	"github.com/guzenok/go-sqltest/sample/store/postgres/migrations"
)

//go:generate go run github.com/kevinburke/go-bindata/go-bindata -o ./migrations/migrations.bindata.go -pkg migrations -ignore=\\*.go ./migrations/...

const (
	driverName    = "postgres"
	sourceName    = "go-bindata"
	migrationsDir = "migrations"
)

// New constructor for connection to postgresql RDBMS
func New(cfg *store.DatabaseConfig) (store store.Store, err error) {
	if cfg == nil {
		return nil, errors.New("config not set")
	}

	db, err := sqlx.Connect(driverName, cfg.URI)
	if err != nil {
		err = errors.Wrap(err, "can't connect to postgresql database")
		return
	}

	names, err := migrations.AssetDir(migrationsDir)
	if err != nil {
		return
	}

	sourceInstance, err := bindata.WithInstance(bindata.Resource(names, asset))
	if err != nil {
		return
	}

	driver, err := postgres.WithInstance(db.DB, new(postgres.Config))
	if err != nil {
		return
	}

	m, err := migrate.NewWithInstance(sourceName, sourceInstance, driverName, driver)
	if err != nil {
		return
	}

	if cfg.Version == 0 {
		err = m.Up()
	} else {
		err = m.Migrate(cfg.Version)
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
		return nil, errors.New("migrations is dirty")
	}
	log.Printf("migrations is applied: current version %d", version)

	return &connection{db}, nil
}

type connection struct {
	db *sqlx.DB
}

func (c *connection) Close() error {
	return c.db.Close()
}

func (c *connection) Ping() error {
	return c.db.Ping()
}

func asset(name string) ([]byte, error) {
	data, err := migrations.Asset(fmt.Sprintf("%s/%s", migrationsDir, name))
	if err != nil {
		return nil, err
	}
	return data, nil
}
