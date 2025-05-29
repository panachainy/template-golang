package usecases

import (
	"template-golang/config"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
)

func setupJWTUsecase(t *testing.T) *jwtUsecaseImpl {
	conf := &config.Config{
		Auth: config.AuthConfig{
			PrivateKeyPath: "../../../ecdsa_private_key.pem",
		},
	}
	return Provide(conf)
}

func TestGenerateJWT(t *testing.T) {
	jwtUsecase := setupJWTUsecase(t)

	userID := "test-user-123"
	token, err := jwtUsecase.GenerateJWT(userID)

	assert.NoError(t, err)
	assert.NotEmpty(t, token)

	// Verify the token contains the expected user ID
	result, err := jwtUsecase.VerifyToken(token)
	assert.NoError(t, err)
	assert.True(t, result.Valid)
	assert.False(t, result.Expired)
	assert.False(t, result.NotExist)
	assert.Equal(t, userID, result.UserID)
	assert.NotNil(t, result.Claims)
}

func TestVerifyToken_ValidToken(t *testing.T) {
	jwtUsecase := setupJWTUsecase(t)

	// Generate a valid token
	userID := "test-user-456"
	token, err := jwtUsecase.GenerateJWT(userID)
	assert.NoError(t, err)

	// Verify the token
	result, err := jwtUsecase.VerifyToken(token)
	assert.NoError(t, err)
	assert.True(t, result.Valid)
	assert.False(t, result.Expired)
	assert.False(t, result.NotExist)
	assert.Equal(t, userID, result.UserID)

	// Check claims
	assert.NotNil(t, result.Claims)
	assert.Equal(t, userID, result.Claims["sub"])
}

func TestVerifyToken_EmptyToken(t *testing.T) {
	jwtUsecase := setupJWTUsecase(t)

	result, err := jwtUsecase.VerifyToken("")
	assert.NoError(t, err)
	assert.False(t, result.Valid)
	assert.False(t, result.Expired)
	assert.True(t, result.NotExist)
	assert.Empty(t, result.UserID)
	assert.Nil(t, result.Claims)
}

func TestVerifyToken_ExpiredToken(t *testing.T) {
	jwtUsecase := setupJWTUsecase(t)

	// Create an expired token manually
	token := jwt.NewWithClaims(jwt.SigningMethodES256, jwt.MapClaims{
		"sub": "test-user-expired",
		"exp": jwt.NewNumericDate(time.Now().Add(-1 * time.Hour)), // Expired 1 hour ago
		"iss": "my-auth-server-issuer",
	})

	expiredTokenString, err := token.SignedString(jwtUsecase.privateKey)
	assert.NoError(t, err)

	// Verify the expired token
	result, err := jwtUsecase.VerifyToken(expiredTokenString)
	assert.NoError(t, err)
	assert.False(t, result.Valid)
	assert.True(t, result.Expired)
	assert.False(t, result.NotExist)
}

func TestVerifyToken_InvalidToken(t *testing.T) {
	jwtUsecase := setupJWTUsecase(t)

	// Test with malformed token
	result, err := jwtUsecase.VerifyToken("invalid-token")
	assert.Error(t, err)
	assert.False(t, result.Valid)
	assert.False(t, result.Expired)
	assert.False(t, result.NotExist)
}

func TestVerifyToken_InvalidSignature(t *testing.T) {
	jwtUsecase := setupJWTUsecase(t)

	// Create a token with wrong signature (using a different key)
	wrongToken := "eyJhbGciOiJFUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiJ0ZXN0LXVzZXIiLCJleHAiOjk5OTk5OTk5OTksImlzcyI6Im15LWF1dGgtc2VydmVyLWlzc3VlciJ9.wrong_signature"

	result, err := jwtUsecase.VerifyToken(wrongToken)
	assert.Error(t, err)
	assert.False(t, result.Valid)
	assert.False(t, result.Expired)
	assert.False(t, result.NotExist)
}
