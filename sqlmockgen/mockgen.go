package main

import (
	"flag"
	"io"
	"io/ioutil"
	"log"
	"os"

	"github.com/guzenok/go-sqltest/sqlmockgen/generator"
	"github.com/guzenok/go-sqltest/sqlmockgen/model"
)

const (
	usageText = `sqlmockgen generates mock interfaces by building a program
that uses reflection to understand interfaces. It is enabled
by passing non-flag arguments: an import path.
Example:
	sqlmockgen github.com/company/sql/driver

`
)

var (
	destination   = flag.String("destination", "", "Output file; defaults to stdout.")
	copyrightFile = flag.String("copyright_file", "", "Copyright file used to add copyright header")
)

type Descr struct {
	ImportPath      string
	Pkg             *model.Package
	CopyrightHeader string
	OutputPath      string

	SpecTestName  string
	GeneratorPath string
}

func main() {
	flag.Usage = usage
	flag.Parse()

	dsc := newDescr()
	dsc.OutputPath = *destination

	if flag.NArg() != 1 {
		usage()
		log.Fatal("Expected exactly one arguments")
	}
	dsc.ImportPath = flag.Arg(0)

	var err error
	dsc.Pkg, err = model.Parse(dsc.ImportPath)
	if err != nil {
		log.Fatalf("Failed reading import path: %v", err)
	}

	if *copyrightFile != "" {
		header, err := ioutil.ReadFile(*copyrightFile)
		if err != nil {
			log.Fatalf("Failed reading copyright file: %v", err)
		}
		dsc.CopyrightHeader = string(header)
	}

	err = generate(dsc)
	if err != nil {
		log.Fatal(err)
	}

}

func newDescr() *Descr {
	return &Descr{
		GeneratorPath: generator.ImportPath,
		SpecTestName:  "TestXxx",
	}
}

func usage() {
	io.WriteString(os.Stderr, usageText)
	flag.PrintDefaults()
}
