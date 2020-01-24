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

	"database/sql/driver",
	"time",
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
		imports recorder.ImportList
		code    *bytes.Buffer
	}
)

func New() Generator {
	return &generator{
		imports: make(recorder.ImportList),
		code:    new(bytes.Buffer),
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

	rawsrc := append(g.header(), g.body()...)

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
	g.i("testing")

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
	fcode := new(bytes.Buffer)

	rec, err := recorder.Open(drv, dbUrl, g.imports, fcode, nil)
	if err != nil {
		t.Fatal(err)
	}
	defer rec.Close()

	f(t, rec)

	g.i("database/sql")
	g.i("github.com/DATA-DOG/go-sqlmock")

	g.p("func %s() (*sql.DB, error) {", name)
	g.p("opt := sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual)")
	g.p("db, mock, err := sqlmock.New(opt)")
	g.p("if err != nil {")
	g.p("  return nil, err")
	g.p("}")
	g.p("")
	g.p(string(fcode.Bytes()))
	g.p("")
	g.p("  return db, nil")
	g.p("}")
	g.p("")
}

func (g *generator) header() []byte {
	buf := new(bytes.Buffer)

	buf.WriteString("import (\n")
	for i := range g.imports {
		fmt.Fprintf(buf, "%q\n", i)
	}
	buf.WriteString(")\n")

	return buf.Bytes()
}

func (g *generator) body() []byte {
	return g.code.Bytes()
}

func (g *generator) i(pkg string) {
	g.imports[pkg] = struct{}{}
}

func (g *generator) p(format string, args ...interface{}) {
	fmt.Fprintf(g.code, format+"\n", args...)
}
