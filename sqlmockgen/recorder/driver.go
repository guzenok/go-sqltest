package recorder

import (
	"database/sql"
	"database/sql/driver"

	"github.com/DATA-DOG/go-sqlmock"
)

// Driver wraps the database driver.
type Driver struct {
	orig driver.Driver
	mock sqlmock.Sqlmock
}

// Open opens a new connection to the database. name is a connection string.
func (d *Driver) Open(name string) (driver.Conn, error) {
	connection, err := d.orig.Open(name)
	if err != nil {
		return nil, err
	}

	return newConn(connection, d.mock), nil
}

// Wrap driver with mock recorder.
func Wrap(orig driver.Driver, mock sqlmock.Sqlmock) (name string, drv driver.Driver) {
	drv = &Driver{
		orig: orig,
		mock: mock,
	}

	name = "aesdfzsdfsd" // TODO: make it random uniq
	sql.Register(name, drv)

	return
}
