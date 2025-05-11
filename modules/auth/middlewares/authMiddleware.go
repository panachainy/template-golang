package middlewares

import (
	"github.com/gin-gonic/gin"
)

type AuthMiddleware interface {
	Handle() gin.HandlerFunc
}
