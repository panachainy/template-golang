package context

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func setupTestGin() (*gin.Engine, *gin.Context, *httptest.ResponseRecorder) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	req := httptest.NewRequest("GET", "/test", nil)
	c.Request = req

	router := gin.New()
	return router, c, w
}

func TestNewRequestContext(t *testing.T) {
	_, c, _ := setupTestGin()

	// Set some values in gin context
	c.Set("request_id", "test-request-123")
	c.Set("userID", "user-456")
	c.Set("claims", map[string]interface{}{"role": "admin"})

	rc := NewRequestContext(c)

	assert.NotNil(t, rc)
	assert.NotNil(t, rc.Context)
	assert.NotNil(t, rc.ctx)

	// Check that values were copied
	assert.Equal(t, "test-request-123", rc.GetRequestID())
	assert.Equal(t, "user-456", rc.GetUserID())
	assert.NotNil(t, rc.GetUserClaims())
	assert.NotEmpty(t, rc.GetIPAddress())
	assert.NotZero(t, rc.GetStartTime())
}

func TestRequestContext_WithValue(t *testing.T) {
	_, c, _ := setupTestGin()
	rc := NewRequestContext(c)

	rc.WithValue(TraceIDKey, "trace-123")

	assert.Equal(t, "trace-123", rc.GetTraceID())
}

func TestRequestContext_WithTimeout(t *testing.T) {
	_, c, _ := setupTestGin()
	rc := NewRequestContext(c)

	timeout := 5 * time.Second
	newRc, cancel := rc.WithTimeout(timeout)
	defer cancel()

	assert.NotNil(t, newRc)
	assert.NotEqual(t, rc.ctx, newRc.ctx)

	deadline, ok := newRc.ctx.Deadline()
	assert.True(t, ok)
	assert.True(t, deadline.After(time.Now()))
}

func TestRequestContext_WithDeadline(t *testing.T) {
	_, c, _ := setupTestGin()
	rc := NewRequestContext(c)

	deadline := time.Now().Add(10 * time.Second)
	newRc, cancel := rc.WithDeadline(deadline)
	defer cancel()

	assert.NotNil(t, newRc)

	ctxDeadline, ok := newRc.ctx.Deadline()
	assert.True(t, ok)
	assert.Equal(t, deadline.Unix(), ctxDeadline.Unix())
}

func TestRequestContext_WithCancel(t *testing.T) {
	_, c, _ := setupTestGin()
	rc := NewRequestContext(c)

	newRc, cancel := rc.WithCancel()
	defer cancel()

	assert.NotNil(t, newRc)
	assert.NotEqual(t, rc.ctx, newRc.ctx)

	// Test cancellation
	cancel()

	select {
	case <-newRc.ctx.Done():
		assert.Equal(t, context.Canceled, newRc.ctx.Err())
	default:
		t.Error("Context should be cancelled")
	}
}

func TestRequestContext_Getters(t *testing.T) {
	_, c, _ := setupTestGin()
	rc := NewRequestContext(c)

	// Test string getters with empty values
	assert.Empty(t, rc.GetRequestID())
	assert.Empty(t, rc.GetUserID())
	assert.Empty(t, rc.GetTraceID())
	assert.Empty(t, rc.GetUserRole())
	assert.Nil(t, rc.GetUserClaims())

	// Test non-empty start time
	assert.False(t, rc.GetStartTime().IsZero())
	assert.NotEmpty(t, rc.GetIPAddress())

	// Test elapsed time
	time.Sleep(1 * time.Millisecond)
	elapsed := rc.GetElapsedTime()
	assert.True(t, elapsed > 0)
}

func TestRequestContext_Setters(t *testing.T) {
	_, c, _ := setupTestGin()
	rc := NewRequestContext(c)

	// Set values
	rc.SetRequestID("req-123")
	rc.SetUserID("user-456")
	rc.SetTraceID("trace-789")
	rc.SetUserRole("admin")
	rc.SetUserClaims(map[string]interface{}{"role": "admin"})

	// Verify values in RequestContext
	assert.Equal(t, "req-123", rc.GetRequestID())
	assert.Equal(t, "user-456", rc.GetUserID())
	assert.Equal(t, "trace-789", rc.GetTraceID())
	assert.Equal(t, "admin", rc.GetUserRole())
	assert.NotNil(t, rc.GetUserClaims())

	// Verify values were also set in gin context
	assert.Equal(t, "req-123", c.GetString("request_id"))
	assert.Equal(t, "user-456", c.GetString("userID"))
	assert.Equal(t, "trace-789", c.GetString("trace_id"))
	assert.Equal(t, "admin", c.GetString("user_role"))

	claims, exists := c.Get("claims")
	assert.True(t, exists)
	assert.NotNil(t, claims)
}

func TestRequestIDMiddleware(t *testing.T) {
	router, _, w := setupTestGin()

	middleware := RequestIDMiddleware()
	router.Use(middleware)

	router.GET("/test", func(c *gin.Context) {
		requestID := GetRequestIDFromGin(c)
		assert.NotEmpty(t, requestID)
		c.JSON(200, gin.H{"request_id": requestID})
	})

	req := httptest.NewRequest("GET", "/test", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
	assert.NotEmpty(t, w.Header().Get("X-Request-ID"))
}

func TestRequestIDMiddleware_WithExistingHeader(t *testing.T) {
	router, _, w := setupTestGin()

	middleware := RequestIDMiddleware()
	router.Use(middleware)

	router.GET("/test", func(c *gin.Context) {
		requestID := GetRequestIDFromGin(c)
		assert.Equal(t, "existing-request-id", requestID)
		c.JSON(200, gin.H{"request_id": requestID})
	})

	req := httptest.NewRequest("GET", "/test", nil)
	req.Header.Set("X-Request-ID", "existing-request-id")
	router.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
	assert.Equal(t, "existing-request-id", w.Header().Get("X-Request-ID"))
}

func TestTraceIDMiddleware(t *testing.T) {
	router, _, w := setupTestGin()

	middleware := TraceIDMiddleware()
	router.Use(middleware)

	router.GET("/test", func(c *gin.Context) {
		traceID := GetTraceIDFromGin(c)
		assert.NotEmpty(t, traceID)
		c.JSON(200, gin.H{"trace_id": traceID})
	})

	req := httptest.NewRequest("GET", "/test", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
	assert.NotEmpty(t, w.Header().Get("X-Trace-ID"))
}

func TestTimeoutMiddleware(t *testing.T) {
	router, _, w := setupTestGin()

	timeout := 100 * time.Millisecond
	middleware := TimeoutMiddleware(timeout)
	router.Use(middleware)

	router.GET("/test", func(c *gin.Context) {
		ctx := c.Request.Context()
		deadline, ok := ctx.Deadline()
		assert.True(t, ok)
		assert.True(t, deadline.After(time.Now()))
		c.JSON(200, gin.H{"status": "ok"})
	})

	req := httptest.NewRequest("GET", "/test", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
}

func TestMetricsMiddleware(t *testing.T) {
	router, _, w := setupTestGin()

	middleware := MetricsMiddleware()
	router.Use(middleware)

	router.GET("/test", func(c *gin.Context) {
		startTime := GetStartTimeFromGin(c)
		assert.False(t, startTime.IsZero())

		elapsed := GetElapsedTimeFromGin(c)
		assert.True(t, elapsed >= 0)

		c.JSON(200, gin.H{"status": "ok"})
	})

	req := httptest.NewRequest("GET", "/test", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
	assert.NotEmpty(t, w.Header().Get("X-Response-Time"))
}

func TestGetRequestContext(t *testing.T) {
	_, c, _ := setupTestGin()

	rc := GetRequestContext(c)
	assert.NotNil(t, rc)
	assert.Equal(t, c, rc.Context)
}

func TestHelperFunctions(t *testing.T) {
	_, c, _ := setupTestGin()

	// Test with empty context
	assert.Empty(t, GetRequestIDFromGin(c))
	assert.Empty(t, GetUserIDFromGin(c))
	assert.Empty(t, GetTraceIDFromGin(c))
	assert.Nil(t, GetUserClaimsFromGin(c))
	assert.True(t, GetStartTimeFromGin(c).IsZero())
	assert.Equal(t, time.Duration(0), GetElapsedTimeFromGin(c))

	// Set values
	c.Set("request_id", "req-123")
	c.Set("userID", "user-456")
	c.Set("trace_id", "trace-789")
	c.Set("claims", map[string]interface{}{"role": "admin"})
	c.Set("start_time", time.Now())

	// Test with values
	assert.Equal(t, "req-123", GetRequestIDFromGin(c))
	assert.Equal(t, "user-456", GetUserIDFromGin(c))
	assert.Equal(t, "trace-789", GetTraceIDFromGin(c))
	assert.NotNil(t, GetUserClaimsFromGin(c))
	assert.False(t, GetStartTimeFromGin(c).IsZero())
	assert.True(t, GetElapsedTimeFromGin(c) >= 0)
}

func TestSetGetContextValue(t *testing.T) {
	_, c, _ := setupTestGin()

	// Set value
	SetContextValue(c, "test_key", "test_value")

	// Get value
	value, exists := GetContextValue(c, "test_key")
	assert.True(t, exists)
	assert.Equal(t, "test_value", value)

	// Get non-existent value
	value, exists = GetContextValue(c, "non_existent")
	assert.False(t, exists)
	assert.Nil(t, value)
}

func TestGenerateFunctions(t *testing.T) {
	requestID := GenerateRequestID()
	assert.NotEmpty(t, requestID)
	assert.Len(t, requestID, 36) // UUID length

	traceID := GenerateTraceID()
	assert.NotEmpty(t, traceID)
	assert.Len(t, traceID, 36) // UUID length

	// Ensure they're different
	assert.NotEqual(t, requestID, traceID)

	// Ensure multiple calls generate different IDs
	requestID2 := GenerateRequestID()
	assert.NotEqual(t, requestID, requestID2)
}

func TestContextKey(t *testing.T) {
	// Test that context keys are defined
	assert.Equal(t, ContextKey("request_id"), RequestIDKey)
	assert.Equal(t, ContextKey("user_id"), UserIDKey)
	assert.Equal(t, ContextKey("trace_id"), TraceIDKey)
	assert.Equal(t, ContextKey("user_role"), UserRoleKey)
	assert.Equal(t, ContextKey("user_claims"), UserClaimsKey)
	assert.Equal(t, ContextKey("start_time"), StartTimeKey)
	assert.Equal(t, ContextKey("ip_address"), IPAddressKey)
	assert.Equal(t, ContextKey("user_agent"), UserAgentKey)
}

func TestGetStringValue_NonString(t *testing.T) {
	_, c, _ := setupTestGin()
	rc := NewRequestContext(c)

	// Set a non-string value
	rc.WithValue(RequestIDKey, 123)

	// Should return empty string for non-string values
	assert.Empty(t, rc.GetStringValue(RequestIDKey))
}
