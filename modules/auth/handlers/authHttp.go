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
}

func Provide(jwtUsecase usecases.JWTUsecase, conf *config.Config) *authHttpHandler {
	goth.UseProviders(
		line.New(conf.Auth.Line.ClientID, conf.Auth.Line.ClientSecret, "http://localhost:3000/auth/line/callback", "profile", "openid", "email"),
	)

	return &authHttpHandler{
		jwtUsecase: jwtUsecase,
	}
}

func (h *authHttpHandler) Login(c *gin.Context) {
	// reqBody := new(models.LoginRequest)

	// if err := c.ShouldBindJSON(reqBody); err != nil {
	// 	c.JSON(
	// 		http.StatusBadRequest,
	// 		gin.H{"message": err.Error()},
	// 	)
	// 	c.Error(err)
	// 	return
	// }

	// validate := validator.New(validator.WithRequiredStructEnabled())

	// if err := validate.Struct(reqBody); err != nil {
	// 	c.JSON(
	// 		http.StatusBadRequest,
	// 		gin.H{"message": err.Error()},
	// 	)
	// 	c.Error(err)
	// 	return
	// }

	// === new

	provider := c.Param("provider")
	if provider == "" {
		c.JSON(400, gin.H{"message": "Provider is required"})
		return
	}

	q := c.Request.URL.Query()
	q.Add("provider", c.Param("provider"))
	c.Request.URL.RawQuery = q.Encode()

	// c.Request = c.Request.WithContext(gothic.WithProvider(c.Request.Context(), provider))

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
	// provider := c.Param("provider")
	// if provider == "" {
	// 	c.JSON(400, gin.H{"message": "Provider is required"})
	// 	return
	// }

	// FIXME: check provider
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

	c.JSON(http.StatusOK, gin.H{
		"message": "Authentication successful",
		"token":   token,
		// FIXME: user not sure it important or not?
		"user": user,
	})
}

func (h *authHttpHandler) Logout(c *gin.Context) {
	// provider := c.Param("provider")
	// if provider == "" {
	// 	c.JSON(400, gin.H{"message": "Provider is required"})
	// 	return
	// }

	// c.Request = c.Request.WithContext(gothic.WithProvider(c.Request.Context(), provider))
	// FIXME: check provider
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
