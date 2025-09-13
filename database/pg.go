package database

import (
	"context"
	"fmt"
	"template-golang/config"

	"github.com/jackc/pgx/v5/pgxpool"
)

type postgresDatabase struct {
	pool *pgxpool.Pool
}

func NewPostgresDatabase(cfg *config.Config) (Database, error) {
	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s TimeZone=%s",
		cfg.Db.Host,
		cfg.Db.Port,
		cfg.Db.UserName,
		cfg.Db.Password,
		cfg.Db.DBName,
		cfg.Db.SSLMode,
		cfg.Db.TimeZone,
	)

	pool, err := pgxpool.New(context.Background(), dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to create connection pool: %w", err)
	}

	// Test the connection
	if err := pool.Ping(context.Background()); err != nil {
		pool.Close()
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	return &postgresDatabase{
		pool: pool,
	}, nil
}

func (d *postgresDatabase) GetPool() *pgxpool.Pool {
	return d.pool
}

func (d *postgresDatabase) Close() {
	if d.pool != nil {
		d.pool.Close()
	}
}

func (d *postgresDatabase) Ping(ctx context.Context) error {
	return d.pool.Ping(ctx)
}
