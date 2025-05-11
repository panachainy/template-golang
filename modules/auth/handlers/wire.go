package handlers

import (
	"github.com/google/wire"
)

var AuthProviderSet = wire.NewSet(
	Provide,
	wire.Bind(new(AuthHandler), new(*authHttpHandler)),
)
