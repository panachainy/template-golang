package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/markbates/goth/gothic"
)

type authHttpHandler struct {
}

func Provide() *authHttpHandler {
	return &authHttpHandler{}
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
	provider := c.Param("provider")
	if provider == "" {
		c.JSON(400, gin.H{"message": "Provider is required"})
		return
	}

	// FIXME: check provider
	user, err := gothic.CompleteUserAuth(c.Writer, c.Request)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"user": user,
	})
}

func (h *authHttpHandler) Logout(c *gin.Context) {
	provider := c.Param("provider")
	if provider == "" {
		c.JSON(400, gin.H{"message": "Provider is required"})
		return
	}

	// c.Request = c.Request.WithContext(gothic.WithProvider(c.Request.Context(), provider))
	// FIXME: check provider
	err := gothic.Logout(c.Writer, c.Request)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "logged out"})

	if err != nil {
		c.JSON(500, gin.H{"message": "Logout failed", "error": err.Error()})
		return
	}

	c.JSON(200, gin.H{"message": "Logout successful"})
}

func (h *authHttpHandler) Routes(routerGroup *gin.RouterGroup) {
	authGroup := routerGroup.Group("/auth/:provider")
	{
		authGroup.GET("/login", h.Login)
		authGroup.GET("/callback", h.AuthCallback)
		authGroup.GET("/logout", h.Logout)
	}
}
