package integration

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"testing"
	"time"

	"template-golang/config"
	db "template-golang/db/sqlc"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/lib/pq"
	"github.com/stretchr/testify/require"
)

// TestDBConfig holds test database configuration
type TestDBConfig struct {
	Host     string
	Port     string
	DBName   string
	Username string
	Password string
	SSLMode  string
}

// DefaultTestDBConfig returns default test database configuration
func DefaultTestDBConfig() *TestDBConfig {
	return &TestDBConfig{
		Host:     getEnv("TEST_DB_HOST", "localhost"),
		Port:     getEnv("TEST_DB_PORT", "5432"),
		DBName:   getEnv("TEST_DB_DBNAME", "template_golang_test"),
		Username: getEnv("TEST_DB_USERNAME", "postgres"),
		Password: getEnv("TEST_DB_PASSWORD", "postgres"),
		SSLMode:  getEnv("TEST_DB_SSLMODE", "disable"),
	}
}

// getEnv gets environment variable with fallback
func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}

// DSN returns the database connection string
func (c *TestDBConfig) DSN() string {
	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s",
		c.Username, c.Password, c.Host, c.Port, c.DBName, c.SSLMode)
}

// PostgresDSN returns the connection string to postgres database (for creating test db)
func (c *TestDBConfig) PostgresDSN() string {
	return fmt.Sprintf("postgres://%s:%s@%s:%s/postgres?sslmode=%s",
		c.Username, c.Password, c.Host, c.Port, c.SSLMode)
}

// SetupTestDB creates and migrates a test database
func SetupTestDB(t *testing.T) (*pgxpool.Pool, func()) {
	t.Helper()

	dbConfig := DefaultTestDBConfig()

	// Create test database
	createTestDatabase(t, dbConfig)

	// Connect to test database
	pool, err := pgxpool.New(context.Background(), dbConfig.DSN())
	require.NoError(t, err, "Failed to connect to test database")

	// Run migrations
	runMigrations(t, dbConfig.DSN())

	// Return cleanup function
	cleanup := func() {
		pool.Close()
		dropTestDatabase(t, dbConfig)
	}

	return pool, cleanup
}

// createTestDatabase creates the test database
func createTestDatabase(t *testing.T, dbConfig *TestDBConfig) {
	t.Helper()

	db, err := sql.Open("postgres", dbConfig.PostgresDSN())
	require.NoError(t, err, "Failed to connect to postgres")
	defer func() {
		if err := db.Close(); err != nil {
			t.Logf("Failed to close db: %v", err)
		}
	}()

	// Drop database if exists
	_, err = db.Exec(fmt.Sprintf("DROP DATABASE IF EXISTS %s", dbConfig.DBName))
	require.NoError(t, err, "Failed to drop existing test database")

	// Create database
	_, err = db.Exec(fmt.Sprintf("CREATE DATABASE %s", dbConfig.DBName))
	require.NoError(t, err, "Failed to create test database")

	log.Printf("Created test database: %s", dbConfig.DBName)
}

// dropTestDatabase drops the test database
func dropTestDatabase(t *testing.T, dbConfig *TestDBConfig) {
	t.Helper()

	db, err := sql.Open("postgres", dbConfig.PostgresDSN())
	if err != nil {
		t.Logf("Failed to connect to postgres for cleanup: %v", err)
		return
	}
	defer func() {
		if err := db.Close(); err != nil {
			t.Logf("Failed to close db: %v", err)
		}
	}()

	_, err = db.Exec(fmt.Sprintf("DROP DATABASE IF EXISTS %s", dbConfig.DBName))
	if err != nil {
		t.Logf("Failed to drop test database: %v", err)
	} else {
		log.Printf("Dropped test database: %s", dbConfig.DBName)
	}
}

// runMigrations runs database migrations
func runMigrations(t *testing.T, dsn string) {
	t.Helper()

	db, err := sql.Open("postgres", dsn)
	require.NoError(t, err, "Failed to connect to test database for migration")
	defer func() {
		if err := db.Close(); err != nil {
			t.Logf("Failed to close db: %v", err)
		}
	}()

	driver, err := postgres.WithInstance(db, &postgres.Config{})
	require.NoError(t, err, "Failed to create postgres driver")

	m, err := migrate.NewWithDatabaseInstance(
		"file://../../db/migrations",
		"postgres", driver)
	require.NoError(t, err, "Failed to create migration instance")
	defer func() {
		if err, _ := m.Close(); err != nil {
			t.Logf("Failed to close migration: %v", err)
		}
	}()

	err = m.Up()
	if err != nil && err != migrate.ErrNoChange {
		require.NoError(t, err, "Failed to run migrations")
	}

	log.Println("Migrations completed successfully")
}

// SetupTestConfig creates a test configuration
func SetupTestConfig(t *testing.T) *config.Config {
	t.Helper()

	// Set session secret for Goth/Gothic
	if err := os.Setenv("SESSION_SECRET", "test_session_secret_123456789"); err != nil {
		t.Fatalf("failed to set SESSION_SECRET: %v", err)
	}

	dbConfig := DefaultTestDBConfig()

	return &config.Config{
		Server: config.ServerConfig{
			Port: 8080,
			Mode: "test",
		},
		Db: config.DbConfig{
			Host:          dbConfig.Host,
			Port:          5432,
			UserName:      dbConfig.Username,
			Password:      dbConfig.Password,
			DBName:        dbConfig.DBName,
			SSLMode:       dbConfig.SSLMode,
			TimeZone:      "Asia/Bangkok",
			MigrationPath: "file://../../db/migrations",
		},
		Auth: config.AuthConfig{
			LineClientID:      "test_client_id",
			LineClientSecret:  "test_client_secret",
			LineCallbackURL:   "http://localhost:8080/api/v1/auth/line/callback",
			LineFECallbackURL: "http://localhost:3000/auth/callback",
			PrivateKeyPath:    "../../ecdsa_private_key.pem",
		},
	}
}

// CreateTestDatabase creates a database instance for testing
func CreateTestDatabase(t *testing.T, pool *pgxpool.Pool) *db.Queries {
	t.Helper()

	return db.New(pool)
}

// WaitForDB waits for database to be ready
func WaitForDB(t *testing.T, pool *pgxpool.Pool, timeout time.Duration) {
	t.Helper()

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	ticker := time.NewTicker(500 * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			t.Fatal("Timeout waiting for database to be ready")
		case <-ticker.C:
			if err := pool.Ping(ctx); err == nil {
				return
			}
		}
	}
}
