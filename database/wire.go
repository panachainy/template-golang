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

func WirePostgres(conf *config.Config) (Database, error) {
	wire.Build(PostgresProviderSet)
	return &postgresDatabase{}, nil
}
