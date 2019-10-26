package store

import (
	"github.com/spf13/viper"
)

const (
	dbVersion = "db_version"
	dbURI     = "db_uri"
)

var cfg *viper.Viper

func init() {
	cfg = viper.New()
	cfg.AutomaticEnv()

	_ = cfg.BindEnv(dbURI, "DB_URI")
	cfg.SetDefault(dbURI, "postgresql://postgres:postgres@localhost:5432/test?sslmode=disable")

	_ = cfg.BindEnv(dbVersion, "DB_VERSION")
	cfg.SetDefault(dbVersion, 0)
}

// DatabaseConfig connection database configuration
type DatabaseConfig struct {
	URI     string
	Version uint
}

// NewDatabaseConfig constructor for DatabaseConfig
func NewDatabaseConfig() *DatabaseConfig {
	version := cfg.GetInt(dbVersion)
	if version < 0 {
		version = 0
	}
	return &DatabaseConfig{
		URI:     cfg.GetString(dbURI),
		Version: uint(version),
	}
}
