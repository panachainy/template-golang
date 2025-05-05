package usecases

import (
	"github.com/google/wire"
)

var UsecaseSet = wire.NewSet(
	NewCockroachUsecaseImpl,
	wire.Bind(new(CockroachUsecase), new(*cockroachUsecaseImpl)),
)
