package generator

import (
	"bytes"
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

	db, err := init(dbUrl)
	if err != nil {
		t.Fatal(err)
	}
	drv := db.Driver()
	db.Close()

	g.imports()

	const pref = "Test"
	for impl, f := range tests {
		test := pref + impl[len(pref):]
		mock := impl + "SqlMock"

		ok := t.Run(impl, func(t *testing.T) {
			g.writeTestFunc(test, mock, impl)
			g.writeMockFunc(t, dbUrl, mock, drv, f)
		})
		if !ok {
			break
		}
	}

	rawsrc := g.buf.Bytes()
	if t.Failed() {
		return rawsrc
	}

	src, err := format.Source(rawsrc)
	if err != nil {
		t.Error(err)
		return rawsrc
	}
	return src
}

func (g *generator) writeTestFunc(name, mock, impl string) {
	g.p("func %s(t *testing.T) {", name)
	g.p("db, err := %s()", mock)
	g.p("if err != nil {")
	g.p("  t.Fatal(err)")
	g.p("}")
	g.p("%s(t, db)", impl)
	g.p("}")
	g.p("")
}

func (g *generator) writeMockFunc(t *testing.T, dbUrl, name string, drv driver.Driver, f model.TestDbFunc) {
	code := new(bytes.Buffer)

	rec, err := recorder.Open(drv, dbUrl, nil, code)
	if err != nil {
		t.Fatal(err)
	}
	defer rec.Close()

	f(t, rec)

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
