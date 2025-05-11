package auth

import (
	"template-golang/modules/auth/middlewares"
)

// Dependencies contains all dependencies for the module
type Auth struct {
	Handler middlewares.AuthMiddleware
}
