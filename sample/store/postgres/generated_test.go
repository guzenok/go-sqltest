package postgres

import (
	"database/sql"
	"database/sql/driver"
	"testing"

	"time"

	"github.com/pkg/errors"

	"github.com/DATA-DOG/go-sqlmock"
	//"github.com/guzenok/go-sqltest/sqlmockgen/recorder"
)

const (
	uri = "postgresql://postgres:postgres@localhost:5432/test?sslmode=disable"
)

/*
 * This code file should be generated.
 */

func UsersTestDb(t *testing.T) (*sql.DB, sqlmock.Sqlmock, error) {
	/*
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

		name, _ := recorder.Wrap(db.Driver(), mock)
		rec, err := sql.Open(name, uri)
		if err != nil {
			return nil, nil, errors.Wrap(err, "can't connect to database")
		}

		testStore_Users(t, rec)
	*/
	mockDb, mock, err := sqlmock.New()
	if err != nil {
		return nil, nil, err
	}
	generate(mock)

	return mockDb, mock, nil
}

func generate(mock sqlmock.Sqlmock) {
	// === RUN   TestStoreMock_Users/already_exists
	mock.ExpectBegin()
	//mock.ExpectRollback() // ????
	mock.ExpectQuery(`
INSERT INTO users (id, login, password, is_super) 
VALUES ($1, $2, $3, $4)
RETURNING created_at;`).WithArgs([]driver.NamedValue{
		driver.NamedValue{Name: "", Ordinal: 1, Value: 1},
		driver.NamedValue{Name: "", Ordinal: 2, Value: "user01"},
		driver.NamedValue{Name: "", Ordinal: 3, Value: "123456"},
		driver.NamedValue{Name: "", Ordinal: 4, Value: false},
	}).
		WillReturnError(errors.New("pq: duplicate key value violates unique constraint \"users_pkey\""))
	mock.ExpectRollback()

	// === RUN   TestStoreMock_Users/delete
	mock.ExpectBegin()
	mock.ExpectExec(`
DELETE FROM users 
WHERE id = $1;`).
		WithArgs([]driver.NamedValue{driver.NamedValue{Name: "", Ordinal: 1, Value: 1}}).
		WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectCommit()

	// === RUN   TestStoreMock_Users/create
	mock.ExpectBegin()
	mock.ExpectQuery(`
INSERT INTO users (id, login, password, is_super) 
VALUES ($1, $2, $3, $4)
RETURNING created_at;`).
		WithArgs([]driver.NamedValue{
			driver.NamedValue{Name: "", Ordinal: 1, Value: 1},
			driver.NamedValue{Name: "", Ordinal: 2, Value: "user01"},
			driver.NamedValue{Name: "", Ordinal: 3, Value: "123456"},
			driver.NamedValue{Name: "", Ordinal: 4, Value: false},
		}).
		WillReturnRows(func() *sqlmock.Rows {
			rr := sqlmock.NewRows([]string{"created_at"})
			rr.AddRow([]driver.Value{time.Unix(15, 775399000)}...)
			return rr
		}())
	mock.ExpectCommit()

	// === RUN   TestStoreMock_Users/set_password
	mock.ExpectBegin()
	mock.ExpectExec(`
UPDATE users SET password = $1 
WHERE id = $2;`).WithArgs([]driver.NamedValue{driver.NamedValue{Name: "", Ordinal: 1, Value: "654321"}, driver.NamedValue{Name: "", Ordinal: 2, Value: 1}}).WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectCommit()

	// === RUN   TestStoreMock_Users/get_by_login
	mock.ExpectQuery(`
SELECT id, password, created_at, is_super
FROM users
WHERE login = $1
LIMIT 1;`).WithArgs([]driver.NamedValue{driver.NamedValue{Name: "", Ordinal: 1, Value: "user01"}}).WillReturnRows(func() *sqlmock.Rows {
		rr := sqlmock.NewRows([]string{"id", "password", "created_at", "is_super"})
		rr.AddRow([]driver.Value{1, "654321", time.Unix(15, 775399000), false}...)
		return rr
	}())

	// === RUN   TestStoreMock_Users/get_by_id
	mock.ExpectQuery(`
SELECT login, password, created_at, is_super
FROM users
WHERE id = $1
LIMIT 1;`).WithArgs([]driver.NamedValue{driver.NamedValue{Name: "", Ordinal: 1, Value: 1}}).WillReturnRows(func() *sqlmock.Rows {
		rr := sqlmock.NewRows([]string{"login", "password", "created_at", "is_super"})
		rr.AddRow([]driver.Value{"user01", "654321", time.Unix(15, 775399000), false}...)
		return rr
	}())
}
