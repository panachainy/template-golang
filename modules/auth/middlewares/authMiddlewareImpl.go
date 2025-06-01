package middlewares

import (
	"net/http"
	"strings"
	"template-golang/modules/auth/usecases"

	"github.com/gin-gonic/gin"
	"github.com/labstack/gommon/log"
)

type userAuthMiddleware struct {
	jwtUsecase usecases.JWTUsecase
}

func Provide(jwtUsecase usecases.JWTUsecase) *userAuthMiddleware {
	return &userAuthMiddleware{
		jwtUsecase: jwtUsecase,
	}
}

func (m *userAuthMiddleware) Handle() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Extract token from Authorization header
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			log.Warn("Missing Authorization header")
			c.JSON(http.StatusUnauthorized, gin.H{
				"error":   "Unauthorized",
				"message": "Missing authorization header",
			})
			c.Abort()
			return
		}

		// Check for Bearer token format
		tokenParts := strings.Split(authHeader, " ")
		if len(tokenParts) != 2 || tokenParts[0] != "Bearer" || strings.TrimSpace(tokenParts[1]) == "" {
			log.Warn("Invalid Authorization header format")
			c.JSON(http.StatusUnauthorized, gin.H{
				"error":   "Unauthorized",
				"message": "Invalid authorization header format",
			})
			c.Abort()
			return
		}

		tokenString := tokenParts[1]

		// Verify the token
		result, err := m.jwtUsecase.ValidateJWT(tokenString)
		if err != nil {
			log.Errorf("Token verification error: %v", err)
			c.JSON(http.StatusUnauthorized, gin.H{
				"error":   "Unauthorized",
				"message": "Token verification failed",
			})
			c.Abort()
			return
		}

		// Handle different validation states
		if result.NotExist {
			log.Warn("Token not provided")
			c.JSON(http.StatusUnauthorized, gin.H{
				"error":   "Unauthorized",
				"message": "Token not provided",
			})
			c.Abort()
			return
		}

		if result.Expired {
			log.Warn("Token has expired")
			c.JSON(http.StatusUnauthorized, gin.H{
				"error":   "Unauthorized",
				"message": "Token has expired",
			})
			c.Abort()
			return
		}

		if !result.Valid {
			log.Warn("Invalid token")
			c.JSON(http.StatusUnauthorized, gin.H{
				"error":   "Unauthorized",
				"message": "Invalid token",
			})
			c.Abort()
			return
		}

		// TODO: check this

		// Token is valid, set user context
		c.Set("userID", result.UserID)
		c.Set("claims", result.Claims)

		log.Infof("Successfully authenticated user: %s", result.UserID)

		c.Next()
	}
}
