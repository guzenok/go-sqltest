package recorder

import (
	"database/sql/driver"
	"io"

	"github.com/DATA-DOG/go-sqlmock"
)

// Driver wraps the database driver.
type Driver struct {
	orig driver.Driver
	rec
}

// newDriver is the wrap-recorder of original driver.
func newDriver(orig driver.Driver, mock sqlmock.Sqlmock, out io.Writer) driver.Driver {
	if mock == nil {
		var err error
		_, mock, err = sqlmock.New()
		if err != nil {
			panic(err)
		}
	}

	return &Driver{
		orig: orig,
		rec:  rec{out, mock},
	}
}

// Open opens a new connection to the database. name is a connection string.
func (d *Driver) Open(name string) (driver.Conn, error) {
	connection, err := d.orig.Open(name)
	if err != nil {
		return nil, err
	}

	return d.newConn(connection), nil
}
