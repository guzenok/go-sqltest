package recorder

import (
	"database/sql/driver"
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
)

type output struct {
	imports ImportList
	code    io.Writer
	mock    sqlmock.Sqlmock
}

func newOutput(imports ImportList, code io.Writer, mock sqlmock.Sqlmock) *output {
	if mock == nil {
		var err error
		_, mock, err = sqlmock.New()
		if err != nil {
			panic(err)
		}
	}

	return &output{
		imports: imports,
		code:    code,
		mock:    mock,
	}
}

func (r *output) i(pkg string) {
	r.imports[pkg] = struct{}{}
}

func (r *output) p(format string, a ...interface{}) {
	_, err := fmt.Fprintf(r.code, format, a...)
	if err != nil {
		panic(err)
	}
}

func argsToString(args []driver.Value) string {
	return fmt.Sprintf("%#v", args)
}

func (out *output) namedToValues(args []driver.NamedValue) []driver.Value {
	vv := make([]driver.Value, len(args), len(args))
	for i, nv := range args {
		vv[i] = nv.Value
	}
	return vv
}

func (out *output) namedToString(args []driver.NamedValue) string {
	var s string
	for _, a := range args {
		out.i("database/sql/driver")
		s = s + fmt.Sprintf("driver.Value(%s),\n", out.valToString(a.Value))
	}
	return s
}

func (out *output) rowsToString(rr *rows) string {
	var ss []string
	for _, vv := range rr.vals {
		var s string
		for _, v := range vv {
			s = s + out.valToString(v) + ", "
		}
		s = fmt.Sprintf("rr.AddRow(%s)", s)
		ss = append(ss, s)
	}

	out.i("github.com/DATA-DOG/go-sqlmock")
	return fmt.Sprintf(`func() *sqlmock.Rows {
		rr := sqlmock.NewRows(%#v)
		%s
		return rr
	}()`, rr.cols, strings.Join(ss, "\n"))
}

func (out *output) valToString(v interface{}) string {
	if v == nil {
		return "nil"
	}

	switch x := v.(type) {
	case time.Time:
		out.i("time")
		return fmt.Sprintf("time.Unix(%d, %d)", x.Unix(), x.Nanosecond())
	default:
		return fmt.Sprintf("%#v", v)
	}
}

func (out *output) errToString(err error) string {
	out.i("errors")
	return fmt.Sprintf("errors.New(%#v)", err.Error())
}

func (out *output) resultToString(res driver.Result) string {
	lastId, err := res.LastInsertId()
	if err != nil {
		lastId = 0
	}

	n, err := res.RowsAffected()
	if err != nil {
		panic(err)
	}

	out.i("github.com/DATA-DOG/go-sqlmock")
	return fmt.Sprintf("sqlmock.NewResult(%d, %d)", lastId, n)
}
