// This module don't use usecases because library is interact directly with Gin Framework
package auth

import (
	"template-golang/modules/auth/handlers"
	"template-golang/modules/auth/middlewares"
)

type Auth struct {
	Handler    handlers.AuthHandler
	Middleware middlewares.AuthMiddleware
}
