package database

import (
	"fmt"
	"template-golang/config"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/labstack/gommon/log"
)

// MigrationManager handles database migrations using golang-migrate
type MigrationManager struct {
	db     Database
	config *config.Config
}

// ProvideMigrationManager creates a new migration manager
func ProvideMigrationManager(db Database, conf *config.Config) *MigrationManager {
	return &MigrationManager{
		db:     db,
		config: conf,
	}
}

// RunMigrations runs all pending migrations
func (m *MigrationManager) RunMigrations() error {
	log.Info("Running database migrations...")

	// Get the underlying SQL DB from GORM
	sqlDB, err := m.db.GetDb().DB()
	if err != nil {
		return fmt.Errorf("failed to get underlying sql.DB: %w", err)
	}

	// Create postgres driver instance
	driver, err := postgres.WithInstance(sqlDB, &postgres.Config{})
	if err != nil {
		return fmt.Errorf("failed to create postgres driver: %w", err)
	}

	// Create migrate instance
	migration, err := migrate.NewWithDatabaseInstance(
		"file://migrations",
		"postgres", driver)
	if err != nil {
		return fmt.Errorf("failed to create migrate instance: %w", err)
	}
	defer migration.Close()

	// Run migrations
	if err := migration.Up(); err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("failed to run migrations: %w", err)
	}

	log.Info("Database migrations completed successfully")
	return nil
}

// RollbackMigrations rolls back migrations by specified number of steps
func (m *MigrationManager) RollbackMigrations(steps int) error {
	log.Infof("Rolling back %d migration(s)...", steps)

	// Get the underlying SQL DB from GORM
	sqlDB, err := m.db.GetDb().DB()
	if err != nil {
		return fmt.Errorf("failed to get underlying sql.DB: %w", err)
	}

	// Create postgres driver instance
	driver, err := postgres.WithInstance(sqlDB, &postgres.Config{})
	if err != nil {
		return fmt.Errorf("failed to create postgres driver: %w", err)
	}

	// Create migrate instance
	migration, err := migrate.NewWithDatabaseInstance(
		"file://migrations",
		"postgres", driver)
	if err != nil {
		return fmt.Errorf("failed to create migrate instance: %w", err)
	}
	defer migration.Close()

	// Roll back migrations
	if err := migration.Steps(-steps); err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("failed to rollback migrations: %w", err)
	}

	log.Infof("Rollback of %d migration(s) completed successfully", steps)
	return nil
}

// GetVersion returns the current migration version
func (m *MigrationManager) GetVersion() (uint, bool, error) {
	// Get the underlying SQL DB from GORM
	sqlDB, err := m.db.GetDb().DB()
	if err != nil {
		return 0, false, fmt.Errorf("failed to get underlying sql.DB: %w", err)
	}

	// Create postgres driver instance
	driver, err := postgres.WithInstance(sqlDB, &postgres.Config{})
	if err != nil {
		return 0, false, fmt.Errorf("failed to create postgres driver: %w", err)
	}

	// Create migrate instance
	migration, err := migrate.NewWithDatabaseInstance(
		"file://migrations",
		"postgres", driver)
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

// CreateMigration creates a new migration file pair (up and down)
func (m *MigrationManager) CreateMigration(name string) error {
	// This is typically done using the migrate CLI tool
	// But we can provide a helper function here
	log.Infof("To create a new migration, use the migrate CLI:")
	log.Infof("migrate create -ext sql -dir migrations -seq %s", name)
	return nil
}
