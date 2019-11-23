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
	GeneratorPath   string
}

func main() {
	flag.Usage = usage
	flag.Parse()

	if flag.NArg() != 1 {
		usage()
		log.Fatal("Expected exactly one argument")
	}

	descr := &Descr{
		ImportPath:    flag.Arg(0),
		GeneratorPath: generator.ImportPath,
	}

	var err error
	descr.Pkg, err = model.Parse(descr.ImportPath)
	if err != nil {
		log.Fatalf("Failed reading import path: %v", err)
	}

	if *copyrightFile != "" {
		header, err := ioutil.ReadFile(*copyrightFile)
		if err != nil {
			log.Fatalf("Failed reading copyright file: %v", err)
		}
		descr.CopyrightHeader = string(header)
	}

	err = generate(descr, *destination)
	if err != nil {
		log.Fatal(err)
	}

}

func usage() {
	io.WriteString(os.Stderr, usageText)
	flag.PrintDefaults()
}
