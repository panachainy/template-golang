package handlers

import (
	"net/http"
	"template-golang/config"
	"template-golang/modules/auth/usecases"

	"github.com/gin-gonic/gin"

	"github.com/markbates/goth"
	"github.com/markbates/goth/gothic"
	"github.com/markbates/goth/providers/line"
)

// TODO: fix goth/gothic: no SESSION_SECRET environment variable is set. The default cookie store is not available and any calls will fail. Ignore this warning if you are using a different store.

type authHttpHandler struct {
	jwtUsecase usecases.JWTUsecase
	conf       *config.Config
}

func Provide(jwtUsecase usecases.JWTUsecase, conf *config.Config) *authHttpHandler {
	goth.UseProviders(
		line.New(conf.Auth.LineClientID, conf.Auth.LineClientSecret, conf.Auth.LineCallbackURL, "profile", "openid", "email"),
	)

	return &authHttpHandler{
		jwtUsecase: jwtUsecase,
		conf:       conf,
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

func (h *authHttpHandler) Information(c *gin.Context) {
	// Extract JWT from Authorization header
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header is required"})
		return
	}

	// Remove "Bearer " prefix if present
	token := authHeader
	if len(authHeader) > 7 && authHeader[:7] == "Bearer " {
		token = authHeader[7:]
	}

	// Validate and parse JWT
	result, err := h.jwtUsecase.ValidateJWT(token)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
		return
	}

	// Check validation result
	if !result.Valid || result.Expired || result.NotExist {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"result":  result,
		"message": "User authenticated successfully",
	})
}

func (h *authHttpHandler) Routes(routerGroup *gin.RouterGroup) {
	authProviderGroup := routerGroup.Group("/auth/:provider")

	authProviderGroup.GET("/login", h.Login)
	authProviderGroup.GET("/callback", h.AuthCallback)
	authProviderGroup.GET("/logout", h.Logout)

	authGroup := routerGroup.Group("/auth")
	// need auth
	// TODO: should't call this client should be check with it self.
	authGroup.GET("/info", h.Information)

}
