package recorder

import (
	"context"
	"database/sql/driver"
	"errors"
)

var (
	// TODO: replace with different types on conn for each set of implemented interfaces
	// by original connection.
	ErrIsNotImplemented = errors.New("is not implemented")
)

// conn wraps the database connection.
type conn struct {
	orig driver.Conn
	*output
}

func (d *Driver) newConn(c driver.Conn) *conn {
	return &conn{
		output: d.output,
		orig:   c,
	}
}

// Begin implements driver.Conn.
func (c *conn) Begin() (driver.Tx, error) {
	after := c.mock.ExpectBegin()
	c.p("mock.ExpectBegin()")
	tx, err := c.orig.Begin()
	if err != nil {
		after.WillReturnError(err)
		c.p(".WillReturnError(%s)\n", c.errToString(err))
		return nil, err
	}
	c.p("\n")

	return c.newTx(tx), nil
}

// Prepare implements driver.Conn.
func (c *conn) Prepare(q string) (driver.Stmt, error) {
	after := c.mock.ExpectPrepare(q)
	c.p("mock.ExpectPrepare(\"%s\")", q)
	stmt, err := c.orig.Prepare(q)
	if err != nil {
		after.WillReturnError(err)
		c.p(".WillReturnError(%s)", c.errToString(err))
	}
	c.p("\n")
	return stmt, err
}

// Implement the "ConnBeginTx" interface.
func (c *conn) BeginTx(ctx context.Context, opts driver.TxOptions) (driver.Tx, error) {
	i, ok := c.orig.(driver.ConnBeginTx)
	if !ok {
		return nil, ErrIsNotImplemented
	}

	after := c.mock.ExpectBegin()
	c.p("mock.ExpectBegin()")
	tx, err := i.BeginTx(ctx, opts)
	if err != nil {
		after.WillReturnError(err)
		c.p(".WillReturnError(%s)\n", c.errToString(err))
		return nil, err
	}
	c.p("\n")

	return c.newTx(tx), nil
}

// Close implements driver.Conn.
func (c *conn) Close() error {
	after := c.mock.ExpectClose()
	c.p("mock.ExpectClose()")
	err := c.orig.Close()
	if err != nil {
		after.WillReturnError(err)
		c.p(".WillReturnError(%s)", c.errToString(err))
	}
	c.p("\n")
	return err
}

// Implement the optional "Execer" interface for one-shot queries.
func (c *conn) Exec(query string, args []driver.Value) (driver.Result, error) {
	i, ok := c.orig.(driver.Execer)
	if !ok {
		return nil, ErrIsNotImplemented
	}

	after := c.mock.ExpectExec(query).WithArgs(args...)
	c.p("mock.ExpectExec(`%s`).WithArgs(\n%s)", query, argsToString(args))
	res, err := i.Exec(query, args)
	if err != nil {
		after.WillReturnError(err)
		c.p(".WillReturnError(%s)", c.errToString(err))
	} else {
		after.WillReturnResult(res)
		c.p(".WillReturnResult(%s)", c.resultToString(res))
	}
	c.p("\n")

	return res, err
}

// Implement the "ExecerContext" interface.
func (c *conn) ExecContext(ctx context.Context, query string, args []driver.NamedValue) (driver.Result, error) {
	i, ok := c.orig.(driver.ExecerContext)
	if !ok {
		return nil, ErrIsNotImplemented
	}

	after := c.mock.ExpectExec(query).WithArgs(c.namedToValues(args)...)
	c.p("mock.ExpectExec(`%s`).WithArgs(\n%s)", query, c.namedToString(args))
	res, err := i.ExecContext(ctx, query, args)
	if err != nil {
		after.WillReturnError(err)
		c.p(".WillReturnError(%s)", c.errToString(err))
	} else {
		after.WillReturnResult(res)
		c.p(".WillReturnResult(%s)", c.resultToString(res))
	}
	c.p("\n")

	return res, err
}

// Implement the "Queryer" interface.
func (c *conn) Query(query string, args []driver.Value) (driver.Rows, error) {
	i, ok := c.orig.(driver.Queryer)
	if !ok {
		return nil, ErrIsNotImplemented
	}

	after := c.mock.ExpectQuery(query).WithArgs(args)
	c.p("mock.ExpectQuery(`%s`).WithArgs(\n%s)", query, argsToString(args))
	res, err := i.Query(query, args)
	if err != nil {
		after.WillReturnError(err)
		c.p(".WillReturnError(%s)\n", c.errToString(err))
		return nil, err
	}

	rr := parseRows(res)
	after.WillReturnRows(rr.MockRows())
	c.p(".WillReturnRows(%s)\n", c.rowsToString(rr))

	return rr, err
}

// Implement the "QueryerContext" interface.
func (c *conn) QueryContext(ctx context.Context, query string, args []driver.NamedValue) (driver.Rows, error) {
	i, ok := c.orig.(driver.QueryerContext)
	if !ok {
		return nil, ErrIsNotImplemented
	}

	after := c.mock.ExpectQuery(query).WithArgs(c.namedToValues(args)...)
	c.p("mock.ExpectQuery(`%s`).WithArgs(\n%s)", query, c.namedToString(args))
	res, err := i.QueryContext(ctx, query, args)
	if err != nil {
		after.WillReturnError(err)
		c.p(".WillReturnError(%s)\n", c.errToString(err))
		return nil, err
	}

	rr := parseRows(res)
	after.WillReturnRows(rr.MockRows())
	c.p(".WillReturnRows(%s)\n", c.rowsToString(rr))

	return rr, err
}
