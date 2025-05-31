//go:generate mockgen -source=jwtUsecase.go -destination=../../../mock/mock_jwt_usecase.go -package=mock

package usecases

import "github.com/golang-jwt/jwt/v5"

// TokenValidationResult represents the result of token validation
type TokenValidationResult struct {
	Valid    bool
	Expired  bool
	NotExist bool
	Claims   jwt.MapClaims
	UserID   string
}

type JWTUsecase interface {
	GenerateJWT(userID string) (string, error)
	ValidateJWT(tokenString string) (*TokenValidationResult, error)
}
