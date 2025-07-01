package database

import (
	"fmt"
	"template-golang/pkg/logger"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/sqlite3"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	gormLogger "gorm.io/gorm/logger"
)

// Config holds the SQLite database configuration
type Config struct {
	DSN           string
	LogMode       bool
	MigrationPath string // Path to migration files, defaults to "file://db/migrations"
}

// SqliteDatabase implements the Database interface for SQLite
type SqliteDatabase struct {
	db     *gorm.DB
	config *Config
}

// NewSqliteDatabase creates a new SQLite database instance
func NewSqliteDatabase(config *Config) (*SqliteDatabase, error) {
	if config == nil {
		return nil, fmt.Errorf("config cannot be nil")
	}

	gormConfig := &gorm.Config{}
	if config.LogMode {
		gormConfig.Logger = gormLogger.Default.LogMode(gormLogger.Info)
	} else {
		gormConfig.Logger = gormLogger.Default.LogMode(gormLogger.Silent)
	}

	db, err := gorm.Open(sqlite.Open(config.DSN), gormConfig)
	if err != nil {
		logger.Errorf("failed to connect to SQLite database: %v", err)
		return nil, fmt.Errorf("failed to connect to SQLite database: %w", err)
	}

	logger.Info("successfully connected to SQLite database")

	return &SqliteDatabase{
		db:     db,
		config: config,
	}, nil
}

// ProvideSqliteDatabase creates a new SQLite database instance with default configuration
// This follows the project's dependency injection pattern
func ProvideSqliteDatabase(dsn string, logMode bool) (*SqliteDatabase, error) {
	config := &Config{
		DSN:           dsn,
		LogMode:       logMode,
		MigrationPath: "", // Will use default "file://db/migrations"
	}
	return NewSqliteDatabase(config)
}

// ProvideSqliteDatabaseWithMigrationPath creates a new SQLite database instance with custom migration path
func ProvideSqliteDatabaseWithMigrationPath(dsn string, logMode bool, migrationPath string) (*SqliteDatabase, error) {
	config := &Config{
		DSN:           dsn,
		LogMode:       logMode,
		MigrationPath: migrationPath,
	}
	return NewSqliteDatabase(config)
}

// GetDb returns the GORM database instance
func (s *SqliteDatabase) GetDb() *gorm.DB {
	return s.db
}

// getMigrationPath returns the migration path, defaulting to "file://db/migrations" if not set
func (s *SqliteDatabase) getMigrationPath() string {
	if s.config.MigrationPath == "" {
		return "file://db/migrations"
	}
	return s.config.MigrationPath
}

// MigrateUp applies all pending migrations
func (s *SqliteDatabase) MigrateUp() error {
	logger.Info("running database migrations up")

	// Get the underlying SQL DB from GORM
	sqlDB, err := s.db.DB()
	if err != nil {
		return fmt.Errorf("failed to get underlying sql.DB: %w", err)
	}

	// Create SQLite driver instance
	driver, err := sqlite3.WithInstance(sqlDB, &sqlite3.Config{})
	if err != nil {
		return fmt.Errorf("failed to create sqlite3 driver: %w", err)
	}

	// Create migrate instance
	migration, err := migrate.NewWithDatabaseInstance(
		s.getMigrationPath(),
		"sqlite3", driver)
	if err != nil {
		return fmt.Errorf("failed to create migrate instance: %w", err)
	}
	defer migration.Close()

	// Run migrations
	if err := migration.Up(); err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("failed to run migrations: %w", err)
	}

	logger.Info("SQLite migration up completed")
	return nil
}

// MigrateDown rolls back the specified number of migration steps
func (s *SqliteDatabase) MigrateDown(steps int) error {
	logger.Infof("rolling back %d migration steps", steps)

	// Get the underlying SQL DB from GORM
	sqlDB, err := s.db.DB()
	if err != nil {
		return fmt.Errorf("failed to get underlying sql.DB: %w", err)
	}

	// Create SQLite driver instance
	driver, err := sqlite3.WithInstance(sqlDB, &sqlite3.Config{})
	if err != nil {
		return fmt.Errorf("failed to create sqlite3 driver: %w", err)
	}

	// Create migrate instance
	migration, err := migrate.NewWithDatabaseInstance(
		s.getMigrationPath(),
		"sqlite3", driver)
	if err != nil {
		return fmt.Errorf("failed to create migrate instance: %w", err)
	}
	defer migration.Close()

	// Roll back migrations
	if err := migration.Steps(-steps); err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("failed to rollback migrations: %w", err)
	}

	logger.Infof("SQLite migration down completed for %d steps", steps)
	return nil
}

// GetVersion returns the current migration version
func (s *SqliteDatabase) GetVersion() (uint, bool, error) {
	logger.Info("getting current migration version")

	// Get the underlying SQL DB from GORM
	sqlDB, err := s.db.DB()
	if err != nil {
		return 0, false, fmt.Errorf("failed to get underlying sql.DB: %w", err)
	}

	// Create SQLite driver instance
	driver, err := sqlite3.WithInstance(sqlDB, &sqlite3.Config{})
	if err != nil {
		return 0, false, fmt.Errorf("failed to create sqlite3 driver: %w", err)
	}

	// Create migrate instance
	migration, err := migrate.NewWithDatabaseInstance(
		s.getMigrationPath(),
		"sqlite3", driver)
	if err != nil {
		return 0, false, fmt.Errorf("failed to create migrate instance: %w", err)
	}
	defer migration.Close()

	// Get current version
	version, dirty, err := migration.Version()
	if err != nil && err != migrate.ErrNilVersion {
		return 0, false, fmt.Errorf("failed to get migration version: %w", err)
	}

	return version, dirty, nil
}

// Close closes the database connection
func (s *SqliteDatabase) Close() error {
	sqlDB, err := s.db.DB()
	if err != nil {
		return fmt.Errorf("failed to get underlying sql.DB: %w", err)
	}

	if err := sqlDB.Close(); err != nil {
		logger.Errorf("failed to close SQLite database connection: %v", err)
		return fmt.Errorf("failed to close database connection: %w", err)
	}

	logger.Info("SQLite database connection closed")
	return nil
}
