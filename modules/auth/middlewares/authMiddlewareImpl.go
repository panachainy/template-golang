package middlewares

import (
	"net/http"
	"strings"
	"template-golang/modules/auth/models"
	"template-golang/modules/auth/usecases"
	"template-golang/pkg/logger"

	"github.com/gin-gonic/gin"
)

type userAuthMiddleware struct {
	jwtUsecase usecases.JWTUsecase
}

func NewAuthMiddleware(jwtUsecase usecases.JWTUsecase) AuthMiddleware {
	return &userAuthMiddleware{
		jwtUsecase: jwtUsecase,
	}
}

func (m *userAuthMiddleware) Handle() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Extract token from Authorization header
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			logger.Warn("Missing Authorization header")
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
			logger.Warn("Invalid Authorization header format")
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
			logger.Errorf("Token verification error: %v", err)
			c.JSON(http.StatusUnauthorized, gin.H{
				"error":   "Unauthorized",
				"message": "Token verification failed",
			})
			c.Abort()
			return
		}

		// Handle different validation states
		if result.NotExist {
			logger.Warn("Token not provided")
			c.JSON(http.StatusUnauthorized, gin.H{
				"error":   "Unauthorized",
				"message": "Token not provided",
			})
			c.Abort()
			return
		}

		if result.Expired {
			logger.Warn("Token has expired")
			c.JSON(http.StatusUnauthorized, gin.H{
				"error":   "Unauthorized",
				"message": "Token has expired",
			})
			c.Abort()
			return
		}

		if !result.Valid {
			logger.Warn("Invalid token")
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

		logger.Infof("Successfully authenticated user: %s", result.UserID)

		c.Next()
	}
}

func (m *userAuthMiddleware) Allows(roles []models.Role) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get user claims from context (set by Handle middleware)
		claims, exists := c.Get("claims")
		if !exists {
			logger.Warn("No user claims found in context")
			c.JSON(http.StatusUnauthorized, gin.H{
				"error":   "Unauthorized",
				"message": "No user claims found",
			})
			c.Abort()
			return
		}

		userClaims, ok := claims.(map[string]interface{})
		if !ok {
			logger.Warn("Invalid claims format")
			c.JSON(http.StatusForbidden, gin.H{
				"error":   "Forbidden",
				"message": "Invalid user claims",
			})
			c.Abort()
			return
		}

		// Extract user role from claims
		userRole, exists := userClaims["role"]
		if !exists {
			logger.Warn("No role found in user claims")
			c.JSON(http.StatusForbidden, gin.H{
				"error":   "Forbidden",
				"message": "No role found in user claims",
			})
			c.Abort()
			return
		}

		userRoleStr, ok := userRole.(string)
		if !ok {
			logger.Warn("Invalid role format in claims")
			c.JSON(http.StatusForbidden, gin.H{
				"error":   "Forbidden",
				"message": "Invalid role format",
			})
			c.Abort()
			return
		}

		// Check if user role is in allowed roles
		for _, allowedRole := range roles {
			if userRoleStr == allowedRole.ToString() {
				logger.Infof("User role %s is authorized", userRoleStr)
				c.Next()
				return
			}
		}

		logger.Warnf("User role %s is not authorized for this resource", userRoleStr)
		c.JSON(http.StatusForbidden, gin.H{
			"error":   "Forbidden",
			"message": "Insufficient permissions",
		})
		c.Abort()
	}
}
