package postgres

import (
	"database/sql"
	"database/sql/driver"
	"io"
	"io/ioutil"
	"testing"
	"time"

	"github.com/pkg/errors"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/guzenok/go-sqltest/sqlmockgen/recorder"
)

const (
	uri = "postgresql://postgres:postgres@localhost:5432/test?sslmode=disable"
)

/*
 * This code should be generated.
 */

func TestStoreUsers(t *testing.T) {
	opt := sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual)
	db, mock, err := sqlmock.New(opt)
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	setStoreUsersMock(mock)

	testStoreUsers(t, db)
}

func setStoreUsersMock(mock sqlmock.Sqlmock) {
	mock.MatchExpectationsInOrder(true)

	mock.ExpectBegin()
	mock.ExpectQuery(`
INSERT INTO users (id, login, password, is_super) 
VALUES ($1, $2, $3, $4)
RETURNING created_at;`).WithArgs(driver.Value(1),
		driver.Value("user01"),
		driver.Value("first-P"),
		driver.Value(false),
	).WillReturnError(errors.New("pq: duplicate key value violates unique constraint \"users_pkey\""))
	mock.ExpectRollback()

	mock.ExpectBegin()
	mock.ExpectExec(`
DELETE FROM users 
WHERE id = $1;`).WithArgs(
		driver.Value(1),
	).WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectCommit()

	mock.ExpectBegin()
	mock.ExpectQuery(`
INSERT INTO users (id, login, password, is_super) 
VALUES ($1, $2, $3, $4)
RETURNING created_at;`).WithArgs(
		driver.Value(1),
		driver.Value("user01"),
		driver.Value("second-P"),
		driver.Value(false),
	).WillReturnRows(func() *sqlmock.Rows {
		rr := sqlmock.NewRows([]string{"created_at"})
		rr.AddRow([]driver.Value{time.Unix(1572357152, 177755000)}...)
		return rr
	}())
	mock.ExpectCommit()

	mock.ExpectBegin()
	mock.ExpectExec(`
UPDATE users SET password = $1 
WHERE id = $2;`).WithArgs(
		driver.Value("third-P"),
		driver.Value(1),
	).WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectCommit()

	mock.ExpectQuery(`
SELECT id, password, created_at, is_super
FROM users
WHERE login = $1
LIMIT 1;`).WithArgs(driver.Value("user01")).WillReturnRows(func() *sqlmock.Rows {
		rr := sqlmock.NewRows([]string{"id", "password", "created_at", "is_super"})
		rr.AddRow(1, "third-P", time.Unix(1572357152, 177755000), false)
		return rr
	}())

	mock.ExpectQuery(`
SELECT login, password, created_at, is_super
FROM users
WHERE id = $1
LIMIT 1;`).WithArgs(driver.Value(1)).WillReturnRows(func() *sqlmock.Rows {
		rr := sqlmock.NewRows([]string{"login", "password", "created_at", "is_super"})
		rr.AddRow([]driver.Value{"user01", "third-P", time.Unix(1572357152, 177755000), false}...)
		return rr
	}())

}

/*
 * This code should be execute by generator.
 */

func TestGenerator(t *testing.T) {
	getStoreUsersMock(t, ioutil.Discard)
}

func getStoreUsersMock(t *testing.T, out io.Writer) {
	realDb, err := initTestDb(uri)
	if err != nil {
		t.Fatal(err)
	}
	defer realDb.Close()

	uid, _ := recorder.Wrap(realDb.Driver(), nil, out)
	recDb, err := sql.Open(uid, uri)
	if err != nil {
		t.Fatal(err)
	}

	testStoreUsers(t, recDb)
}
