package userauth

import (
	"template-golang/modules/userauth/middlewares"
)

// Dependencies contains all dependencies for the module
type UserAuth struct {
	Handler middlewares.UserAuthMiddleware
}
