package middlewares

import (
	"github.com/gin-gonic/gin"
)

type userAuthMiddleware struct {
}

func Provide() UserAuthMiddleware {
	return &userAuthMiddleware{}
}

func (m *userAuthMiddleware) Handle() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Add your authentication logic here
		c.Next()
	}
}
