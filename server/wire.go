//go:build wireinject
// +build wireinject

//go:generate wire
package server

import (
	"template-golang/config"
	cockroachHandlers "template-golang/modules/cockroach/handlers"

	"github.com/google/wire"
)

var ProviderSet = wire.NewSet(
	NewGinServer,
	wire.Bind(new(Server), new(*ginServer)),
)

func Wire(conf *config.Config, cockroachH cockroachHandlers.CockroachHandler) (Server, error) {
	wire.Build(ProviderSet)
	return &ginServer{}, nil
}
