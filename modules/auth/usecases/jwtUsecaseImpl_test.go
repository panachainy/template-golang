package usecases

import (
	"template-golang/config"
	"template-golang/mock"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func setupJWTUsecase(t *testing.T, ctrl *gomock.Controller) *jwtUsecaseImpl {
	conf := &config.Config{
		Auth: config.AuthConfig{
			PrivateKeyPath: "../../../config/ecdsa_private_key_test.pem",
		},
	}

	// Create a mock auth repository for testing
	mockRepo := mock.NewMockAuthRepository(ctrl)

	return Provide(conf, mockRepo)
}

func TestGenerateJWT(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	jwtUsecase := setupJWTUsecase(t, ctrl)

	userID := "test-user-123"
	token, err := jwtUsecase.GenerateJWT(userID)

	assert.NoError(t, err)
	assert.NotEmpty(t, token)

	// Verify the token contains the expected user ID
	result, err := jwtUsecase.ValidateJWT(token)
	assert.NoError(t, err)
	assert.True(t, result.Valid)
	assert.False(t, result.Expired)
	assert.False(t, result.NotExist)
	assert.Equal(t, userID, result.UserID)
	assert.NotNil(t, result.Claims)
}

func TestVerifyToken_ValidToken(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	jwtUsecase := setupJWTUsecase(t, ctrl)

	// Generate a valid token
	userID := "test-user-456"
	token, err := jwtUsecase.GenerateJWT(userID)
	assert.NoError(t, err)

	// Verify the token
	result, err := jwtUsecase.ValidateJWT(token)
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
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	jwtUsecase := setupJWTUsecase(t, ctrl)

	result, err := jwtUsecase.ValidateJWT("")
	assert.NoError(t, err)
	assert.False(t, result.Valid)
	assert.False(t, result.Expired)
	assert.True(t, result.NotExist)
	assert.Empty(t, result.UserID)
	assert.Nil(t, result.Claims)
}

func TestVerifyToken_ExpiredToken(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	jwtUsecase := setupJWTUsecase(t, ctrl)

	// Create an expired token manually
	token := jwt.NewWithClaims(jwt.SigningMethodES256, jwt.MapClaims{
		"sub": "test-user-expired",
		"exp": jwt.NewNumericDate(time.Now().Add(-1 * time.Hour)), // Expired 1 hour ago
		"iss": "my-auth-server-issuer",
	})

	expiredTokenString, err := token.SignedString(jwtUsecase.privateKey)
	assert.NoError(t, err)

	// Verify the expired token
	result, err := jwtUsecase.ValidateJWT(expiredTokenString)
	assert.NoError(t, err)
	assert.False(t, result.Valid)
	assert.True(t, result.Expired)
	assert.False(t, result.NotExist)
}

func TestVerifyToken_InvalidToken(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	jwtUsecase := setupJWTUsecase(t, ctrl)

	// Test with malformed token
	result, err := jwtUsecase.ValidateJWT("invalid-token")
	assert.Error(t, err)
	assert.False(t, result.Valid)
	assert.False(t, result.Expired)
	assert.False(t, result.NotExist)
}

func TestVerifyToken_InvalidSignature(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	jwtUsecase := setupJWTUsecase(t, ctrl)

	// Create a token with wrong signature (using a different key)
	wrongToken := "eyJhbGciOiJFUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiJ0ZXN0LXVzZXIiLCJleHAiOjk5OTk5OTk5OTksImlzcyI6Im15LWF1dGgtc2VydmVyLWlzc3VlciJ9.wrong_signature"

	result, err := jwtUsecase.ValidateJWT(wrongToken)
	assert.Error(t, err)
	assert.False(t, result.Valid)
	assert.False(t, result.Expired)
	assert.False(t, result.NotExist)
}
