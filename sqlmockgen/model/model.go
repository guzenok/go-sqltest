// Package model contains the data model necessary for generating sqlmock implementations.
package model

import (
	"database/sql"
	"go/importer"
	"go/token"
	"go/types"
	"reflect"
	"strings"
	"testing"
)

const (
	thisPackage = "github.com/guzenok/go-sqltest/sqlmockgen/model"
	compiler    = "source"
)

type (
	// Signiture agreement here.

	InitDbFunc func(db *sql.DB) error
	TestDbFunc func(*testing.T, *sql.DB)
)

const (
	// Naming agreement here.

	initDbFuncName   = "InitTestDb"
	testDbFuncSuffix = "Test"
)

var (
	typeofInitDbFunc types.Type
	typeofTestDbFunc types.Type
)

func init() {
	goImporter := importer.ForCompiler(token.NewFileSet(), compiler, nil)
	pkg, err := goImporter.Import(thisPackage)
	if err != nil {
		panic(err)
	}
	scope := pkg.Scope()

	var (
		initFunc InitDbFunc
		initName = reflect.TypeOf(initFunc).Name()
		testFunc TestDbFunc
		testName = reflect.TypeOf(testFunc).Name()
	)
	typeofInitDbFunc = scope.Lookup(initName).Type()
	typeofTestDbFunc = scope.Lookup(testName).Type()
}

func Build(path string) (model *Package, err error) {
	goImporter := importer.ForCompiler(token.NewFileSet(), compiler, nil)
	pkg, err := goImporter.Import(path)
	if err != nil {
		return
	}

	model = &Package{
		Name: pkg.Name(),
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

		// NOTE: AssignableTo() does not properly work
		// because /usr/local/go/src/go/types/predicates.go:286
		// (fixed locally)

		if obj.Name() == initDbFuncName &&
			types.AssignableTo(funcType, typeofInitDbFunc) {
			model.Inits = append(model.Inits, name)
			continue
		}

		if strings.HasSuffix(obj.Name(), testDbFuncSuffix) &&
			types.AssignableTo(funcType, typeofTestDbFunc) {
			model.Tests = append(model.Tests, name)
			continue
		}
	}

	return
}
