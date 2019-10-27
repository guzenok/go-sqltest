package postgres

import (
	"database/sql"
	"path/filepath"

	"github.com/pkg/errors"

	"github.com/DATA-DOG/go-sqlmock"
	"gopkg.in/testfixtures.v2"
)

const (
	uri         = "postgresql://postgres:postgres@localhost:5432/test?sslmode=disable"
	fixturesDir = "testdata"
)

func InitDbUsers(db *sql.DB) (err error) {
	err = Migrate(db)
	if err != nil {
		return
	}

	err = loadFixtures(db, "users")
	if err != nil {
		return
	}

	return err
}

func SqlsDictUsers() ([]string, error) {
	return nil, nil
}

func UsersTestDb() (*sql.DB, sqlmock.Sqlmock, error) {
	db, err := sql.Open(driverName, uri)
	if err != nil {
		return nil, nil, errors.Wrap(err, "can't connect to database")
	}

	err = InitDbUsers(db)
	if err != nil {
		return nil, nil, errors.Wrap(err, "can't init database")
	}

	return db, nil, nil

	db, mock, err := sqlmock.New()
	if err != nil {
		return nil, nil, err
	}

	return db, mock, err
}

func loadFixtures(db *sql.DB, name string) (err error) {
	fixtures, err := testfixtures.NewFolder(
		db,
		&testfixtures.PostgreSQL{},
		filepath.Join(fixturesDir, name),
	)
	if err != nil {
		return
	}

	err = fixtures.DetectTestDatabase()
	if err != nil {
		return
	}

	err = fixtures.Load()
	return
}
