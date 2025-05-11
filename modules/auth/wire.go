//go:build wireinject
// +build wireinject

//go:generate wire
package auth

import (
	"template-golang/database"
	"template-golang/modules/auth/middlewares"

	"github.com/google/wire"
)

var ProviderSet = wire.NewSet(
	middlewares.ProviderSet,
	wire.Struct(new(Auth), "*"),
)

func Wire(db database.Database) (*Auth, error) {
	wire.Build(ProviderSet)
	return &Auth{}, nil
}
