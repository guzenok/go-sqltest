// Package model contains the data model necessary for generating sqlmock implementations.
package model

import (
	"database/sql"
	"errors"
	"go/importer"
	"go/token"
	"go/types"
	"log"
	"reflect"
	"strings"
	"testing"
)

const (
	// ImportPath of current package.
	ImportPath = "github.com/guzenok/go-sqltest/sqlmockgen/model"
	compiler   = "source"
)

// Signiture agreement here.
type (
	InitDbFunc func(dbUrl string) (*sql.DB, error)
	TestDbFunc func(*testing.T, *sql.DB)
)

// Naming agreement here.
const (
	initDbFuncName   = "initTestDb"
	testDbFuncPrefix = "test"
)

func Parse(path string) (model *Package, err error) {
	golang := importer.ForCompiler(token.NewFileSet(), compiler, nil)
	typeofInitDbFunc, typeofTestDbFunc := loadTypes(golang)
	model = &Package{}

	model.SrcDir, err = avoidTesting(path)
	defer func() {
		err := restoreTesting(path)
		if err != nil {
			log.Printf("failed to remove temp-files: %s", err)
		}
	}()
	if err != nil {
		return
	}

	pkg, err := golang.Import(path)
	if err != nil {
		return
	}

	model.Name = pkg.Name()

	scope := pkg.Scope()
	for _, name := range scope.Names() {
		obj := scope.Lookup(name)
		if obj.Exported() {
			continue
		}

		funcType, ok := obj.Type().(*types.Signature)
		if !ok {
			continue
		}

		if obj.Name() == initDbFuncName &&
			types.AssignableTo(funcType, typeofInitDbFunc) {
			model.Init = name
			continue
		}

		if strings.HasPrefix(obj.Name(), testDbFuncPrefix) &&
			types.AssignableTo(funcType, typeofTestDbFunc) {
			model.Tests = append(model.Tests, name)
			continue
		}
	}

	if model.Init == "" {
		return nil, errors.New(initDbFuncName + " function not found")
	}

	if len(model.Tests) < 1 {
		return nil, errors.New(testDbFuncPrefix + "* function not found")
	}

	return
}

func loadTypes(golang types.Importer) (
	typeofInitDbFunc types.Type,
	typeofTestDbFunc types.Type,
) {
	pkg, err := golang.Import(ImportPath)
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
	return
}
