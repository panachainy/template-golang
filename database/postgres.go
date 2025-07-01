package database

import (
	"fmt"
	"sync"
	"template-golang/config"
	"template-golang/pkg/logger"

	"github.com/golang-migrate/migrate/v4"
	pgMigrate "github.com/golang-migrate/migrate/v4/database/postgres"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type postgresDatabase struct {
	Db   *gorm.DB
	conf *config.Config
}

var (
	once       sync.Once
	dbInstance *postgresDatabase
)

func NewPostgres(conf *config.Config) *postgresDatabase {
	once.Do(func() {
		dsn := fmt.Sprintf(
			"host=%s user=%s password=%s dbname=%s port=%d sslmode=%s TimeZone=%s",
			conf.Db.Host,
			conf.Db.UserName,
			conf.Db.Password,
			conf.Db.DBName,
			conf.Db.Port,
			conf.Db.SSLMode,
			conf.Db.TimeZone,
		)

		db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
		if err != nil {
			panic("failed to connect database")
		}

		dbInstance = &postgresDatabase{Db: db, conf: conf}
	})

	return dbInstance
}

func (p *postgresDatabase) GetDb() *gorm.DB {
	return dbInstance.Db
}

func (p *postgresDatabase) MigrateUp() error {
	defer p.Close()
	logger.Info("Running database migrations...")

	// Get the underlying SQL DB from GORM
	sqlDB, err := p.Db.DB()
	if err != nil {
		return fmt.Errorf("failed to get underlying sql.DB: %w", err)
	}

	// Create postgres driver instance
	driver, err := pgMigrate.WithInstance(sqlDB, &pgMigrate.Config{})
	if err != nil {
		return fmt.Errorf("failed to create postgres driver: %w", err)
	}

	// Create migrate instance
	migration, err := migrate.NewWithDatabaseInstance(
		p.conf.Db.MigrationPath,
		"postgres", driver)
	if err != nil {
		return fmt.Errorf("failed to create migrate instance: %w", err)
	}
	defer migration.Close()

	// Run migrations
	if err := migration.Up(); err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("failed to run migrations: %w", err)
	}

	logger.Info("Database migrations completed successfully")
	return nil
}

func (p *postgresDatabase) MigrateDown(steps int) error {
	defer p.Close()
	logger.Infof("Rolling back %d migration(s)...", steps)

	// Get the underlying SQL DB from GORM
	sqlDB, err := p.Db.DB()
	if err != nil {
		return fmt.Errorf("failed to get underlying sql.DB: %w", err)
	}

	// Create postgres driver instance
	driver, err := pgMigrate.WithInstance(sqlDB, &pgMigrate.Config{})
	if err != nil {
		return fmt.Errorf("failed to create postgres driver: %w", err)
	}

	// Create migrate instance
	migration, err := migrate.NewWithDatabaseInstance(
		p.conf.Db.MigrationPath,
		"postgres", driver)
	if err != nil {
		return fmt.Errorf("failed to create migrate instance: %w", err)
	}
	defer migration.Close()

	// Roll back migrations
	if err := migration.Steps(-steps); err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("failed to rollback migrations: %w", err)
	}

	logger.Infof("Rollback of %d migration(s) completed successfully", steps)
	return nil
}

func (p *postgresDatabase) GetVersion() (uint, bool, error) {
	// Get the underlying SQL DB from GORM
	sqlDB, err := p.Db.DB()
	if err != nil {
		return 0, false, fmt.Errorf("failed to get underlying sql.DB: %w", err)
	}

	// Create postgres driver instance
	driver, err := pgMigrate.WithInstance(sqlDB, &pgMigrate.Config{})
	if err != nil {
		return 0, false, fmt.Errorf("failed to create postgres driver: %w", err)
	}

	// Create migrate instance
	migration, err := migrate.NewWithDatabaseInstance(
		p.conf.Db.MigrationPath,
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

func (p *postgresDatabase) Close() error {
	// Get the underlying SQL DB from GORM
	sqlDB, err := p.Db.DB()
	if err != nil {
		return fmt.Errorf("failed to get underlying sql.DB: %w", err)
	}

	if err := sqlDB.Close(); err != nil {
		return fmt.Errorf("failed to close database connection: %w", err)
	}

	// Reset the singleton instance and sync.Once
	once = sync.Once{}
	dbInstance = nil

	logger.Info("Database connection closed successfully")
	return nil
}
