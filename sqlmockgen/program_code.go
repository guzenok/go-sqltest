package main

// imports specialized functions and runs generator.
var srcCodeTemplate = `
package main

import (
	"flag"
	"fmt"
	"os"

	generator {{printf "%q" .GeneratorPath}}
	{{.Pkg.Name}} {{printf "%q" .ImportPath}}
)

var output = flag.String("output", "", "The output file name, or empty to use stdout.")

var copyright = ` + "`{{.CopyrightHeader}}`" + `

func main() {
	flag.Parse()
	
	{{$pkg := .Pkg.Name}}
	inits := generator.InitDataFunctions{
		{{range $f, $_ := .Pkg.Data}}
			{{printf "%q" $f}}: {{$pkg}}.{{$f}},
		{{end}}
	}
	
	sqls := generator.SqlsDictFunctions{
		{{range $f, $_ := .Pkg.Sqls}}
			{{printf "%q" $f}}: {{$pkg}}.{{$f}},
		{{end}}
	}

	out, close := generator.OutWriter(*output)
	defer close()
	
	if copyright!="" {
		fmt.Fprintln(out, copyright)
	}
	
	if err := generator.WriteCode({{printf "%q" .Pkg.Name}}, inits, sqls, out); err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
}
`
