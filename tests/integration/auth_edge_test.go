package integration

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"template-golang/modules/auth/handlers"
	"template-golang/modules/auth/middlewares"
	"template-golang/modules/auth/repositories"
	"template-golang/modules/auth/usecases"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAuthHandler_EdgeCases_Integration(t *testing.T) {
	// Setup test database
	pool, cleanup := SetupTestDB(t)
	defer cleanup()

	// Wait for database to be ready
	WaitForDB(t, pool, 10*time.Second)

	// Setup test configuration
	conf := SetupTestConfig(t)

	// Create database instance
	queries := CreateTestDatabase(t, pool)

	// Setup dependencies
	authRepo := repositories.NewAuthRepository(queries)
	jwtUsecase := usecases.NewJWTUsecase(conf, authRepo)
	authMiddleware := middlewares.NewAuthMiddleware(jwtUsecase)

	// Create auth handler
	authHandler := handlers.NewAuthHttpHandler(jwtUsecase, conf, authMiddleware, authRepo)

	// Setup Gin router
	gin.SetMode(gin.TestMode)
	router := gin.New()
	api := router.Group("/api/v1")
	authHandler.Routes(api)

	t.Run("Test HTTP Methods", func(t *testing.T) {
		tests := []struct {
			name           string
			method         string
			path           string
			expectedStatus int
		}{
			{
				name:           "POST to login endpoint",
				method:         "POST",
				path:           "/api/v1/auth/line/login",
				expectedStatus: http.StatusNotFound, // Gin returns 404 for unregistered method/route combinations
			},
			{
				name:           "PUT to callback endpoint",
				method:         "PUT",
				path:           "/api/v1/auth/line/callback",
				expectedStatus: http.StatusNotFound,
			},
			{
				name:           "DELETE to logout endpoint",
				method:         "DELETE",
				path:           "/api/v1/auth/line/logout",
				expectedStatus: http.StatusNotFound,
			},
			{
				name:           "POST to example endpoint",
				method:         "POST",
				path:           "/api/v1/auth/example",
				expectedStatus: http.StatusNotFound,
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				req, err := http.NewRequest(tt.method, tt.path, nil)
				require.NoError(t, err)

				w := httptest.NewRecorder()
				router.ServeHTTP(w, req)

				assert.Equal(t, tt.expectedStatus, w.Code)
			})
		}
	})

	t.Run("Test Special Characters in Provider", func(t *testing.T) {
		specialProviders := []string{
			"provider-with-dash",
			"provider_with_underscore",
			"provider.with.dots",
			"provider@with@symbols",
			"provider with spaces",
			"provider/with/slashes",
		}

		for _, provider := range specialProviders {
			t.Run("Provider: "+provider, func(t *testing.T) {
				// URL encode the provider for the request
				encodedProvider := strings.ReplaceAll(provider, " ", "%20")
				encodedProvider = strings.ReplaceAll(encodedProvider, "@", "%40")

				req, err := http.NewRequest("GET", "/api/v1/auth/"+encodedProvider+"/login", nil)
				require.NoError(t, err)

				w := httptest.NewRecorder()
				router.ServeHTTP(w, req)

				// Should either redirect or return an error, but not crash
				assert.True(t, w.Code >= 300, "Expected non-2xx status for special provider: %s", provider)
			})
		}
	})

	t.Run("Test Malformed Authorization Headers", func(t *testing.T) {
		malformedHeaders := []string{
			"Bearer",                              // Missing token
			"Bearer ",                             // Empty token
			"BasicAuthNotBearer xyz",              // Wrong auth type
			"Bearer token.with.only.two.parts",    // Invalid JWT format
			"Bearer " + strings.Repeat("x", 1000), // Very long token
		}

		for _, header := range malformedHeaders {
			t.Run("Header: "+header[:min(len(header), 50)], func(t *testing.T) {
				req, err := http.NewRequest("GET", "/api/v1/auth/example", nil)
				require.NoError(t, err)

				req.Header.Set("Authorization", header)

				w := httptest.NewRecorder()
				router.ServeHTTP(w, req)

				// Should return unauthorized for malformed auth headers
				assert.Equal(t, http.StatusUnauthorized, w.Code)
			})
		}
	})

	t.Run("Test Very Long Provider Names", func(t *testing.T) {
		longProvider := strings.Repeat("a", 1000)

		req, err := http.NewRequest("GET", "/api/v1/auth/"+longProvider+"/login", nil)
		require.NoError(t, err)

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		// Should handle long provider names gracefully
		assert.True(t, w.Code >= 300, "Expected non-2xx status for very long provider name")
	})
}

func TestAuthHandler_ConcurrentRequests_Integration(t *testing.T) {
	// Setup test database
	pool, cleanup := SetupTestDB(t)
	defer cleanup()

	// Wait for database to be ready
	WaitForDB(t, pool, 10*time.Second)

	// Setup test configuration
	conf := SetupTestConfig(t)

	// Create database instance
	queries := CreateTestDatabase(t, pool)

	// Setup dependencies
	authRepo := repositories.NewAuthRepository(queries)
	jwtUsecase := usecases.NewJWTUsecase(conf, authRepo)
	authMiddleware := middlewares.NewAuthMiddleware(jwtUsecase)

	// Create auth handler
	authHandler := handlers.NewAuthHttpHandler(jwtUsecase, conf, authMiddleware, authRepo)

	// Setup Gin router
	gin.SetMode(gin.TestMode)
	router := gin.New()
	api := router.Group("/api/v1")
	authHandler.Routes(api)

	t.Run("Concurrent login requests", func(t *testing.T) {
		const numGoroutines = 10
		done := make(chan bool, numGoroutines)

		for i := 0; i < numGoroutines; i++ {
			go func(id int) {
				defer func() { done <- true }()

				req, err := http.NewRequest("GET", "/api/v1/auth/line/login", nil)
				require.NoError(t, err)

				w := httptest.NewRecorder()
				router.ServeHTTP(w, req)

				// Each request should get a consistent response
				assert.True(t, w.Code >= 300, "Request %d: Expected redirect or error", id)
			}(i)
		}

		// Wait for all goroutines to complete
		for i := 0; i < numGoroutines; i++ {
			<-done
		}
	})
}

// Helper function for Go < 1.21 compatibility
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
