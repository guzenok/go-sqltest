package postgres

import (
	"database/sql"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq" // postgresql driver
	"github.com/pkg/errors"

	store "github.com/guzenok/go-sqltest/sample"
)

//go:generate go run ../../sqlmockgen -out=sql_test.go -db=postgresql://postgres:postgres@localhost:5432/test?sslmode=disable .

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
