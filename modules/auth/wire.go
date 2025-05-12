//go:build wireinject
// +build wireinject

//go:generate wire
package auth

import (
	"template-golang/config"
	"template-golang/database"
	"template-golang/modules/auth/handlers"
	"template-golang/modules/auth/middlewares"

	"github.com/google/wire"
)

var ProviderSet = wire.NewSet(
	middlewares.ProviderSet,
	handlers.AuthProviderSet,
	wire.Struct(new(Auth), "*"),
)

func Wire(db database.Database, conf *config.Config) (*Auth, error) {
	wire.Build(ProviderSet)
	return &Auth{}, nil
}
