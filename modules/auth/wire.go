//go:build wireinject
// +build wireinject

//go:generate wire
package auth

import (
	"template-golang/config"
	"template-golang/database"
	"template-golang/modules/auth/handlers"
	"template-golang/modules/auth/middlewares"
	"template-golang/modules/auth/repositories"
	"template-golang/modules/auth/usecases"

	"github.com/google/wire"
)

var ProviderSet = wire.NewSet(
	middlewares.ProviderSet,
	handlers.AuthProviderSet,
	usecases.ProviderSet,
	repositories.ProviderSet,
	wire.Struct(new(Auth), "*"),
)

func Wire(db database.Database, conf *config.Config) (*Auth, error) {
	wire.Build(ProviderSet)
	return &Auth{}, nil
}
