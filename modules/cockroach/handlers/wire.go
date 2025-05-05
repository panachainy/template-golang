package handlers

import (
	"github.com/google/wire"
)

var ProviderSet = wire.NewSet(
	Provide,
	wire.Bind(new(CockroachHandler), new(*cockroachHttpHandler)),
)
