package generator

import (
	"bytes"
	"database/sql"
	"database/sql/driver"
	"fmt"
	"go/format"
	"testing"

	"github.com/guzenok/go-sqltest/sqlmockgen/model"
	"github.com/guzenok/go-sqltest/sqlmockgen/recorder"
)

const (
	// ImportPath of current package.
	ImportPath = "github.com/guzenok/go-sqltest/sqlmockgen/generator"
)

var imports = []string{
	"database/sql",
	"database/sql/driver",
	"errors",
	"testing",
	"time",
	"github.com/DATA-DOG/go-sqlmock",
}

type (
	Generator interface {
		GenCode(
			t *testing.T,
			dbUrl string,
			init model.InitDbFunc,
			tests map[string]model.TestDbFunc,
		) []byte
	}

	generator struct {
		buf *bytes.Buffer
	}
)

func New() Generator {
	return &generator{
		buf: new(bytes.Buffer),
	}
}

func (g *generator) GenCode(
	t *testing.T,
	dbUrl string,
	init model.InitDbFunc,
	tests map[string]model.TestDbFunc,
) []byte {
	g.imports()

	db, err := init(dbUrl)
	if err != nil {
		t.Fatal(err)
	}

	for implFunc, f := range tests {
		const pref = "Test"
		testFunc := pref + implFunc[len(pref):]
		mockFunc := implFunc + "Mock"
		g.test(t, testFunc, mockFunc, implFunc)
		g.mock(t, dbUrl, mockFunc, db.Driver(), f)
	}

	src, err := format.Source(g.buf.Bytes())
	if err != nil {
		t.Fatal(err)
	}

	return src
}

func (g *generator) test(t *testing.T, name, mock, impl string) {
	g.p("func %s(t *testing.T) {", name)
	g.p("db, err := %s()", mock)
	g.p("if err != nil {")
	g.p("  t.Fatal(err)")
	g.p("}")
	g.p("%s(t, db)", impl)
	g.p("}")
	g.p("")
}

func (g *generator) mock(t *testing.T,
	dbUrl, name string, driver driver.Driver, f model.TestDbFunc) {
	code := new(bytes.Buffer)
	uid, _ := recorder.Wrap(driver, nil, code)
	rec, err := sql.Open(uid, dbUrl)
	if err != nil {
		t.Fatal(err)
	}
	defer rec.Close()

	f(t, rec)
	if t.Failed() {
		return
	}

	g.p("func %s() (*sql.DB, error) {", name)
	g.p("opt := sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual)")
	g.p("db, mock, err := sqlmock.New(opt)")
	g.p("if err != nil {")
	g.p("  return nil, err")
	g.p("}")
	g.p("")
	g.p(string(code.Bytes()))
	g.p("")
	g.p("  return db, nil")
	g.p("}")
	g.p("")
}

func (g *generator) imports() {
	g.p("import (")
	for _, i := range imports {
		g.p(`"%s"`, i)
	}
	g.p(")")
	g.p("")
}

func (g *generator) p(format string, args ...interface{}) {
	fmt.Fprintf(g.buf, format+"\n", args...)
}
