package repositories

import (
	"github.com/google/wire"
)

var ProviderSet = wire.NewSet(
	NewCockroachFCMMessaging,
	wire.Bind(new(CockroachMessaging), new(*cockroachFCMMessaging)),
	NewCockroachPostgresRepository,
	wire.Bind(new(CockroachRepository), new(*cockroachPostgresRepository)),
)
