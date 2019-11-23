package model

import (
	"database/sql"
	"log"
	"testing"

	"github.com/stretchr/testify/assert"
)

// thisPkgPath is the importable path for package model
const thisPkgPath = "github.com/guzenok/go-sqltest/sqlmockgen/model"

func TestBuild(t *testing.T) {
	assert := assert.New(t)

	expect := &Package{
		Name: "model",
		Inits: []string{
			"InitTestDb",
		},
		Tests: []string{
			"Store1MockTest",
			"Store2MockTest",
		},
	}

	_, err := avoidTesting(thisPkgPath)
	defer func() {
		err := restoreTesting(thisPkgPath)
		if err != nil {
			log.Printf("failed to remove *"+tempfile+" files: %s", err)
		}
	}()
	if err != nil {
		t.Fatal(err)
	}

	got, err := Parse(thisPkgPath)
	if !assert.NoError(err) {
		return
	}

	_ = true &&
		assert.Equal(expect.Name, got.Name) &&
		assert.ElementsMatch(expect.Inits, got.Inits) &&
		assert.ElementsMatch(expect.Tests, got.Tests)
}

// InitTestDb is for importer test.
func InitTestDb(db *sql.DB) error {
	return nil
}

// Store1MockTest is for importer test.
func Store1MockTest(t *testing.T, db *sql.DB) {
	return
}

// Store2MockTest is for importer test.
func Store2MockTest(t *testing.T, db *sql.DB) {
	return
}
