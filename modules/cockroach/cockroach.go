package cockroach

import (
	"template-golang/modules/cockroach/handlers"
	"template-golang/modules/cockroach/repositories"
	"template-golang/modules/cockroach/usecases"
)

// Dependencies contains all dependencies for the module
type Cockroach struct {
	Handler    handlers.CockroachHandler
	Repository repositories.CockroachRepository
	Messaging  repositories.CockroachMessaging
	Usecase    usecases.CockroachUsecase
}

// NewDependencies initializes all dependencies
func NewCockroach(
	h handlers.CockroachHandler,
	r repositories.CockroachRepository,
	m repositories.CockroachMessaging,
	u usecases.CockroachUsecase,
) *Cockroach {

	return &Cockroach{
		Handler:    h,
		Repository: r,
		Messaging:  m,
		Usecase:    u,
	}
}
