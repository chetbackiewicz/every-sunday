package database

import (
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/jackc/pgx/v5/stdlib"
)

// Config holds database configuration
type Config struct {
	Host     string
	Port     string
	User     string
	Password string
	Database string
	SSLMode  string
}

// NewPostgresDB creates a new PostgreSQL connection with connection pooling
// Configuration follows research/13-database-integration.md recommendations
func NewPostgresDB(cfg Config) (*sqlx.DB, error) {
	dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s",
		cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.Database, cfg.SSLMode)

	db, err := sqlx.Connect("pgx", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// Budget app is low-traffic, use conservative connection pool settings
	// Reference: research/13-database-integration.md
	db.SetMaxOpenConns(25)                     // Max concurrent connections
	db.SetMaxIdleConns(5)                      // Idle connections ready
	db.SetConnMaxLifetime(30 * time.Minute)    // Recycle connections
	db.SetConnMaxIdleTime(5 * time.Minute)     // Close unused idle connections

	// Verify connection
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	return db, nil
}
