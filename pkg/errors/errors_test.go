package errors

import (
	"errors"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAppError_Error(t *testing.T) {
	// Test without cause
	err := &AppError{
		Type:    ErrorTypeValidation,
		Message: "validation failed",
	}
	assert.Equal(t, "validation: validation failed", err.Error())

	// Test with cause
	cause := errors.New("original error")
	err = &AppError{
		Type:    ErrorTypeValidation,
		Message: "validation failed",
		Cause:   cause,
	}
	assert.Equal(t, "validation: validation failed (caused by: original error)", err.Error())
}

func TestAppError_Unwrap(t *testing.T) {
	cause := errors.New("original error")
	err := &AppError{
		Type:    ErrorTypeValidation,
		Message: "validation failed",
		Cause:   cause,
	}

	assert.Equal(t, cause, err.Unwrap())
}

func TestAppError_Is(t *testing.T) {
	err1 := &AppError{Type: ErrorTypeValidation, Message: "test"}
	err2 := &AppError{Type: ErrorTypeValidation, Message: "test2"}
	err3 := &AppError{Type: ErrorTypeNotFound, Message: "test"}

	assert.True(t, err1.Is(err2))
	assert.False(t, err1.Is(err3))
	assert.False(t, err1.Is(nil))
}

func TestAppError_As(t *testing.T) {
	err := &AppError{Type: ErrorTypeValidation, Message: "test"}

	var target *AppError
	assert.True(t, err.As(&target))
	assert.Equal(t, err, target)

	var invalidTarget *string
	assert.False(t, err.As(&invalidTarget))
}

func TestAppError_WithContext(t *testing.T) {
	err := &AppError{
		Type:    ErrorTypeValidation,
		Message: "test",
		Context: make(map[string]interface{}),
	}

	err.WithContext("field", "username")
	assert.Equal(t, "username", err.Context["field"])
}

func TestAppError_WithDetails(t *testing.T) {
	err := &AppError{
		Type:    ErrorTypeValidation,
		Message: "test",
	}

	err.WithDetails("detailed explanation")
	assert.Equal(t, "detailed explanation", err.Details)
}

func TestAppError_WithStack(t *testing.T) {
	err := &AppError{
		Type:    ErrorTypeValidation,
		Message: "test",
	}

	err.WithStack()
	assert.NotEmpty(t, err.Stack)
	assert.Contains(t, err.Stack, "TestAppError_WithStack")
}

func TestNew(t *testing.T) {
	err := New(ErrorTypeValidation, "test message")

	assert.Equal(t, ErrorTypeValidation, err.Type)
	assert.Equal(t, "test message", err.Message)
	assert.Equal(t, http.StatusBadRequest, err.StatusCode)
	assert.NotNil(t, err.Context)
}

func TestWrap(t *testing.T) {
	// Test wrapping nil error
	result := Wrap(nil, ErrorTypeValidation, "test")
	assert.Nil(t, result)

	// Test wrapping regular error
	originalErr := errors.New("original error")
	wrapped := Wrap(originalErr, ErrorTypeValidation, "validation failed")

	assert.Equal(t, ErrorTypeValidation, wrapped.Type)
	assert.Equal(t, "validation failed", wrapped.Message)
	assert.Equal(t, originalErr, wrapped.Cause)
	assert.Equal(t, http.StatusBadRequest, wrapped.StatusCode)

	// Test wrapping AppError preserves type when not specified
	appErr := &AppError{Type: ErrorTypeNotFound, Message: "not found"}
	wrapped = Wrap(appErr, "", "wrapped message")
	assert.Equal(t, ErrorTypeNotFound, wrapped.Type)
}

func TestWrapWithStack(t *testing.T) {
	originalErr := errors.New("original error")
	wrapped := WrapWithStack(originalErr, ErrorTypeInternal, "internal error")

	assert.Equal(t, ErrorTypeInternal, wrapped.Type)
	assert.Equal(t, "internal error", wrapped.Message)
	assert.Equal(t, originalErr, wrapped.Cause)
	assert.NotEmpty(t, wrapped.Stack)
}

func TestConvenienceFunctions(t *testing.T) {
	tests := []struct {
		name     string
		fn       func() *AppError
		expected ErrorType
		status   int
	}{
		{
			name:     "Validation",
			fn:       func() *AppError { return Validation("validation error") },
			expected: ErrorTypeValidation,
			status:   http.StatusBadRequest,
		},
		{
			name:     "NotFound",
			fn:       func() *AppError { return NotFound("not found") },
			expected: ErrorTypeNotFound,
			status:   http.StatusNotFound,
		},
		{
			name:     "Unauthorized",
			fn:       func() *AppError { return Unauthorized("unauthorized") },
			expected: ErrorTypeUnauthorized,
			status:   http.StatusUnauthorized,
		},
		{
			name:     "Forbidden",
			fn:       func() *AppError { return Forbidden("forbidden") },
			expected: ErrorTypeForbidden,
			status:   http.StatusForbidden,
		},
		{
			name:     "Conflict",
			fn:       func() *AppError { return Conflict("conflict") },
			expected: ErrorTypeConflict,
			status:   http.StatusConflict,
		},
		{
			name:     "BadRequest",
			fn:       func() *AppError { return BadRequest("bad request") },
			expected: ErrorTypeBadRequest,
			status:   http.StatusBadRequest,
		},
		{
			name:     "Timeout",
			fn:       func() *AppError { return Timeout("timeout") },
			expected: ErrorTypeTimeout,
			status:   http.StatusRequestTimeout,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.fn()
			assert.Equal(t, tt.expected, err.Type)
			assert.Equal(t, tt.status, err.StatusCode)
		})
	}
}

func TestInternal(t *testing.T) {
	err := Internal("internal error")
	assert.Equal(t, ErrorTypeInternal, err.Type)
	assert.Equal(t, "internal error", err.Message)
	assert.Equal(t, http.StatusInternalServerError, err.StatusCode)
	assert.NotEmpty(t, err.Stack) // Should have stack trace
}

func TestInternalWithCause(t *testing.T) {
	cause := errors.New("database connection failed")
	err := InternalWithCause("internal error", cause)

	assert.Equal(t, ErrorTypeInternal, err.Type)
	assert.Equal(t, "internal error", err.Message)
	assert.Equal(t, cause, err.Cause)
	assert.NotEmpty(t, err.Stack)
}

func TestValidationWithDetails(t *testing.T) {
	err := ValidationWithDetails("validation failed", "field is required")
	assert.Equal(t, ErrorTypeValidation, err.Type)
	assert.Equal(t, "validation failed", err.Message)
	assert.Equal(t, "field is required", err.Details)
}

func TestDatabase(t *testing.T) {
	cause := errors.New("connection timeout")
	err := Database("database error", cause)

	assert.Equal(t, ErrorTypeDatabase, err.Type)
	assert.Equal(t, "database error", err.Message)
	assert.Equal(t, cause, err.Cause)
	assert.Equal(t, http.StatusInternalServerError, err.StatusCode)
}

func TestExternal(t *testing.T) {
	cause := errors.New("API timeout")
	err := External("external service error", cause)

	assert.Equal(t, ErrorTypeExternal, err.Type)
	assert.Equal(t, "external service error", err.Message)
	assert.Equal(t, cause, err.Cause)
	assert.Equal(t, http.StatusBadGateway, err.StatusCode)
}

func TestIsType(t *testing.T) {
	err := &AppError{Type: ErrorTypeValidation, Message: "test"}

	assert.True(t, IsType(err, ErrorTypeValidation))
	assert.False(t, IsType(err, ErrorTypeNotFound))

	// Test with non-AppError
	regularErr := errors.New("regular error")
	assert.False(t, IsType(regularErr, ErrorTypeValidation))
}

func TestGetStatusCode(t *testing.T) {
	err := &AppError{Type: ErrorTypeValidation, StatusCode: http.StatusBadRequest}
	assert.Equal(t, http.StatusBadRequest, GetStatusCode(err))

	// Test with non-AppError
	regularErr := errors.New("regular error")
	assert.Equal(t, http.StatusInternalServerError, GetStatusCode(regularErr))
}

func TestTypeToStatusCode(t *testing.T) {
	tests := []struct {
		errorType  ErrorType
		statusCode int
	}{
		{ErrorTypeValidation, http.StatusBadRequest},
		{ErrorTypeNotFound, http.StatusNotFound},
		{ErrorTypeUnauthorized, http.StatusUnauthorized},
		{ErrorTypeForbidden, http.StatusForbidden},
		{ErrorTypeConflict, http.StatusConflict},
		{ErrorTypeBadRequest, http.StatusBadRequest},
		{ErrorTypeTimeout, http.StatusRequestTimeout},
		{ErrorTypeDatabase, http.StatusInternalServerError},
		{ErrorTypeExternal, http.StatusBadGateway},
		{ErrorTypeInternal, http.StatusInternalServerError},
		{ErrorType("unknown"), http.StatusInternalServerError},
	}

	for _, tt := range tests {
		t.Run(string(tt.errorType), func(t *testing.T) {
			assert.Equal(t, tt.statusCode, typeToStatusCode(tt.errorType))
		})
	}
}

func TestErrorList(t *testing.T) {
	el := NewErrorList()
	assert.Empty(t, el.Errors)
	assert.False(t, el.HasErrors())
	assert.Equal(t, "no errors", el.Error())
	assert.Nil(t, el.First())
}

func TestErrorList_Add(t *testing.T) {
	el := NewErrorList()

	// Add nil error (should be ignored)
	el.Add(nil)
	assert.False(t, el.HasErrors())

	// Add valid error
	err := Validation("test error")
	el.Add(err)
	assert.True(t, el.HasErrors())
	assert.Len(t, el.Errors, 1)
	assert.Equal(t, err, el.First())
}

func TestErrorList_AddValidation(t *testing.T) {
	el := NewErrorList()
	el.AddValidation("username", "is required")

	assert.True(t, el.HasErrors())
	assert.Len(t, el.Errors, 1)

	err := el.First()
	assert.Equal(t, ErrorTypeValidation, err.Type)
	assert.Equal(t, "is required", err.Message)
	assert.Equal(t, "username", err.Context["field"])
}

func TestErrorList_Error(t *testing.T) {
	el := NewErrorList()

	// Single error
	el.Add(Validation("test error"))
	assert.Equal(t, "validation: test error", el.Error())

	// Multiple errors
	el.Add(NotFound("not found"))
	assert.Contains(t, el.Error(), "multiple errors:")
	assert.Contains(t, el.Error(), "validation: test error")
	assert.Contains(t, el.Error(), "not_found: not found")
}
