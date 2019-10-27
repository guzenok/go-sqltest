package model

import (
	"log"
	"testing"

	"github.com/stretchr/testify/assert"
)

// thisPkgPath is the importable path for package model
const thisPkgPath = "github.com/guzenok/go-sqltest/sqlmockgen/model"

func TestBuild(t *testing.T) {
	assert := assert.New(t)

	_, err := AvoidTesting(thisPkgPath)
	defer func() {
		err := RestoreTesting(thisPkgPath)
		if err != nil {
			log.Printf("failed to remove *"+tempfile+" files: %s", err)
		}
	}()
	if err != nil {
		t.Fatal(err)
	}

	expect := &Package{
		Name: "model",
		Data: map[string]struct{}{
			"InitDbExample1": struct{}{},
		},
		Sqls: map[string]struct{}{
			"SqlsDictExample1": struct{}{},
		},
	}

	got, err := Build(thisPkgPath)
	if !assert.NoError(err) {
		return
	}

	assert.EqualValues(expect, got)
}

// InitDbExample1 is for importer test.
func InitDbExample1(x string, y string) error {
	return nil
}

// SqlsDictExample1 is for importer test.
func SqlsDictExample1() ([]string, error) {
	return []string{
		"select * from table1",
		"select * from table2",
	}, nil
}
