package repositories

import (
	"github.com/google/wire"
)

var ProviderSet = wire.NewSet(
	ProvideAuthRepository,
	wire.Bind(new(AuthRepository), new(*authPostgresRepository)),
)
