package postgres

import (
	"database/sql"
	"database/sql/driver"
	"fmt"
	"reflect"

	"github.com/pkg/errors"

	"github.com/DATA-DOG/go-sqlmock"
)

const (
	uri = "postgresql://postgres:postgres@localhost:5432/test?sslmode=disable"
)

/*
 * This code file should be generated.
 */

func UsersTestDb() (*sql.DB, sqlmock.Sqlmock, error) {
	db, err := sql.Open(driverName, uri)
	if err != nil {
		return nil, nil, errors.Wrap(err, "can't connect to database")
	}

	ctx, err := InitDbUsers(db)
	if err != nil {
		return nil, nil, errors.Wrap(err, "can't init database")
	}

	mockdb, mock, err := sqlmock.New()
	if err != nil {
		return nil, nil, err
	}

	for _, q := range SqlsDictUsers(ctx) {
		if q.Tx {
			mock.ExpectBegin()
			fmt.Println(">>>>>>>>> mock.ExpectBegin()")
		}

		vals := make([]driver.Value, len(q.Args), len(q.Args))
		for i, v := range q.Args {
			vals[i] = v
		}
		after := mock.ExpectQuery(q.SQL).WithArgs(vals...)
		fmt.Print(">>>>>>>>> mock.ExpectQuery(q.SQL).WithArgs(vals...)")
		rows, err := db.Query(q.SQL, q.Args...)
		if err != nil {
			after.WillReturnError(err)
			fmt.Println(".WillReturnError(err)")
			if q.Tx {
				mock.ExpectRollback()
				fmt.Println(">>>>>>>>> mock.ExpectRollback()")
			}
			continue
		}

		if q.Tx {
			mock.ExpectCommit()
			fmt.Println(">>>>>>>>> mock.ExpectCommit()")
		}

		cc, err := rows.Columns()
		if err != nil {
			panic(err)
		}
		rr := sqlmock.NewRows(cc)

		tt, err := rows.ColumnTypes()
		if err != nil {
			panic(err)
		}
		dst0 := make([]interface{}, len(tt), len(tt))
		for i, t := range tt {
			dst0[i] = reflect.New(t.ScanType()).Interface()
		}

		for rows.Next() {
			err = rows.Scan(dst0...)
			if err != nil {
				panic(err)
			}

			dst1 := make([]driver.Value, len(cc), len(cc))
			for i, v := range dst0 {
				dst1[i] = v
			}
			rr.AddRow(dst1...)
		}

		after.WillReturnRows(rr)

	}

	return mockdb, mock, nil
}
