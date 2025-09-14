package usecases

import (
	"template-golang/config"
	"testing"

	"github.com/stretchr/testify/assert"
)

func setupJWTUsecase(t *testing.T) JWTUsecase {
	conf := &config.Config{
		Auth: config.AuthConfig{
			PrivateKeyPath: "../../../config/ecdsa_private_key_test.pem",
		},
	}

	return NewJWTUsecase(conf, nil)
}

func TestGenerateJWT(t *testing.T) {
	jwtUsecase := setupJWTUsecase(t)

	userID := "test-user-123"
	token, err := jwtUsecase.GenerateJWT(userID)

	assert.NoError(t, err)
	assert.NotEmpty(t, token)
}

func TestValidateJWT_ValidToken(t *testing.T) {
	jwtUsecase := setupJWTUsecase(t)

	// Generate a valid token first
	userID := "test-user-123"
	token, err := jwtUsecase.GenerateJWT(userID)
	assert.NoError(t, err)

	// Validate the token
	result, err := jwtUsecase.ValidateJWT(token)
	assert.NoError(t, err)
	assert.True(t, result.Valid)
	assert.False(t, result.Expired)
	assert.False(t, result.NotExist)
}

func TestValidateJWT_InvalidToken(t *testing.T) {
	jwtUsecase := setupJWTUsecase(t)

	// Test with invalid token
	_, err := jwtUsecase.ValidateJWT("invalid-token")
	// The implementation returns an error for malformed tokens
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to parse token")
}
