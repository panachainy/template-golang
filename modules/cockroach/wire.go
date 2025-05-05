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

var ProviderSet = wire.NewSet(
	handlers.ProviderSet,
	repositories.ProviderSet,
	usecases.ProviderSet,
	wire.Struct(new(Cockroach), "*"),
)

func Wire(db database.Database) (*Cockroach, error) {
	wire.Build(ProviderSet)
	return &Cockroach{}, nil
}
