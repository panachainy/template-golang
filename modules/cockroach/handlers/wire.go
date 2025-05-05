package handlers

import (
	"github.com/google/wire"
)

var HandlerSet = wire.NewSet(
	NewCockroachHttpHandler,
	wire.Bind(new(CockroachHandler), new(*cockroachHttpHandler)),
)
