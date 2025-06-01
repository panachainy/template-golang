//go:generate mockgen -source=authMiddleware.go -destination=../../../mock/mock_auth_middleware.go -package=mock
package middlewares

import (
	"github.com/gin-gonic/gin"
)

type AuthMiddleware interface {
	Handle() gin.HandlerFunc
}
