package usecases

import (
	"github.com/google/wire"
)

var ProviderSet = wire.NewSet(
	Provide,
	wire.Bind(new(CockroachUsecase), new(*cockroachUsecaseImpl)),
)
