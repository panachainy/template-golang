//go:build wireinject
// +build wireinject

//go:generate wire
package userauth

import (
	"template-golang/database"
	"template-golang/modules/userauth/middlewares"

	"github.com/google/wire"
)

var ProviderSet = wire.NewSet(
	middlewares.ProviderSet,
	wire.Struct(new(UserAuth), "*"),
)

func Wire(db database.Database) (*UserAuth, error) {
	wire.Build(ProviderSet)
	return &UserAuth{}, nil
}
