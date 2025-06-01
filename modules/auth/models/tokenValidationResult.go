package models

import "github.com/golang-jwt/jwt/v5"

// TokenValidationResult represents the result of token validation
type TokenValidationResult struct {
	Valid    bool
	Expired  bool
	NotExist bool
	Claims   jwt.MapClaims
	UserID   string
}
