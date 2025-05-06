package middlewares

import (
	"github.com/gin-gonic/gin"
)

type UserAuthMiddleware interface {
	Handle() gin.HandlerFunc
}
