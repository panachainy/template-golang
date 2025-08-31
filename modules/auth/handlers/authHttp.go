package handlers

import (
	"net/http"
	"strconv"
	"template-golang/config"
	"template-golang/modules/auth/middlewares"
	"template-golang/modules/auth/models"
	"template-golang/modules/auth/repositories"
	"template-golang/modules/auth/usecases"

	"github.com/gin-gonic/gin"

	"github.com/markbates/goth"
	"github.com/markbates/goth/gothic"
	"github.com/markbates/goth/providers/line"
)

// TODO: fix goth/gothic: no SESSION_SECRET environment variable is set. The default cookie store is not available and any calls will fail. Ignore this warning if you are using a different store.

type authHttpHandler struct {
	jwtUsecase     usecases.JWTUsecase
	conf           *config.Config
	authMiddleware middlewares.AuthMiddleware
	authRepo       repositories.AuthRepository
}

func NewAuthHttpHandler(jwtUsecase usecases.JWTUsecase, conf *config.Config,
	authMiddleware middlewares.AuthMiddleware, authRepo repositories.AuthRepository) AuthHandler {
	goth.UseProviders(
		line.New(conf.Auth.LineClientID, conf.Auth.LineClientSecret, conf.Auth.LineCallbackURL, "profile", "openid", "email"),
	)

	return &authHttpHandler{
		jwtUsecase:     jwtUsecase,
		conf:           conf,
		authMiddleware: authMiddleware,
		authRepo:       authRepo,
	}
}

func (h *authHttpHandler) Login(c *gin.Context) {
	// Translate provider
	provider := c.Param("provider")
	if provider == "" {
		c.JSON(400, gin.H{"message": "Provider is required"})
		return
	}

	q := c.Request.URL.Query()
	q.Add("provider", c.Param("provider"))
	c.Request.URL.RawQuery = q.Encode()

	gothic.BeginAuthHandler(c.Writer, c.Request)
}

func (h *authHttpHandler) AuthCallback(c *gin.Context) {
	// Translate provider
	provider := c.Param("provider")
	if provider == "" {
		c.JSON(400, gin.H{"message": "Provider is required"})
		return
	}

	q := c.Request.URL.Query()
	q.Add("provider", c.Param("provider"))
	c.Request.URL.RawQuery = q.Encode()

	user, err := gothic.CompleteUserAuth(c.Writer, c.Request)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	// Insert or update user in the database
	err = h.jwtUsecase.UpsertUser(user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to upsert user"})
		return
	}
	// // Retrieve the user from the database
	// user, err = h.jwtUsecase.GetUserByID(user.UserID)
	// if err != nil {
	// 	c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve user"})
	// 	return
	// }
	// // If user is not found, return unauthorized
	// if user == nil {
	// 	c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found"})
	// 	return
	// }
	// // If user is not active, return unauthorized
	// if !user.IsActive {
	// 	c.JSON(http.StatusUnauthorized, gin.H{"error": "User is not active"})
	// 	return
	// }

	// Generate JWT for the authenticated user
	token, err := h.jwtUsecase.GenerateJWT(user.UserID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	// Redirect with the token as a query parameter
	redirectURL := h.conf.Auth.LineFECallbackURL + "?token=" + token
	c.Redirect(http.StatusFound, redirectURL)
	return
}

func (h *authHttpHandler) Logout(c *gin.Context) {
	// Translate provider
	provider := c.Param("provider")
	if provider == "" {
		c.JSON(400, gin.H{"message": "Provider is required"})
		return
	}

	q := c.Request.URL.Query()
	q.Add("provider", c.Param("provider"))
	c.Request.URL.RawQuery = q.Encode()

	err := gothic.Logout(c.Writer, c.Request)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "logged out"})
}

func (h *authHttpHandler) Example(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "example",
	})
}

// ======================== Admin Routes ========================

// GetUsers retrieves multiple users with pagination
func (h *authHttpHandler) GetUsers(c *gin.Context) {
	limitStr := c.DefaultQuery("limit", "10")
	limit := 10
	if l, err := strconv.Atoi(limitStr); err == nil && l > 0 {
		limit = l
	}

	users, err := h.authRepo.Gets(limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve users"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"users": users,
		"count": len(users),
	})
}

// ======================== Admin Routes ========================

func (h *authHttpHandler) Routes(routerGroup *gin.RouterGroup) {
	authProviderGroup := routerGroup.Group("/auth/:provider")

	authProviderGroup.GET("/login", h.Login)
	authProviderGroup.GET("/callback", h.AuthCallback)
	authProviderGroup.GET("/logout", h.Logout)

	authGroup := routerGroup.Group("/auth")
	authGroup.Use(h.authMiddleware.Handle())
	authGroup.GET("/example", h.Example)

	authAdminGroup := routerGroup.Group("/admin/auth")
	authAdminGroup.Use(h.authMiddleware.Handle(),
		h.authMiddleware.Allows(
			[]models.Role{models.RoleAdmin}),
	)
	authAdminGroup.GET("/users", h.GetUsers)
}
