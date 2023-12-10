package models

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

type PostgresConfig struct {
	Type string

	User     string
	Password string
	Host     string
	Port     string
	Database string
	sslMode  string
}

func (cfg PostgresConfig) String() string {
	// postgresql://<username>:<password>@<database_ip>/todos?sslmode=disable
	return fmt.Sprintf("%s://%s:%s@%s/%s?sslmode=%s", cfg.Type, cfg.User, cfg.Password, cfg.Host, cfg.Database, cfg.sslMode)
}

func DefaultPostgresConfig() PostgresConfig {
	return PostgresConfig{
		Type:     "postgres",
		User:     "baloo",
		Password: "junglebook",
		Host:     "localhost",
		Port:     "3306",
		Database: "lenslocked",
		sslMode:  "disable",
	}
}

// Open will open a SQL connnection with the provided
// Postgres database. cakkers of Open need to ensure
// db.Close() method
func Open(config PostgresConfig) (*sql.DB, error) {
	db, err := sql.Open(config.Type, config.String())

	if err != nil {
		return nil, fmt.Errorf("Open: %w", err)
	}
	return db, nil
}
