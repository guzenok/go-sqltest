package model

import (
	"database/sql"
	"log"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParse(t *testing.T) {
	assert := assert.New(t)

	expect := &Package{
		Name: "model",
		Init: "initTestDb",
		Tests: []string{
			"testStore1Mock",
			"testStore2Mock",
		},
	}

	_, err := avoidTesting(ImportPath)
	defer func() {
		err := restoreTesting(ImportPath)
		if err != nil {
			log.Printf("failed to remove *"+tempfile+" files: %s", err)
		}
	}()
	if err != nil {
		t.Fatal(err)
	}

	got, err := Parse(ImportPath)
	if !assert.NoError(err) {
		return
	}

	_ = true &&
		assert.Equal(expect.Name, got.Name) &&
		assert.Equal(expect.Init, got.Init) &&
		assert.ElementsMatch(expect.Tests, got.Tests)
}

// initTestDb is for importer test.
func initTestDb(dbUrl string) (*sql.DB, error) {
	return nil, nil
}

// testStore1Mock is for importer test.
func testStore1Mock(t *testing.T, db *sql.DB) {
	return
}

// testStore2Mock is for importer test.
func testStore2Mock(t *testing.T, db *sql.DB) {
	return
}
