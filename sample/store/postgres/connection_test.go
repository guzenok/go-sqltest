package postgres

import (
	"database/sql"
	"testing"

	"gopkg.in/testfixtures.v2"

	"github.com/guzenok/go-sqltest/sample/store"
)

const (
	ErrForValueExpectedGot = "for %v: expected %v, got %v"
)

func testConnection(t *testing.T) store.Store {
	// Create new sqlx db connection to apply latest migrations
	cfg := store.NewDatabaseConfig()
	s, err := New(cfg)
	if err != nil {
		t.Fatalf("can't apply migrations: %v", err)
	}

	// Create new sql db connection to load fixtures
	db, err := sql.Open(driverName, cfg.URI)
	if err != nil {
		t.Fatalf("can't connect to postgresql test database: %v", err)
	}

	fixtures, err := testfixtures.NewFolder(db, &testfixtures.PostgreSQL{}, "testdata/fixtures")
	if err != nil {
		t.Fatalf("can't read fixtures folder: %v", err)
	}
	if err := fixtures.DetectTestDatabase(); err != nil {
		t.Fatalf("database is not for tests: %v", err)
	}

	if err := fixtures.Load(); err != nil {
		t.Fatalf("can't load fixtures: %v", err)
	}

	return s
}
