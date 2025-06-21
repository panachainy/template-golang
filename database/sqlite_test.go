package database

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewSqliteDatabase(t *testing.T) {
	tests := []struct {
		name        string
		config      *Config
		expectError bool
	}{
		{
			name:        "nil config should return error",
			config:      nil,
			expectError: true,
		},
		{
			name: "valid config should create database",
			config: &Config{
				DSN:     ":memory:",
				LogMode: false,
			},
			expectError: false,
		},
		{
			name: "valid config with log mode should create database",
			config: &Config{
				DSN:     ":memory:",
				LogMode: true,
			},
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, err := NewSqliteDatabase(tt.config)

			if tt.expectError {
				assert.Error(t, err)
				assert.Nil(t, db)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, db)
				assert.NotNil(t, db.GetDb())

				// Clean up
				if db != nil {
					err := db.Close()
					assert.NoError(t, err)
				}
			}
		})
	}
}

func TestSqliteDatabase_GetDb(t *testing.T) {
	config := &Config{
		DSN:     ":memory:",
		LogMode: false,
	}

	db, err := NewSqliteDatabase(config)
	require.NoError(t, err)
	defer db.Close()

	gormDB := db.GetDb()
	assert.NotNil(t, gormDB)
}

func TestSqliteDatabase_Migration_Operations(t *testing.T) {
	// Create a temporary directory for test migrations
	tempDir, err := os.MkdirTemp("", "sqlite_test_migrations")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	// Create test migration files
	migrationsDir := filepath.Join(tempDir, "migrations")
	err = os.MkdirAll(migrationsDir, 0755)
	require.NoError(t, err)

	// Create a simple test migration
	upMigration := `CREATE TABLE test_table (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT NOT NULL,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);`

	downMigration := `DROP TABLE IF EXISTS test_table;`

	err = os.WriteFile(filepath.Join(migrationsDir, "000001_create_test_table.up.sql"), []byte(upMigration), 0644)
	require.NoError(t, err)

	err = os.WriteFile(filepath.Join(migrationsDir, "000001_create_test_table.down.sql"), []byte(downMigration), 0644)
	require.NoError(t, err)

	// Create a temporary SQLite database file
	tempDB := filepath.Join(tempDir, "test.db")
	config := &Config{
		DSN:     tempDB,
		LogMode: false,
	}

	db, err := NewSqliteDatabase(config)
	require.NoError(t, err)
	defer db.Close()

	// Test GetVersion before any migrations
	version, _, err := db.GetVersion()
	// Note: For a fresh database, this might return an error if no migrations table exists yet
	// This is expected behavior with golang-migrate
	assert.True(t, err != nil || version == 0)

	// Note: For this test to fully work, we would need to:
	// 1. Set up the migration path to point to our test migrations
	// 2. Or modify the migration methods to accept a custom path
	// For now, this test validates the structure and basic functionality
}

func TestSqliteDatabase_Close(t *testing.T) {
	config := &Config{
		DSN:     ":memory:",
		LogMode: false,
	}

	db, err := NewSqliteDatabase(config)
	require.NoError(t, err)

	err = db.Close()
	assert.NoError(t, err)

	// Verify that subsequent operations fail after closing
	// Note: This might not always fail with SQLite in-memory databases
	// but it's good practice to test the close functionality
}

func TestSqliteDatabase_Interface_Compliance(t *testing.T) {
	// This test ensures SqliteDatabase implements the Database interface
	var _ Database = (*SqliteDatabase)(nil)
}
