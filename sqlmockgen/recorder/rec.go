package recorder

import (
	"database/sql"
	"database/sql/driver"
	"io"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
)

type ImportList map[string]struct{}

func Open(orig driver.Driver, dataSourceName string, imports ImportList, code io.Writer, mock sqlmock.Sqlmock) (*sql.DB, error) {
	drv := newDriver(orig, imports, code, mock)

	name := uuid.New().String()
	sql.Register(name, drv)

	return sql.Open(name, dataSourceName)
}
