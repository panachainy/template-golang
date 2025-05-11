package auth

import (
	"template-golang/modules/auth/handlers"
	"template-golang/modules/auth/middlewares"
	"template-golang/modules/auth/usecases"
)

type Auth struct {
	Handler    handlers.AuthHandler
	Middleware middlewares.AuthMiddleware
	Usecase    usecases.AuthUsecase
}
