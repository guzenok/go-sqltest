// Package model contains the data model necessary for generating sqlmock implementations.
package model

import (
	"database/sql"
	"fmt"
	"go/importer"
	"go/token"
	"go/types"
	"io"
	"reflect"
)

const (
	thisPackage = "github.com/guzenok/go-sqltest/sqlmockgen/model"
	compiler    = "source"
)

type (
	Query struct {
		Tx   bool
		SQL  string
		Args []interface{}
	}

	InitDbFunc   func(db *sql.DB) error
	SqlsDictFunc func() []Query

	// Package is a Go package.
	Package struct {
		SrcDir string
		Name   string
		Data   map[string]struct{}
		Sqls   map[string]struct{}
	}
)

var (
	typeofInitDbFunc   types.Type
	typeofSqlsDictFunc types.Type
)

func init() {
	goImporter := importer.ForCompiler(token.NewFileSet(), compiler, nil)
	pkg, err := goImporter.Import(thisPackage)
	if err != nil {
		panic(err)
	}
	scope := pkg.Scope()

	var (
		f1     InitDbFunc
		f1name = reflect.TypeOf(f1).Name()
		f2     SqlsDictFunc
		f2name = reflect.TypeOf(f2).Name()
	)
	typeofInitDbFunc = scope.Lookup(f1name).Type()
	typeofSqlsDictFunc = scope.Lookup(f2name).Type()
}

func Build(path string) (model *Package, err error) {
	goImporter := importer.ForCompiler(token.NewFileSet(), compiler, nil)
	pkg, err := goImporter.Import(path)
	if err != nil {
		return
	}

	model = &Package{
		Name: pkg.Name(),
		Data: make(map[string]struct{}),
		Sqls: make(map[string]struct{}),
	}

	scope := pkg.Scope()
	for _, name := range scope.Names() {
		obj := scope.Lookup(name)
		if !obj.Exported() {
			continue
		}

		funcType, ok := obj.Type().(*types.Signature)
		if !ok {
			continue
		}

		// NOTE: AssignableTo() does not properly work because /usr/local/go/src/go/types/predicates.go:286
		// (fixed locally)

		if types.AssignableTo(funcType, typeofInitDbFunc) {
			model.Data[name] = struct{}{}
			continue
		}

		if types.AssignableTo(funcType, typeofSqlsDictFunc) {
			model.Sqls[name] = struct{}{}
			continue
		}
	}

	return
}

func (pkg *Package) Print(w io.Writer) {
	fmt.Fprintf(w, "package %s\n", pkg.Name)
}
