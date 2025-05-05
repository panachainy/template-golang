//go:build wireinject
// +build wireinject

//go:generate wire
package server

import (
	"template-golang/config"
	"template-golang/modules/cockroach"

	"github.com/google/wire"
)

var ProviderSet = wire.NewSet(
	NewGinServer,
	wire.Bind(new(Server), new(*ginServer)),
	// cockroach.ProviderSet,
)

func Wire(conf *config.Config, cockroach *cockroach.Cockroach) (Server, error) {
	wire.Build(ProviderSet)
	return &ginServer{}, nil
}
