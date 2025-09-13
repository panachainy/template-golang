package usecases

import (
	"template-golang/modules/auth/models"

	"github.com/markbates/goth"
)

type JWTUsecase interface {
	GenerateJWT(userID string) (string, error)
	ValidateJWT(tokenString string) (*models.TokenValidationResult, error)
	UpsertUser(user goth.User, role ...models.Role) error
}
