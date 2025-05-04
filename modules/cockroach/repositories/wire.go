//go:build wireinject
// +build wireinject

//go:generate wire
package repositories

import (
	"template-golang/database"

	"github.com/google/wire"
)

var RepositorySet = wire.NewSet(
	NewCockroachFCMMessaging,
	wire.Bind(new(CockroachMessaging), new(*cockroachFCMMessaging)),
	NewCockroachPostgresRepository,
	wire.Bind(new(CockroachRepository), new(*cockroachPostgresRepository)),
)

func Wire(db database.Database) (CockroachMessaging, CockroachRepository, error) {
	wire.Build(RepositorySet)
	return &cockroachFCMMessaging{}, &cockroachPostgresRepository{}, nil
}
