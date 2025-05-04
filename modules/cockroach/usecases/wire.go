//go:build wireinject
// +build wireinject

//go:generate wire
package usecases

import (
	"template-golang/modules/cockroach/repositories"

	"github.com/google/wire"
)

var UsecaseSet = wire.NewSet(
	NewCockroachUsecaseImpl,
	wire.Bind(new(CockroachUsecase), new(*cockroachUsecaseImpl)),
)

func Wire(cockroachRepository repositories.CockroachRepository,
	cockroachMessaging repositories.CockroachMessaging,
) (CockroachUsecase, error) {
	wire.Build(UsecaseSet)
	return &cockroachUsecaseImpl{}, nil
}
