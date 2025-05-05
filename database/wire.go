//go:build wireinject
// +build wireinject

//go:generate wire
package database

import (
	"template-golang/config"

	"github.com/google/wire"
)

var ProviderSet = wire.NewSet(
	NewPostgresDatabase,
	wire.Bind(new(Database), new(*postgresDatabase)),
)

func Wire(conf *config.Config) (Database, error) {
	wire.Build(ProviderSet)
	return &postgresDatabase{}, nil
}
