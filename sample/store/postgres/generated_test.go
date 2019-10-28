package postgres

import (
	"database/sql"
	"testing"

	"github.com/pkg/errors"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/guzenok/go-sqltest/sqlmockgen/driver"
)

const (
	uri = "postgresql://postgres:postgres@localhost:5432/test?sslmode=disable"
)

/*
 * This code file should be generated.
 */

func UsersTestDb(t *testing.T) (*sql.DB, sqlmock.Sqlmock, error) {
	db, err := sql.Open(driverName, uri)
	if err != nil {
		return nil, nil, errors.Wrap(err, "can't connect to database")
	}
	defer db.Close()

	_, err = InitDbUsers(db)
	if err != nil {
		return nil, nil, errors.Wrap(err, "can't init database")
	}

	mockDb, mock, err := sqlmock.New()
	if err != nil {
		return nil, nil, err
	}
	mock.ExpectRollback()

	name, _ := driver.Wrap(db.Driver(), mock)
	recorder, err := sql.Open(name, uri)
	if err != nil {
		return nil, nil, errors.Wrap(err, "can't connect to database")
	}

	testStore_Users(t, recorder)

	return mockDb, mock, nil
}
