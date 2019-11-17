package postgres

import (
	"database/sql"
	"fmt"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/golang-migrate/migrate/v4/source/go_bindata"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq" // postgresql driver
	"github.com/pkg/errors"

	store "github.com/guzenok/go-sqltest/sample"
	"github.com/guzenok/go-sqltest/sample/postgres/migrations"
)

//go:generate go run github.com/kevinburke/go-bindata/go-bindata -o ./migrations/migrations.bindata.go -pkg migrations -ignore=\\*.go ./migrations/...

const (
	actualVersion uint = 1550026905
	driverName         = "postgres"
	migrationsDir      = "migrations"
)

type postgresStore struct {
	db *sqlx.DB
}

func (s *postgresStore) Close() error {
	return s.db.Close()
}

func (s *postgresStore) Ping() error {
	return s.db.Ping()
}

// New store over postgres db.
func New(uri string) (store.Store, error) {
	db, err := sql.Open(driverName, uri)
	if err != nil {
		return nil, errors.Wrap(err, "can't connect to database")
	}

	err = Migrate(db)
	if err != nil {
		return nil, errors.Wrap(err, "can't migrate database")
	}

	return wrap(db), nil
}

// wrap db in store.
func wrap(db *sql.DB) *postgresStore {
	return &postgresStore{
		db: sqlx.NewDb(db, driverName),
	}
}

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
