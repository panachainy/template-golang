package middlewares

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"template-golang/modules/auth/usecases"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Mock JWT Usecase for testing
type mockJWTUsecase struct {
	mock.Mock
}

func (m *mockJWTUsecase) GenerateJWT(userID string) (string, error) {
	args := m.Called(userID)
	return args.String(0), args.Error(1)
}

func (m *mockJWTUsecase) VerifyToken(tokenString string) (*usecases.TokenValidationResult, error) {
	args := m.Called(tokenString)
	return args.Get(0).(*usecases.TokenValidationResult), args.Error(1)
}

func setupTestMiddleware(jwtUsecase usecases.JWTUsecase) (*gin.Engine, gin.HandlerFunc) {
	gin.SetMode(gin.TestMode)
	router := gin.New()

	middleware := Provide(jwtUsecase)
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
	mockJWT := new(mockJWTUsecase)
	router, _ := setupTestMiddleware(mockJWT)

	// Mock successful token verification
	mockResult := &usecases.TokenValidationResult{
		Valid:    true,
		Expired:  false,
		NotExist: false,
		Claims:   jwt.MapClaims{"sub": "test-user-123"},
		UserID:   "test-user-123",
	}
	mockJWT.On("VerifyToken", "valid-token").Return(mockResult, nil)

	// Create request with valid Bearer token
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/protected", nil)
	req.Header.Set("Authorization", "Bearer valid-token")

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Equal(t, "Success", response["message"])
	assert.Equal(t, "test-user-123", response["userID"])

	mockJWT.AssertExpectations(t)
}

func TestAuthMiddleware_MissingAuthorizationHeader(t *testing.T) {
	mockJWT := new(mockJWTUsecase)
	router, _ := setupTestMiddleware(mockJWT)

	// Create request without Authorization header
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/protected", nil)

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Equal(t, "Unauthorized", response["error"])
	assert.Equal(t, "Missing authorization header", response["message"])
}

func TestAuthMiddleware_InvalidAuthorizationFormat(t *testing.T) {
	mockJWT := new(mockJWTUsecase)
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
			json.Unmarshal(w.Body.Bytes(), &response)
			assert.Equal(t, "Unauthorized", response["error"])
			assert.Equal(t, tt.expectedMsg, response["message"])
		})
	}
}

func TestAuthMiddleware_ExpiredToken(t *testing.T) {
	mockJWT := new(mockJWTUsecase)
	router, _ := setupTestMiddleware(mockJWT)

	// Mock expired token verification
	mockResult := &usecases.TokenValidationResult{
		Valid:    false,
		Expired:  true,
		NotExist: false,
		Claims:   nil,
		UserID:   "",
	}
	mockJWT.On("VerifyToken", "expired-token").Return(mockResult, nil)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/protected", nil)
	req.Header.Set("Authorization", "Bearer expired-token")

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Equal(t, "Unauthorized", response["error"])
	assert.Equal(t, "Token has expired", response["message"])

	mockJWT.AssertExpectations(t)
}

func TestAuthMiddleware_InvalidToken(t *testing.T) {
	mockJWT := new(mockJWTUsecase)
	router, _ := setupTestMiddleware(mockJWT)

	// Mock invalid token verification
	mockResult := &usecases.TokenValidationResult{
		Valid:    false,
		Expired:  false,
		NotExist: false,
		Claims:   nil,
		UserID:   "",
	}
	mockJWT.On("VerifyToken", "invalid-token").Return(mockResult, nil)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/protected", nil)
	req.Header.Set("Authorization", "Bearer invalid-token")

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Equal(t, "Unauthorized", response["error"])
	assert.Equal(t, "Invalid token", response["message"])

	mockJWT.AssertExpectations(t)
}
