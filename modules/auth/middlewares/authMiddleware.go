//go:generate mockgen -source=authMiddleware.go -destination=./mocks/mock_auth_middleware.go -package=mock
package middlewares

import (
	"template-golang/modules/auth/models"

	"github.com/gin-gonic/gin"
)

type AuthMiddleware interface {
	Handle() gin.HandlerFunc
	Allows(roles []models.Role) gin.HandlerFunc
}
