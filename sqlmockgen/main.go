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
	usageText = `sqlmockgen runs your tests on real db, records sql-traffic into sqlmock
and makes your tests work without real db.
Your test funcs should be: func testTESTNAME(*testing.T, *sql.DB){}.
Your real db init func should be: func initTestDb(dbUrl string) (*sql.DB, error){}.
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
