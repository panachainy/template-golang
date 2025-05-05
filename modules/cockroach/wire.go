//go:build wireinject
// +build wireinject

//go:generate wire
package cockroach

import (
	"template-golang/database"
	"template-golang/modules/cockroach/handlers"
	"template-golang/modules/cockroach/repositories"
	"template-golang/modules/cockroach/usecases"

	"github.com/google/wire"
)

var CockroachSet = wire.NewSet(
	NewCockroach,
	handlers.HandlerSet,
	repositories.RepositorySet,
	usecases.UsecaseSet,
)

func Wire(db database.Database) (*Cockroach, error) {
	wire.Build(CockroachSet)
	return &Cockroach{}, nil
}
