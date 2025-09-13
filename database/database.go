package database

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Database interface {
	GetPool() *pgxpool.Pool
	Close()
	Ping(ctx context.Context) error
}
