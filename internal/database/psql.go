package database

import (
	"database/sql"
	"fmt"
	"wabot/internal/config"

	_ "github.com/lib/pq"
)

// DB holds the database
type DB struct{ *sql.DB }

// database instance
var defaultDB = &DB{}

// connect sets the db client of database using configuration
func (db *DB) connect(cfg *config.DB) (err error) {
	dbURI := fmt.Sprintf("user=%s password=%s host=%s sslmode=disable port=%d dbname=%s",
		cfg.User,
		cfg.Password,
		cfg.Host,
		cfg.Port,
		cfg.Name,
	)

	db.DB, err = sql.Open("postgres", dbURI)
	if err != nil {
		return err
	}

	// Try to ping database.
	if err := db.Ping(); err != nil {
		defer db.Close()
		return fmt.Errorf("can't sent ping to database, %w", err)
	}

	return nil
}

// GetDB returns db instance
func GetDB() *DB {
	return defaultDB
}

// ConnectDB sets the db client of database using default configuration
func ConnectDB() error {
	return defaultDB.connect(config.DBCfg())
}
