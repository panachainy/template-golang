package context

import (
	"context"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// ContextKey represents a key for context values
type ContextKey string

const (
	// RequestIDKey is the key for request ID in context
	RequestIDKey ContextKey = "request_id"
	// UserIDKey is the key for user ID in context
	UserIDKey ContextKey = "user_id"
	// TraceIDKey is the key for trace ID in context
	TraceIDKey ContextKey = "trace_id"
	// UserRoleKey is the key for user role in context
	UserRoleKey ContextKey = "user_role"
	// UserClaimsKey is the key for user claims in context
	UserClaimsKey ContextKey = "user_claims"
	// StartTimeKey is the key for request start time in context
	StartTimeKey ContextKey = "start_time"
	// IPAddressKey is the key for client IP address in context
	IPAddressKey ContextKey = "ip_address"
	// UserAgentKey is the key for user agent in context
	UserAgentKey ContextKey = "user_agent"
)

// RequestContext wraps gin.Context with additional functionality
type RequestContext struct {
	*gin.Context
	ctx context.Context
}

// NewRequestContext creates a new RequestContext
func NewRequestContext(ginCtx *gin.Context) *RequestContext {
	ctx := context.WithValue(context.Background(), StartTimeKey, time.Now())

	// Extract values from gin context and add to context
	if requestID := ginCtx.GetString("request_id"); requestID != "" {
		ctx = context.WithValue(ctx, RequestIDKey, requestID)
	}

	if userID := ginCtx.GetString("userID"); userID != "" {
		ctx = context.WithValue(ctx, UserIDKey, userID)
	}

	if claims, exists := ginCtx.Get("claims"); exists {
		ctx = context.WithValue(ctx, UserClaimsKey, claims)
	}

	// Add IP address
	ctx = context.WithValue(ctx, IPAddressKey, ginCtx.ClientIP())

	// Add User-Agent
	ctx = context.WithValue(ctx, UserAgentKey, ginCtx.GetHeader("User-Agent"))

	return &RequestContext{
		Context: ginCtx,
		ctx:     ctx,
	}
}

// GetContext returns the underlying context
func (rc *RequestContext) GetContext() context.Context {
	return rc.ctx
}

// WithValue adds a value to the context
func (rc *RequestContext) WithValue(key ContextKey, value interface{}) *RequestContext {
	rc.ctx = context.WithValue(rc.ctx, key, value)
	return rc
}

// WithTimeout adds a timeout to the context
func (rc *RequestContext) WithTimeout(timeout time.Duration) (*RequestContext, context.CancelFunc) {
	ctx, cancel := context.WithTimeout(rc.ctx, timeout)
	return &RequestContext{
		Context: rc.Context,
		ctx:     ctx,
	}, cancel
}

// WithDeadline adds a deadline to the context
func (rc *RequestContext) WithDeadline(deadline time.Time) (*RequestContext, context.CancelFunc) {
	ctx, cancel := context.WithDeadline(rc.ctx, deadline)
	return &RequestContext{
		Context: rc.Context,
		ctx:     ctx,
	}, cancel
}

// WithCancel adds cancellation to the context
func (rc *RequestContext) WithCancel() (*RequestContext, context.CancelFunc) {
	ctx, cancel := context.WithCancel(rc.ctx)
	return &RequestContext{
		Context: rc.Context,
		ctx:     ctx,
	}, cancel
}

// GetValue retrieves a value from the context
func (rc *RequestContext) GetValue(key ContextKey) interface{} {
	return rc.ctx.Value(key)
}

// GetStringValue retrieves a string value from the context
func (rc *RequestContext) GetStringValue(key ContextKey) string {
	if value := rc.ctx.Value(key); value != nil {
		if str, ok := value.(string); ok {
			return str
		}
	}
	return ""
}

// GetRequestID returns the request ID from context
func (rc *RequestContext) GetRequestID() string {
	return rc.GetStringValue(RequestIDKey)
}

// GetUserID returns the user ID from context
func (rc *RequestContext) GetUserID() string {
	return rc.GetStringValue(UserIDKey)
}

// GetTraceID returns the trace ID from context
func (rc *RequestContext) GetTraceID() string {
	return rc.GetStringValue(TraceIDKey)
}

// GetUserRole returns the user role from context
func (rc *RequestContext) GetUserRole() string {
	return rc.GetStringValue(UserRoleKey)
}

// GetUserClaims returns the user claims from context
func (rc *RequestContext) GetUserClaims() interface{} {
	return rc.GetValue(UserClaimsKey)
}

// GetStartTime returns the request start time from context
func (rc *RequestContext) GetStartTime() time.Time {
	if value := rc.ctx.Value(StartTimeKey); value != nil {
		if t, ok := value.(time.Time); ok {
			return t
		}
	}
	return time.Time{}
}

// GetIPAddress returns the client IP address from context
func (rc *RequestContext) GetIPAddress() string {
	return rc.GetStringValue(IPAddressKey)
}

// GetUserAgent returns the user agent from context
func (rc *RequestContext) GetUserAgent() string {
	return rc.GetStringValue(UserAgentKey)
}

// GetElapsedTime returns the elapsed time since request start
func (rc *RequestContext) GetElapsedTime() time.Duration {
	startTime := rc.GetStartTime()
	if startTime.IsZero() {
		return 0
	}
	return time.Since(startTime)
}

// SetRequestID sets the request ID in both contexts
func (rc *RequestContext) SetRequestID(requestID string) {
	rc.ctx = context.WithValue(rc.ctx, RequestIDKey, requestID)
	rc.Context.Set("request_id", requestID)
}

// SetUserID sets the user ID in both contexts
func (rc *RequestContext) SetUserID(userID string) {
	rc.ctx = context.WithValue(rc.ctx, UserIDKey, userID)
	rc.Context.Set("userID", userID)
}

// SetTraceID sets the trace ID in both contexts
func (rc *RequestContext) SetTraceID(traceID string) {
	rc.ctx = context.WithValue(rc.ctx, TraceIDKey, traceID)
	rc.Context.Set("trace_id", traceID)
}

// SetUserRole sets the user role in both contexts
func (rc *RequestContext) SetUserRole(role string) {
	rc.ctx = context.WithValue(rc.ctx, UserRoleKey, role)
	rc.Context.Set("user_role", role)
}

// SetUserClaims sets the user claims in both contexts
func (rc *RequestContext) SetUserClaims(claims interface{}) {
	rc.ctx = context.WithValue(rc.ctx, UserClaimsKey, claims)
	rc.Context.Set("claims", claims)
}

// Middleware functions

// RequestIDMiddleware adds a unique request ID to each request
func RequestIDMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		requestID := c.GetHeader("X-Request-ID")
		if requestID == "" {
			requestID = uuid.New().String()
		}

		c.Set("request_id", requestID)
		c.Header("X-Request-ID", requestID)

		c.Next()
	}
}

// TraceIDMiddleware adds a trace ID to each request for distributed tracing
func TraceIDMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		traceID := c.GetHeader("X-Trace-ID")
		if traceID == "" {
			traceID = uuid.New().String()
		}

		c.Set("trace_id", traceID)
		c.Header("X-Trace-ID", traceID)

		c.Next()
	}
}

// TimeoutMiddleware adds a timeout to each request
func TimeoutMiddleware(timeout time.Duration) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(c.Request.Context(), timeout)
		defer cancel()

		c.Request = c.Request.WithContext(ctx)
		c.Next()
	}
}

// MetricsMiddleware adds metrics context to each request
func MetricsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		startTime := time.Now()
		c.Set("start_time", startTime)

		c.Next()

		// Calculate request duration
		duration := time.Since(startTime)
		c.Header("X-Response-Time", duration.String())
	}
}

// Helper functions for working with gin.Context

// GetRequestContext creates a RequestContext from gin.Context
func GetRequestContext(c *gin.Context) *RequestContext {
	return NewRequestContext(c)
}

// GetRequestIDFromGin gets request ID from gin context
func GetRequestIDFromGin(c *gin.Context) string {
	if requestID, exists := c.Get("request_id"); exists {
		if str, ok := requestID.(string); ok {
			return str
		}
	}
	return ""
}

// GetUserIDFromGin gets user ID from gin context
func GetUserIDFromGin(c *gin.Context) string {
	return c.GetString("userID")
}

// GetTraceIDFromGin gets trace ID from gin context
func GetTraceIDFromGin(c *gin.Context) string {
	if traceID, exists := c.Get("trace_id"); exists {
		if str, ok := traceID.(string); ok {
			return str
		}
	}
	return ""
}

// GetUserClaimsFromGin gets user claims from gin context
func GetUserClaimsFromGin(c *gin.Context) interface{} {
	if claims, exists := c.Get("claims"); exists {
		return claims
	}
	return nil
}

// GetStartTimeFromGin gets start time from gin context
func GetStartTimeFromGin(c *gin.Context) time.Time {
	if startTime, exists := c.Get("start_time"); exists {
		if t, ok := startTime.(time.Time); ok {
			return t
		}
	}
	return time.Time{}
}

// GetElapsedTimeFromGin calculates elapsed time from start time in gin context
func GetElapsedTimeFromGin(c *gin.Context) time.Duration {
	startTime := GetStartTimeFromGin(c)
	if startTime.IsZero() {
		return 0
	}
	return time.Since(startTime)
}

// SetContextValue sets a value in gin context
func SetContextValue(c *gin.Context, key string, value interface{}) {
	c.Set(key, value)
}

// GetContextValue gets a value from gin context
func GetContextValue(c *gin.Context, key string) (interface{}, bool) {
	return c.Get(key)
}

// GenerateRequestID generates a new UUID for request ID
func GenerateRequestID() string {
	return uuid.New().String()
}

// GenerateTraceID generates a new UUID for trace ID
func GenerateTraceID() string {
	return uuid.New().String()
}
