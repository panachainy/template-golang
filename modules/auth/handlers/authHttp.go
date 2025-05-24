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

	// if gothUser, err := gothic.CompleteUserAuth(c.Writer, c.Request); err == nil {
	// 	t, _ := template.New("foo").Parse(userTemplate)
	// 	t.Execute(res, gothUser)
	// } else {

	gothic.BeginAuthHandler(c.Writer, c.Request)
	// }

	// ===

	// if err := h.authUsecase.ProcessLogin(reqBody); err != nil {
	// 	c.JSON(http.StatusUnauthorized, gin.H{"message": "Authentication failed"})
	// 	c.Error(err)
	// 	return
	// }

	// c.JSON(http.StatusOK, gin.H{"message": "Login successful"})
	// return
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

	// TODO: add redirect url to FE page such as home page
	// Redirect with the token as a query parameter
	redirectURL := "http://localhost:3000/auth/callback?token=" + token
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

func (h *authHttpHandler) Routes(routerGroup *gin.RouterGroup) {
	authGroup := routerGroup.Group("/auth/:provider")

	authGroup.GET("/login", h.Login)
	authGroup.GET("/callback", h.AuthCallback)
	authGroup.GET("/logout", h.Logout)
}
