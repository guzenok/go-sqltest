package main

import (
	"flag"
	"io"
	"io/ioutil"
	"log"
	"os"

	"github.com/guzenok/go-sqltest/sqlmockgen/model"
)

const (
	usageText = `sqlmockgen generates mock interfaces by building a program
that uses reflection to understand interfaces. It is enabled
by passing non-flag arguments: an import path.
Example:
	//go:generate sqlmockgen -out=sql_test.go -db=postgresql://postgres:postgres@localhost:5432/test?sslmode=disable .
`
)

var (
	url           = flag.String("db", "", "Real database url.")
	out           = flag.String("out", "", "Output file; defaults to stdout.")
	copyrightFile = flag.String("copyright", "", "Copyright file used to add copyright header.")
)

func main() {
	flag.Usage = usage
	flag.Parse()

	dsc := newDescr()
	dsc.DbUrl = *url
	dsc.OutputPath = *out

	if flag.NArg() != 1 {
		usage()
		log.Fatal("Expected exactly one arguments: import path")
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

func usage() {
	io.WriteString(os.Stderr, usageText)
	flag.PrintDefaults()
}
