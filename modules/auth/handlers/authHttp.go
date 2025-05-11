package handlers

import (
	"template-golang/modules/auth/usecases"

	"github.com/gin-gonic/gin"
)

type authHttpHandler struct {
	authUsecase usecases.AuthUsecase
}

func Provide(authUsecase usecases.AuthUsecase) *authHttpHandler {
	return &authHttpHandler{
		authUsecase: authUsecase,
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

	// if err := h.authUsecase.ProcessLogin(reqBody); err != nil {
	// 	c.JSON(http.StatusUnauthorized, gin.H{"message": "Authentication failed"})
	// 	c.Error(err)
	// 	return
	// }

	// c.JSON(http.StatusOK, gin.H{"message": "Login successful"})
	// return
}
