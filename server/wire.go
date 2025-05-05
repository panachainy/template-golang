//go:build wireinject
// +build wireinject

//go:generate wire
package server

import (
	"template-golang/config"
	"template-golang/database"
	"template-golang/modules/cockroach"

	"github.com/google/wire"
)

var ProviderSet = wire.NewSet(
	NewGinServer,
	wire.Bind(new(Server), new(*ginServer)),

	// cores
	config.ProviderSet,
	database.ProviderSet,

	// modules
	cockroach.ProviderSet,
)

func Wire() (Server, error) {
	wire.Build(ProviderSet)
	return &ginServer{}, nil
}
