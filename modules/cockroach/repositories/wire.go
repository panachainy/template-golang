package repositories

import (
	"github.com/google/wire"
)

var ProviderSet = wire.NewSet(
	ProvideFCMMessaging,
	wire.Bind(new(CockroachMessaging), new(*cockroachFCMMessaging)),
	ProvidePostgresRepository,
	wire.Bind(new(CockroachRepository), new(*cockroachPostgresRepository)),
)
