package handlers

import (
	"net/http"
	"template-golang/config"
	"template-golang/modules/auth/middlewares"
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
}

func Provide(jwtUsecase usecases.JWTUsecase, conf *config.Config, authMiddleware middlewares.AuthMiddleware) *authHttpHandler {
	goth.UseProviders(
		line.New(conf.Auth.LineClientID, conf.Auth.LineClientSecret, conf.Auth.LineCallbackURL, "profile", "openid", "email"),
	)

	return &authHttpHandler{
		jwtUsecase:     jwtUsecase,
		conf:           conf,
		authMiddleware: authMiddleware,
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

func (h *authHttpHandler) Routes(routerGroup *gin.RouterGroup) {
	authProviderGroup := routerGroup.Group("/auth/:provider")

	authProviderGroup.GET("/login", h.Login)
	authProviderGroup.GET("/callback", h.AuthCallback)
	authProviderGroup.GET("/logout", h.Logout)

	authGroup := routerGroup.Group("/auth")
	authGroup.Use(h.authMiddleware.Handle())
	// TODO: should't call this client should be check with it self.
	authGroup.GET("/example", h.Example)
}
