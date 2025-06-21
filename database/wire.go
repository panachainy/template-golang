//go:build wireinject
// +build wireinject

//go:generate wire
package database

import (
	"template-golang/config"

	"github.com/google/wire"
)

var PostgresProviderSet = wire.NewSet(
	NewPostgres,
	wire.Bind(new(Database), new(*postgresDatabase)),
)

// SQLiteProviderSet provides SQLite database implementation
var SQLiteProviderSet = wire.NewSet(
	ProvideSqliteDatabase,
	wire.Bind(new(Database), new(*SqliteDatabase)),
)

func WirePostgres(conf *config.Config) (Database, error) {
	wire.Build(PostgresProviderSet)
	return &postgresDatabase{}, nil
}

// WireSQLite creates a SQLite database instance using dependency injection
func WireSQLite(dsn string, logMode bool) (Database, error) {
	wire.Build(SQLiteProviderSet)
	return &SqliteDatabase{}, nil
}
