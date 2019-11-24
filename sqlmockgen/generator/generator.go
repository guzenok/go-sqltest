package generator

import (
	"bytes"
	"fmt"
	"go/format"
	"testing"

	"github.com/guzenok/go-sqltest/sqlmockgen/model"
)

const ImportPath = "github.com/guzenok/go-sqltest/sqlmockgen/generator"

var imports = []string{
	"testing",
}

type (
	InitDbFunctions map[string]model.InitDbFunc
	TestDbFunctions map[string]model.TestDbFunc

	Generator interface {
		GenCode(*testing.T, InitDbFunctions, TestDbFunctions) ([]byte, error)
	}

	generator struct {
		buf    *bytes.Buffer
		indent string
	}
)

func New() Generator {
	return &generator{
		buf: new(bytes.Buffer),
	}
}

func (g *generator) GenCode(
	t *testing.T,
	inits InitDbFunctions,
	tests TestDbFunctions,
) (
	[]byte, error,
) {

	g.p("import (")
	g.in()
	for _, i := range imports {
		g.p(`"%s"`, i)
	}
	g.out()
	g.p(")")

	return format.Source(g.buf.Bytes())
}

func (g *generator) p(format string, args ...interface{}) {
	fmt.Fprintf(g.buf, g.indent+format+"\n", args...)
}

func (g *generator) in() {
	g.indent += "\t"
}

func (g *generator) out() {
	if len(g.indent) > 0 {
		g.indent = g.indent[0 : len(g.indent)-1]
	}
}
