package postgres

import (
	"database/sql"
	"database/sql/driver"
	"testing"

	"github.com/pkg/errors"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/guzenok/go-sqltest/sqlmockgen/recorder"
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

	name, _ := recorder.Wrap(db.Driver(), mock)
	rec, err := sql.Open(name, uri)
	if err != nil {
		return nil, nil, errors.Wrap(err, "can't connect to database")
	}

	testStore_Users(t, rec)

	mockDb = rec

	return mockDb, mock, nil
}

func s1() {

	_, mock, _ := sqlmock.New()

	mock.ExpectBegin()
	mock.ExpectQuery(`
INSERT INTO users (id, login, password, is_super) 
VALUES ($1, $2, $3, $4)
RETURNING created_at;`).WithArgs([]driver.NamedValue{driver.NamedValue{Name: "", Ordinal: 1, Value: 1}, driver.NamedValue{Name: "", Ordinal: 2, Value: "user01"}, driver.NamedValue{Name: "", Ordinal: 3, Value: "123456"}, driver.NamedValue{Name: "", Ordinal: 4, Value: false}}).WillReturnError(errors.New("pq: duplicate key value violates unique constraint \"users_pkey\""))
	mock.ExpectRollback()

	//=== RUN   TestStoreMock_Users/delete
	mock.ExpectBegin()
	mock.ExpectExec(`
DELETE FROM users 
WHERE id = $1;`).WithArgs([]driver.NamedValue{driver.NamedValue{Name: "", Ordinal: 1, Value: 1}}).WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectCommit()

	//=== RUN   TestStoreMock_Users/create
	mock.ExpectBegin()
	mock.ExpectQuery(`
INSERT INTO users (id, login, password, is_super) 
VALUES ($1, $2, $3, $4)
RETURNING created_at;`).WithArgs([]driver.NamedValue{driver.NamedValue{Name: "", Ordinal: 1, Value: 1}, driver.NamedValue{Name: "", Ordinal: 2, Value: "user01"}, driver.NamedValue{Name: "", Ordinal: 3, Value: "123456"}, driver.NamedValue{Name: "", Ordinal: 4, Value: false}}).WillReturnRows(sqlmock.NewRows([]string{"created_at"}).FromCSVString("time.Time{wall:0x167f9cf0, ext:63707923573, loc:(*time.Location)(0xc00014e660)}"))
	mock.ExpectCommit()

	//=== RUN   TestStoreMock_Users/set_password
	mock.ExpectBegin()
	mock.ExpectExec(`
UPDATE users SET password = $1 
WHERE id = $2;`).WithArgs([]driver.NamedValue{driver.NamedValue{Name: "", Ordinal: 1, Value: "654321"}, driver.NamedValue{Name: "", Ordinal: 2, Value: 1}}).WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectCommit()

	//=== RUN   TestStoreMock_Users/get_by_login
	mock.ExpectQuery(`
SELECT id, password, created_at, is_super
FROM users
WHERE login = $1
LIMIT 1;`).WithArgs([]driver.NamedValue{driver.NamedValue{Name: "", Ordinal: 1, Value: "user01"}}).WillReturnRows(sqlmock.NewRows([]string{"id", "password", "created_at", "is_super"}).FromCSVString("1, \"654321\", time.Time{wall:0x167f9cf0, ext:63707923573, loc:(*time.Location)(0xc00014e660)}, false"))

	//=== RUN   TestStoreMock_Users/get_by_id
	mock.ExpectQuery(`
SELECT login, password, created_at, is_super
FROM users
WHERE id = $1
LIMIT 1;`).WithArgs([]driver.NamedValue{driver.NamedValue{Name: "", Ordinal: 1, Value: 1}}).WillReturnRows(sqlmock.NewRows([]string{"login", "password", "created_at", "is_super"}).FromCSVString("\"user01\", \"654321\", time.Time{wall:0x167f9cf0, ext:63707923573, loc:(*time.Location)(0xc00014e660)}, false"))

	//=== RUN   TestStoreMock_Users/already_exists#01
	mock.ExpectBegin()
	mock.ExpectQuery(`
INSERT INTO users (id, login, password, is_super) 
VALUES ($1, $2, $3, $4)
RETURNING created_at;`).WithArgs([]driver.NamedValue{driver.NamedValue{Name: "", Ordinal: 1, Value: 1}, driver.NamedValue{Name: "", Ordinal: 2, Value: "user01"}, driver.NamedValue{Name: "", Ordinal: 3, Value: "123456"}, driver.NamedValue{Name: "", Ordinal: 4, Value: false}}).WillReturnError(errors.New("pq: duplicate key value violates unique constraint \"users_pkey\""))
	mock.ExpectRollback()

	//=== RUN   TestStoreMock_Users/delete#01
	mock.ExpectBegin()
	mock.ExpectExec(`
DELETE FROM users 
WHERE id = $1;`).WithArgs([]driver.NamedValue{driver.NamedValue{Name: "", Ordinal: 1, Value: 1}}).WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectCommit()

	//=== RUN   TestStoreMock_Users/create#01
	mock.ExpectBegin()
	mock.ExpectQuery(`
INSERT INTO users (id, login, password, is_super) 
VALUES ($1, $2, $3, $4)
RETURNING created_at;`).WithArgs([]driver.NamedValue{driver.NamedValue{Name: "", Ordinal: 1, Value: 1}, driver.NamedValue{Name: "", Ordinal: 2, Value: "user01"}, driver.NamedValue{Name: "", Ordinal: 3, Value: "123456"}, driver.NamedValue{Name: "", Ordinal: 4, Value: false}}).WillReturnRows(sqlmock.NewRows([]string{"created_at"}).FromCSVString("time.Time{wall:0x17500098, ext:63707923573, loc:(*time.Location)(0xc00014e660)}"))
	mock.ExpectCommit()

	//=== RUN   TestStoreMock_Users/set_password#01
	mock.ExpectBegin()
	mock.ExpectExec(`
UPDATE users SET password = $1 
WHERE id = $2;`).WithArgs([]driver.NamedValue{driver.NamedValue{Name: "", Ordinal: 1, Value: "654321"}, driver.NamedValue{Name: "", Ordinal: 2, Value: 1}}).WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectCommit()

	//=== RUN   TestStoreMock_Users/get_by_login#01
	mock.ExpectQuery(`
SELECT id, password, created_at, is_super
FROM users
WHERE login = $1
LIMIT 1;`).WithArgs([]driver.NamedValue{driver.NamedValue{Name: "", Ordinal: 1, Value: "user01"}}).WillReturnRows(sqlmock.NewRows([]string{"id", "password", "created_at", "is_super"}).FromCSVString("1, \"654321\", time.Time{wall:0x17500098, ext:63707923573, loc:(*time.Location)(0xc00014e660)}, false"))

	//=== RUN   TestStoreMock_Users/get_by_id#01
	mock.ExpectQuery(`
SELECT login, password, created_at, is_super
FROM users
WHERE id = $1
LIMIT 1;`).WithArgs([]driver.NamedValue{driver.NamedValue{Name: "", Ordinal: 1, Value: 1}}).WillReturnRows(sqlmock.NewRows([]string{"login", "password", "created_at", "is_super"}).FromCSVString("\"user01\", \"654321\", time.Time{wall:0x17500098, ext:63707923573, loc:(*time.Location)(0xc00014e660)}, false"))
	mock.ExpectClose()

}
