//go:build wireinject
// +build wireinject

//go:generate wire
package server

import (
	"template-golang/config"
	"template-golang/database"
	"template-golang/modules/cockroach"
	"template-golang/modules/userauth"

	"github.com/google/wire"
)

var ProviderSet = wire.NewSet(
	Provide,
	wire.Bind(new(Server), new(*ginServer)),

	// cores
	config.ProviderSet,
	database.ProviderSet,

	// modules
	cockroach.ProviderSet,
	userauth.ProviderSet,
)

func Wire() (Server, error) {
	wire.Build(ProviderSet)
	return &ginServer{}, nil
}
