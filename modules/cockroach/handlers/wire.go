package handlers

import (
	"github.com/google/wire"
)

var ProviderSet = wire.NewSet(
	NewCockroachHttpHandler,
	wire.Bind(new(CockroachHandler), new(*cockroachHttpHandler)),
)
