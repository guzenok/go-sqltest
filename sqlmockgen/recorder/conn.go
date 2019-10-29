package recorder

import (
	"context"
	"database/sql/driver"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strings"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
)

var (
	// TODO: replace with different types on conn for each set of implemented interfaces
	// by original connection.
	ErrIsNotImplemented = errors.New("is not implemented")
)

// conn wraps the database connection.
type conn struct {
	connection driver.Conn
	txns       []driver.Tx

	code io.Writer
	mock sqlmock.Sqlmock
}

func newConn(c driver.Conn, mock sqlmock.Sqlmock) *conn {
	cn := &conn{
		connection: c,

		code: ioutil.Discard,
		mock: mock,
	}

	cn.code = os.Stderr

	cn.txns = []driver.Tx{cn}

	return cn
}

func (cn *conn) push(c driver.Tx) {
	cn.txns = append([]driver.Tx{c}, cn.txns...)
}

func (cn *conn) pop() driver.Tx {
	tx := cn.txns[0]
	cn.txns = cn.txns[1:]
	return tx
}

// Begin implements driver.Conn.
func (cn *conn) Begin() (driver.Tx, error) {
	after := cn.mock.ExpectBegin()
	cn.write("mock.ExpectBegin()")
	tx, err := cn.connection.Begin()
	if err != nil {
		after.WillReturnError(err)
		cn.write(".WillReturnError(%s)\n", errToString(err))
		return nil, err
	}
	cn.write("\n")

	cn.push(tx)

	return cn, nil
}

// Commit implements driver.Tx.
func (cn *conn) Commit() error {
	tx := cn.pop()
	after := cn.mock.ExpectCommit()
	cn.write("mock.ExpectCommit()")
	err := tx.Commit()
	if err != nil {
		after.WillReturnError(err)
		cn.write(".WillReturnError(%s)", errToString(err))
	}
	cn.write("\n")
	return err
}

// Rollback implements driver.Tx.
func (cn *conn) Rollback() error {
	tx := cn.pop()
	after := cn.mock.ExpectRollback()
	cn.write("mock.ExpectRollback()")
	err := tx.Rollback()
	if err != nil {
		after.WillReturnError(err)
		cn.write(".WillReturnError(%s)", errToString(err))
	}
	cn.write("\n")
	return err
}

// Prepare implements driver.Conn.
func (cn *conn) Prepare(q string) (driver.Stmt, error) {
	after := cn.mock.ExpectPrepare(q)
	cn.write("mock.ExpectPrepare(\"%s\")", q)
	stmt, err := cn.connection.Prepare(q)
	if err != nil {
		after.WillReturnError(err)
		cn.write(".WillReturnError(%s)", errToString(err))
	}
	cn.write("\n")
	return stmt, err
}

// Implement the "ConnBeginTx" interface.
func (cn *conn) BeginTx(ctx context.Context, opts driver.TxOptions) (driver.Tx, error) {
	i, ok := cn.connection.(driver.ConnBeginTx)
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

	cn.push(tx)

	return cn, nil
}

// Close implements driver.Conn.
func (cn *conn) Close() error {
	after := cn.mock.ExpectClose()
	cn.write("mock.ExpectClose()")
	err := cn.connection.Close()
	if err != nil {
		after.WillReturnError(err)
		cn.write(".WillReturnError(%s)", errToString(err))
	}
	cn.write("\n")
	return err
}

// Implement the optional "Execer" interface for one-shot queries.
func (cn *conn) Exec(query string, args []driver.Value) (driver.Result, error) {
	i, ok := cn.connection.(driver.Execer)
	if !ok {
		return nil, ErrIsNotImplemented
	}

	after := cn.mock.ExpectExec(query).WithArgs(args)
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
	i, ok := cn.connection.(driver.ExecerContext)
	if !ok {
		return nil, ErrIsNotImplemented
	}

	after := cn.mock.ExpectExec(query).WithArgs(args)
	cn.write("mock.ExpectExec(`%s`).WithArgs(%s)", query, namedToString(args))
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
	i, ok := cn.connection.(driver.Queryer)
	if !ok {
		return nil, ErrIsNotImplemented
	}

	after := cn.mock.ExpectQuery(query).WithArgs(args)
	cn.write("mock.ExpectQuery(`%s`).WithArgs(%s)", query, argsToString(args))
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
	i, ok := cn.connection.(driver.QueryerContext)
	if !ok {
		return nil, ErrIsNotImplemented
	}

	after := cn.mock.ExpectQuery(query).WithArgs(args)
	cn.write("mock.ExpectQuery(`%s`).WithArgs(%s)", query, namedToString(args))
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

func (cn *conn) write(format string, a ...interface{}) {
	_, err := fmt.Fprintf(cn.code, format, a...)
	if err != nil {
		panic(err)
	}
}

func argsToString(args []driver.Value) string {
	return fmt.Sprintf("%#v", args)
}

func namedToString(args []driver.NamedValue) string {
	return fmt.Sprintf("%#v", args)
}

func rowsToString(rr *rows) string {
	var ss []string
	for _, vv := range rr.vals {
		var s string
		for _, v := range vv {
			switch x := v.(type) {
			case time.Time:
				s = s + fmt.Sprintf("time.Unix(%d, %d), ", x.Unix(), x.Nanosecond())
			default:
				s = s + fmt.Sprintf("%#v, ", v)
			}
		}
		s = fmt.Sprintf("rr.AddRow([]driver.Value{%s}...)", s)
		ss = append(ss, s)
	}

	return fmt.Sprintf(`func() *sqlmock.Rows {
		rr := sqlmock.NewRows(%#v)
		%s
		return rr
	}()`, rr.cols, strings.Join(ss, "\n"))
}

func errToString(err error) string {
	return fmt.Sprintf("errors.New(%#v)", err.Error())
}

func resultToString(res driver.Result) string {
	lastId, err := res.LastInsertId()
	if err != nil {
		// panic(err)
		lastId = 0
	}

	n, err := res.RowsAffected()
	if err != nil {
		panic(err)
	}

	return fmt.Sprintf("sqlmock.NewResult(%d, %d)", lastId, n)
}
