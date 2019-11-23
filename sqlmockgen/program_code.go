package main

// imports specialized functions and runs generator.
var srcCodeTemplate = `
package {{.Pkg.Name}}

import (
	"testing"
	"fmt"
	"os"

	generator {{printf "%q" .GeneratorPath}}
)

var copyright = ` + "`{{.CopyrightHeader}}`" + `

func {{.SpecTestName}}(t *testing.T) {
	inits := generator.InitDbFunctions{
		{{range $_, $f := .Pkg.Inits}}
			{{printf "%q" $f}}: {{$f}},
		{{end}}
	}
	
	tests := generator.TestDbFunctions{
		{{range $_, $f := .Pkg.Tests}}
			{{printf "%q" $f}}: {{$f}},
		{{end}}
	}

	out, close := generator.OutWriter({{printf "%q" .OutputPath}})
	defer close()
	
	fmt.Fprintln(out, "package {{.Pkg.Name}}")
	
	if copyright!="" {
		fmt.Fprintln(out, copyright)
	}
	
	if err := generator.WriteCode(t, inits, tests, out); err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
}
`
