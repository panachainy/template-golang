package middlewares

import (
	"github.com/google/wire"
)

var ProviderSet = wire.NewSet(
	Provide,
	wire.Bind(new(UserAuthMiddleware), new(*userAuthMiddleware)),
)
