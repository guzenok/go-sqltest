package recorder

import (
	"database/sql/driver"
	"fmt"
	"io"
	"strings"

	"github.com/DATA-DOG/go-sqlmock"
)

type rows struct {
	cols []string
	vals [][]driver.Value
}

func parseRows(src driver.Rows) *rows {
	dst := &rows{
		cols: src.Columns(),
	}

	n := len(dst.cols)
	for {
		vv := make([]driver.Value, n, n)
		if err := src.Next(vv); err != nil {
			break
		}
		dst.vals = append(dst.vals, vv)
	}

	return dst
}

// MockRows from rows.
func (rr *rows) MockRows() *sqlmock.Rows {
	mock := sqlmock.NewRows(rr.cols)
	for _, vv := range rr.vals {
		mock.AddRow(vv...)
	}
	return mock
}

// CSV of go-values of all rows.
func (rr *rows) CSV() string {
	var buf []string
	for _, vv := range rr.vals {
		for _, v := range vv {
			buf = append(buf, fmt.Sprintf("%v", v))
		}
	}
	return strings.Join(buf, ", ")
}

// Columns returns the names of the columns. The number of
// columns of the result is inferred from the length of the
// slice. If a particular column name isn't known, an empty
// string should be returned for that entry.
func (rr *rows) Columns() []string {
	return rr.cols
}

// Close closes the rows iterator.
func (rr *rows) Close() error {
	rr.vals = nil
	return nil
}

// Next is called to populate the next row of data into
// the provided slice. The provided slice will be the same
// size as the Columns() are wide.
//
// Next should return io.EOF when there are no more rows.
//
// The dest should not be written to outside of Next. Care
// should be taken when closing Rows not to modify
// a buffer held in dest.
func (rr *rows) Next(dst []driver.Value) error {
	if len(rr.vals) == 0 {
		return io.EOF
	}

	copy(dst, rr.vals[0])
	rr.vals = rr.vals[1:]

	return nil
}
