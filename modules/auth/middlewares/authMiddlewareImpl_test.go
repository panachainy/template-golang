package middlewares

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"template-golang/modules/auth/models"
	"template-golang/modules/auth/usecases"
	"template-golang/modules/auth/usecases/mock"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func setupTestMiddleware(jwtUsecase usecases.JWTUsecase) (*gin.Engine, gin.HandlerFunc) {
	gin.SetMode(gin.TestMode)
	router := gin.New()

	middleware := NewAuthMiddleware(jwtUsecase)
	authMiddleware := middleware.Handle()

	// Create a test route that uses the middleware
	router.GET("/protected", authMiddleware, func(c *gin.Context) {
		userID := c.GetString("userID")
		c.JSON(http.StatusOK, gin.H{
			"message": "Success",
			"userID":  userID,
		})
	})

	return router, authMiddleware
}

func TestAuthMiddleware_ValidToken(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockJWT := mock.NewMockJWTUsecase(ctrl)
	router, _ := setupTestMiddleware(mockJWT)

	// Mock successful token validation
	mockResult := &models.TokenValidationResult{
		Valid:    true,
		Expired:  false,
		NotExist: false,
		Claims:   jwt.MapClaims{"sub": "test-user-123"},
		UserID:   "test-user-123",
	}
	mockJWT.EXPECT().ValidateJWT("valid-token").Return(mockResult, nil)

	// Create request with valid Bearer token
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/protected", nil)
	req.Header.Set("Authorization", "Bearer valid-token")

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	_ = json.Unmarshal(w.Body.Bytes(), &response)
	assert.Equal(t, "Success", response["message"])
	assert.Equal(t, "test-user-123", response["userID"])
}

func TestAuthMiddleware_MissingAuthorizationHeader(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockJWT := mock.NewMockJWTUsecase(ctrl)
	router, _ := setupTestMiddleware(mockJWT)

	// Create request without Authorization header
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/protected", nil)

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)

	var response map[string]interface{}
	_ = json.Unmarshal(w.Body.Bytes(), &response)
	assert.Equal(t, "Unauthorized", response["error"])
	assert.Equal(t, "Missing authorization header", response["message"])
}

func TestAuthMiddleware_InvalidAuthorizationFormat(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockJWT := mock.NewMockJWTUsecase(ctrl)
	router, _ := setupTestMiddleware(mockJWT)

	tests := []struct {
		name        string
		authHeader  string
		expectedMsg string
	}{
		{
			name:        "Missing Bearer prefix",
			authHeader:  "invalid-token",
			expectedMsg: "Invalid authorization header format",
		},
		{
			name:        "Wrong prefix",
			authHeader:  "Basic invalid-token",
			expectedMsg: "Invalid authorization header format",
		},
		{
			name:        "Empty token",
			authHeader:  "Bearer ",
			expectedMsg: "Invalid authorization header format",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodGet, "/protected", nil)
			req.Header.Set("Authorization", tt.authHeader)

			router.ServeHTTP(w, req)

			assert.Equal(t, http.StatusUnauthorized, w.Code)

			var response map[string]interface{}
			_ = json.Unmarshal(w.Body.Bytes(), &response)
			assert.Equal(t, "Unauthorized", response["error"])
			assert.Equal(t, tt.expectedMsg, response["message"])
		})
	}
}

func TestAuthMiddleware_ExpiredToken(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockJWT := mock.NewMockJWTUsecase(ctrl)
	router, _ := setupTestMiddleware(mockJWT)

	// Mock expired token validation
	mockResult := &models.TokenValidationResult{
		Valid:    false,
		Expired:  true,
		NotExist: false,
		Claims:   nil,
		UserID:   "",
	}
	mockJWT.EXPECT().ValidateJWT("expired-token").Return(mockResult, nil)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/protected", nil)
	req.Header.Set("Authorization", "Bearer expired-token")

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)

	var response map[string]interface{}
	_ = json.Unmarshal(w.Body.Bytes(), &response)
	assert.Equal(t, "Unauthorized", response["error"])
	assert.Equal(t, "Token has expired", response["message"])
}

func TestAuthMiddleware_InvalidToken(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockJWT := mock.NewMockJWTUsecase(ctrl)
	router, _ := setupTestMiddleware(mockJWT)

	// Mock invalid token validation
	mockResult := &models.TokenValidationResult{
		Valid:    false,
		Expired:  false,
		NotExist: false,
		Claims:   nil,
		UserID:   "",
	}
	mockJWT.EXPECT().ValidateJWT("invalid-token").Return(mockResult, nil)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/protected", nil)
	req.Header.Set("Authorization", "Bearer invalid-token")

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)

	var response map[string]interface{}
	_ = json.Unmarshal(w.Body.Bytes(), &response)
	assert.Equal(t, "Unauthorized", response["error"])
	assert.Equal(t, "Invalid token", response["message"])
}
