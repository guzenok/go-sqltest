package recorder

import (
	"database/sql"
	"database/sql/driver"
	"io"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
)

// Driver wraps the database driver.
type Driver struct {
	orig driver.Driver
	mock sqlmock.Sqlmock
	out  io.Writer
}

// Open opens a new connection to the database. name is a connection string.
func (d *Driver) Open(name string) (driver.Conn, error) {
	connection, err := d.orig.Open(name)
	if err != nil {
		return nil, err
	}

	return newConn(connection, d.mock, d.out), nil
}

// Wrap driver with mock recorder.
func Wrap(orig driver.Driver, mock sqlmock.Sqlmock, out io.Writer) (name string, drv driver.Driver) {
	if mock == nil {
		var err error
		_, mock, err = sqlmock.New()
		if err != nil {
			panic(err)
		}
	}

	drv = &Driver{
		orig: orig,
		mock: mock,
		out:  out,
	}

	name = uuid.New().String()
	sql.Register(name, drv)

	return
}
