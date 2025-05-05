package usecases

import (
	"github.com/google/wire"
)

var ProviderSet = wire.NewSet(
	NewCockroachUsecaseImpl,
	wire.Bind(new(CockroachUsecase), new(*cockroachUsecaseImpl)),
)
