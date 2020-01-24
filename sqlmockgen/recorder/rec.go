package recorder

import (
	"database/sql"
	"database/sql/driver"
	"fmt"
	"io"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
)

type rec struct {
	out  io.Writer
	mock sqlmock.Sqlmock
}

func (r *rec) write(format string, a ...interface{}) {
	_, err := fmt.Fprintf(r.out, format, a...)
	if err != nil {
		panic(err)
	}
}

func Open(orig driver.Driver, dataSourceName string, mock sqlmock.Sqlmock, out io.Writer) (*sql.DB, error) {
	drv := newDriver(orig, mock, out)

	name := uuid.New().String()
	sql.Register(name, drv)

	return sql.Open(name, dataSourceName)
}
