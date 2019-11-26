package postgres

import (
	"database/sql"
	"path/filepath"

	"gopkg.in/testfixtures.v2"
)

const (
	fixturesDir = "testdata"
)

func loadFixtures(db *sql.DB, name string) (err error) {
	fixtures, err := testfixtures.NewFolder(
		db,
		&testfixtures.PostgreSQL{},
		filepath.Join(fixturesDir, name),
	)
	if err != nil {
		return
	}

	err = fixtures.DetectTestDatabase()
	if err != nil {
		return
	}

	err = fixtures.Load()
	return
}
