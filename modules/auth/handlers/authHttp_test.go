package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"template-golang/config"
	"template-golang/modules/auth/usecases"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Mock JWT Usecase
type mockJWTUsecase struct {
	mock.Mock
}

func (m *mockJWTUsecase) GenerateJWT(userID string) (string, error) {
	args := m.Called(userID)
	return args.String(0), args.Error(1)
}

func setupTestRouter(jwtUsecase usecases.JWTUsecase) (*gin.Engine, *authHttpHandler) {
	gin.SetMode(gin.TestMode)
	router := gin.New()

	conf := &config.Config{
		Auth: config.AuthConfig{
			LineClientID:     "test-client-id",
			LineClientSecret: "test-client-secret",
			LineCallbackURL:  "http://localhost:8080/auth/line/callback",
		},
	}

	handler := Provide(jwtUsecase, conf)
	group := router.Group("")
	handler.Routes(group)

	return router, handler
}

func TestLogin(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name           string
		path           string
		expectedStatus int
		expectedBody   map[string]interface{}
	}{
		{
			name:           "Success - should redirect to provider auth when provider is valid",
			path:           "/auth/line/login",
			expectedStatus: http.StatusTemporaryRedirect, // 307
		},
		{
			name:           "Error - should return 400 when provider is missing",
			path:           "/auth//login",
			expectedStatus: http.StatusBadRequest,
			expectedBody: map[string]interface{}{
				"message": "Provider is required",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockJWT := new(mockJWTUsecase)
			router, _ := setupTestRouter(mockJWT)

			w := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodGet, tt.path, nil)
			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)

			if tt.expectedBody != nil {
				var responseBody map[string]interface{}
				json.Unmarshal(w.Body.Bytes(), &responseBody)
				assert.Equal(t, tt.expectedBody, responseBody)
			}
		})
	}
}

func TestAuthCallback(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name           string
		path           string
		mockError      error
		mockToken      string
		setupMock      func(*mockJWTUsecase)
		expectedStatus int
		expectedBody   map[string]interface{}
	}{
		{
			name:           "Error - should return 400 when provider is missing",
			path:           "/auth//callback",
			expectedStatus: http.StatusBadRequest,
			expectedBody: map[string]interface{}{
				"message": "Provider is required",
			},
		},
		{
			name:           "Error - should return 401 when authentication fails",
			path:           "/auth/line/callback",
			mockError:      errors.New("could not find a matching session for this request"),
			expectedStatus: http.StatusUnauthorized,
			expectedBody: map[string]interface{}{
				"error": "could not find a matching session for this request",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockJWT := new(mockJWTUsecase)
			if tt.setupMock != nil {
				tt.setupMock(mockJWT)
			}
			router, _ := setupTestRouter(mockJWT)

			w := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodGet, tt.path, nil)
			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)

			if tt.expectedBody != nil {
				var responseBody map[string]interface{}
				json.Unmarshal(w.Body.Bytes(), &responseBody)
				assert.Equal(t, tt.expectedBody, responseBody)
			}
		})
	}
}

func TestLogout(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name           string
		path           string
		expectedStatus int
		expectedBody   map[string]interface{}
	}{
		{
			name:           "Success - should logout successfully when provider is valid",
			path:           "/auth/line/logout",
			expectedStatus: http.StatusOK,
			expectedBody: map[string]interface{}{
				"message": "logged out",
			},
		},
		{
			name:           "Error - should return 400 when provider is missing",
			path:           "/auth//logout",
			expectedStatus: http.StatusBadRequest,
			expectedBody: map[string]interface{}{
				"message": "Provider is required",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockJWT := new(mockJWTUsecase)
			router, _ := setupTestRouter(mockJWT)

			w := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodGet, tt.path, nil)
			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)

			var responseBody map[string]interface{}
			json.Unmarshal(w.Body.Bytes(), &responseBody)
			assert.Equal(t, tt.expectedBody, responseBody)
		})
	}
}
func (m *mockJWTUsecase) ValidateJWT(token string) (string, error) {
	args := m.Called(token)
	return args.String(0), args.Error(1)
}

func TestInformation(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name           string
		authHeader     string
		mockUserID     string
		mockError      error
		setupMock      func(*mockJWTUsecase)
		expectedStatus int
		expectedBody   map[string]interface{}
	}{
		{
			name:       "Success - should return user info when JWT is valid",
			authHeader: "Bearer valid-token",
			mockUserID: "user123",
			setupMock: func(m *mockJWTUsecase) {
				m.On("ValidateJWT", "valid-token").Return("user123", nil)
			},
			expectedStatus: http.StatusOK,
			expectedBody: map[string]interface{}{
				"user_id": "user123",
				"message": "User authenticated successfully",
			},
		},
		{
			name:       "Success - should handle token without Bearer prefix",
			authHeader: "direct-token",
			mockUserID: "user456",
			setupMock: func(m *mockJWTUsecase) {
				m.On("ValidateJWT", "direct-token").Return("user456", nil)
			},
			expectedStatus: http.StatusOK,
			expectedBody: map[string]interface{}{
				"user_id": "user456",
				"message": "User authenticated successfully",
			},
		},
		{
			name:           "Error - should return 401 when Authorization header is missing",
			authHeader:     "",
			expectedStatus: http.StatusUnauthorized,
			expectedBody: map[string]interface{}{
				"error": "Authorization header is required",
			},
		},
		{
			name:       "Error - should return 401 when JWT validation fails",
			authHeader: "Bearer invalid-token",
			mockError:  errors.New("token expired"),
			setupMock: func(m *mockJWTUsecase) {
				m.On("ValidateJWT", "invalid-token").Return("", errors.New("token expired"))
			},
			expectedStatus: http.StatusUnauthorized,
			expectedBody: map[string]interface{}{
				"error": "Invalid or expired token",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockJWT := new(mockJWTUsecase)
			if tt.setupMock != nil {
				tt.setupMock(mockJWT)
			}
			router, _ := setupTestRouter(mockJWT)

			w := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodGet, "/auth/info", nil)
			if tt.authHeader != "" {
				req.Header.Set("Authorization", tt.authHeader)
			}
			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)

			var responseBody map[string]interface{}
			json.Unmarshal(w.Body.Bytes(), &responseBody)
			assert.Equal(t, tt.expectedBody, responseBody)

			mockJWT.AssertExpectations(t)
		})
	}
}

func TestAuthCallback_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name           string
		path           string
		mockToken      string
		mockError      error
		setupMock      func(*mockJWTUsecase)
		expectedStatus int
		checkLocation  bool
	}{
		{
			name:      "Success - should generate JWT and redirect when authentication succeeds",
			path:      "/auth/line/callback",
			mockToken: "generated-jwt-token",
			setupMock: func(m *mockJWTUsecase) {
				m.On("GenerateJWT", mock.AnythingOfType("string")).Return("generated-jwt-token", nil)
			},
			expectedStatus: http.StatusFound,
			checkLocation:  true,
		},
		{
			name:      "Error - should return 500 when JWT generation fails",
			path:      "/auth/line/callback",
			mockError: errors.New("JWT generation failed"),
			setupMock: func(m *mockJWTUsecase) {
				m.On("GenerateJWT", mock.AnythingOfType("string")).Return("", errors.New("JWT generation failed"))
			},
			expectedStatus: http.StatusInternalServerError,
			checkLocation:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockJWT := new(mockJWTUsecase)
			if tt.setupMock != nil {
				tt.setupMock(mockJWT)
			}
			router, handler := setupTestRouter(mockJWT)

			w := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodGet, tt.path, nil)
			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)

			if tt.checkLocation && tt.expectedStatus == http.StatusFound {
				location := w.Header().Get("Location")
				assert.Contains(t, location, handler.conf.Auth.LineFECallbackURL)
				if tt.mockToken != "" {
					assert.Contains(t, location, "token="+tt.mockToken)
				}
			}

			if tt.expectedStatus == http.StatusInternalServerError {
				var responseBody map[string]interface{}
				json.Unmarshal(w.Body.Bytes(), &responseBody)
				assert.Equal(t, "Failed to generate token", responseBody["error"])
			}

			mockJWT.AssertExpectations(t)
		})
	}
}
