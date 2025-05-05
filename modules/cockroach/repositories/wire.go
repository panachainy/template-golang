package repositories

import (
	"github.com/google/wire"
)

var RepositorySet = wire.NewSet(
	NewCockroachFCMMessaging,
	wire.Bind(new(CockroachMessaging), new(*cockroachFCMMessaging)),
	NewCockroachPostgresRepository,
	wire.Bind(new(CockroachRepository), new(*cockroachPostgresRepository)),
)
