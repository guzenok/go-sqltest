package recorder

import (
	"context"
	"database/sql/driver"
	"errors"
	"fmt"
	"strings"
	"time"
)

var (
	// TODO: replace with different types on conn for each set of implemented interfaces
	// by original connection.
	ErrIsNotImplemented = errors.New("is not implemented")
)

// conn wraps the database connection.
type conn struct {
	orig driver.Conn
	rec
}

func (d *Driver) newConn(c driver.Conn) *conn {
	return &conn{
		rec:  d.rec,
		orig: c,
	}
}

// Begin implements driver.Conn.
func (cn *conn) Begin() (driver.Tx, error) {
	after := cn.mock.ExpectBegin()
	cn.write("mock.ExpectBegin()")
	tx, err := cn.orig.Begin()
	if err != nil {
		after.WillReturnError(err)
		cn.write(".WillReturnError(%s)\n", errToString(err))
		return nil, err
	}
	cn.write("\n")

	return cn.newTx(tx), nil
}

// Prepare implements driver.Conn.
func (cn *conn) Prepare(q string) (driver.Stmt, error) {
	after := cn.mock.ExpectPrepare(q)
	cn.write("mock.ExpectPrepare(\"%s\")", q)
	stmt, err := cn.orig.Prepare(q)
	if err != nil {
		after.WillReturnError(err)
		cn.write(".WillReturnError(%s)", errToString(err))
	}
	cn.write("\n")
	return stmt, err
}

// Implement the "ConnBeginTx" interface.
func (cn *conn) BeginTx(ctx context.Context, opts driver.TxOptions) (driver.Tx, error) {
	i, ok := cn.orig.(driver.ConnBeginTx)
	if !ok {
		return nil, ErrIsNotImplemented
	}

	after := cn.mock.ExpectBegin()
	cn.write("mock.ExpectBegin()")
	tx, err := i.BeginTx(ctx, opts)
	if err != nil {
		after.WillReturnError(err)
		cn.write(".WillReturnError(%s)\n", errToString(err))
		return nil, err
	}
	cn.write("\n")

	return cn.newTx(tx), nil
}

// Close implements driver.Conn.
func (cn *conn) Close() error {
	after := cn.mock.ExpectClose()
	cn.write("mock.ExpectClose()")
	err := cn.orig.Close()
	if err != nil {
		after.WillReturnError(err)
		cn.write(".WillReturnError(%s)", errToString(err))
	}
	cn.write("\n")
	return err
}

// Implement the optional "Execer" interface for one-shot queries.
func (cn *conn) Exec(query string, args []driver.Value) (driver.Result, error) {
	i, ok := cn.orig.(driver.Execer)
	if !ok {
		return nil, ErrIsNotImplemented
	}

	after := cn.mock.ExpectExec(query).WithArgs(args...)
	cn.write("mock.ExpectExec(`%s`).WithArgs(\n%s)", query, argsToString(args))
	res, err := i.Exec(query, args)
	if err != nil {
		after.WillReturnError(err)
		cn.write(".WillReturnError(%s)", errToString(err))
	} else {
		after.WillReturnResult(res)
		cn.write(".WillReturnResult(%s)", resultToString(res))
	}
	cn.write("\n")

	return res, err
}

// Implement the "ExecerContext" interface.
func (cn *conn) ExecContext(ctx context.Context, query string, args []driver.NamedValue) (driver.Result, error) {
	i, ok := cn.orig.(driver.ExecerContext)
	if !ok {
		return nil, ErrIsNotImplemented
	}

	after := cn.mock.ExpectExec(query).WithArgs(namedToValues(args)...)
	cn.write("mock.ExpectExec(`%s`).WithArgs(\n%s)", query, namedToString(args))
	res, err := i.ExecContext(ctx, query, args)
	if err != nil {
		after.WillReturnError(err)
		cn.write(".WillReturnError(%s)", errToString(err))
	} else {
		after.WillReturnResult(res)
		cn.write(".WillReturnResult(%s)", resultToString(res))
	}
	cn.write("\n")

	return res, err
}

// Implement the "Queryer" interface.
func (cn *conn) Query(query string, args []driver.Value) (driver.Rows, error) {
	i, ok := cn.orig.(driver.Queryer)
	if !ok {
		return nil, ErrIsNotImplemented
	}

	after := cn.mock.ExpectQuery(query).WithArgs(args)
	cn.write("mock.ExpectQuery(`%s`).WithArgs(\n%s)", query, argsToString(args))
	res, err := i.Query(query, args)
	if err != nil {
		after.WillReturnError(err)
		cn.write(".WillReturnError(%s)\n", errToString(err))
		return nil, err
	}

	rr := parseRows(res)
	after.WillReturnRows(rr.MockRows())
	cn.write(".WillReturnRows(%s)\n", rowsToString(rr))

	return rr, err
}

// Implement the "QueryerContext" interface.
func (cn *conn) QueryContext(ctx context.Context, query string, args []driver.NamedValue) (driver.Rows, error) {
	i, ok := cn.orig.(driver.QueryerContext)
	if !ok {
		return nil, ErrIsNotImplemented
	}

	after := cn.mock.ExpectQuery(query).WithArgs(namedToValues(args)...)
	cn.write("mock.ExpectQuery(`%s`).WithArgs(\n%s)", query, namedToString(args))
	res, err := i.QueryContext(ctx, query, args)
	if err != nil {
		after.WillReturnError(err)
		cn.write(".WillReturnError(%s)\n", errToString(err))
		return nil, err
	}

	rr := parseRows(res)
	after.WillReturnRows(rr.MockRows())
	cn.write(".WillReturnRows(%s)\n", rowsToString(rr))

	return rr, err
}

func argsToString(args []driver.Value) string {
	return fmt.Sprintf("%#v", args)
}

func namedToValues(args []driver.NamedValue) []driver.Value {
	vv := make([]driver.Value, len(args), len(args))
	for i, nv := range args {
		vv[i] = nv.Value
	}
	return vv
}

func namedToString(args []driver.NamedValue) string {
	var s string
	for _, a := range args {
		s = s + fmt.Sprintf("driver.Value(%s),\n", valToString(a.Value))
	}
	return s
}

func rowsToString(rr *rows) string {
	var ss []string
	for _, vv := range rr.vals {
		var s string
		for _, v := range vv {
			s = s + valToString(v) + ", "
		}
		s = fmt.Sprintf("rr.AddRow(%s)", s)
		ss = append(ss, s)
	}

	return fmt.Sprintf(`func() *sqlmock.Rows {
		rr := sqlmock.NewRows(%#v)
		%s
		return rr
	}()`, rr.cols, strings.Join(ss, "\n"))
}

func valToString(v interface{}) string {
	if v == nil {
		return "nil"
	}

	switch x := v.(type) {
	case time.Time:
		return fmt.Sprintf("time.Unix(%d, %d)", x.Unix(), x.Nanosecond())
	default:
		return fmt.Sprintf("%#v", v)
	}
}

func errToString(err error) string {
	return fmt.Sprintf("errors.New(%#v)", err.Error())
}

func resultToString(res driver.Result) string {
	lastId, err := res.LastInsertId()
	if err != nil {
		lastId = 0
	}

	n, err := res.RowsAffected()
	if err != nil {
		panic(err)
	}

	return fmt.Sprintf("sqlmock.NewResult(%d, %d)", lastId, n)
}
