package errors

import (
	"errors"
	"fmt"
	"net/http"
	"runtime"
	"strings"
)

// ErrorType represents the type of error
type ErrorType string

const (
	// ErrorTypeValidation represents validation errors
	ErrorTypeValidation ErrorType = "validation"
	// ErrorTypeNotFound represents not found errors
	ErrorTypeNotFound ErrorType = "not_found"
	// ErrorTypeUnauthorized represents unauthorized errors
	ErrorTypeUnauthorized ErrorType = "unauthorized"
	// ErrorTypeForbidden represents forbidden errors
	ErrorTypeForbidden ErrorType = "forbidden"
	// ErrorTypeConflict represents conflict errors
	ErrorTypeConflict ErrorType = "conflict"
	// ErrorTypeInternal represents internal server errors
	ErrorTypeInternal ErrorType = "internal"
	// ErrorTypeBadRequest represents bad request errors
	ErrorTypeBadRequest ErrorType = "bad_request"
	// ErrorTypeTimeout represents timeout errors
	ErrorTypeTimeout ErrorType = "timeout"
	// ErrorTypeDatabase represents database errors
	ErrorTypeDatabase ErrorType = "database"
	// ErrorTypeExternal represents external service errors
	ErrorTypeExternal ErrorType = "external"
)

// AppError represents an application error with additional context
type AppError struct {
	Type       ErrorType              `json:"type"`
	Message    string                 `json:"message"`
	Details    string                 `json:"details,omitempty"`
	StatusCode int                    `json:"status_code"`
	Cause      error                  `json:"-"`
	Context    map[string]interface{} `json:"context,omitempty"`
	Stack      string                 `json:"stack,omitempty"`
}

// Error implements the error interface
func (e *AppError) Error() string {
	if e.Cause != nil {
		return fmt.Sprintf("%s: %s (caused by: %s)", e.Type, e.Message, e.Cause.Error())
	}
	return fmt.Sprintf("%s: %s", e.Type, e.Message)
}

// Unwrap returns the underlying cause
func (e *AppError) Unwrap() error {
	return e.Cause
}

// Is checks if the error is of a specific type
func (e *AppError) Is(target error) bool {
	if target == nil {
		return false
	}

	if appErr, ok := target.(*AppError); ok {
		return e.Type == appErr.Type
	}

	return errors.Is(e.Cause, target)
}

// As finds the first error in the chain that matches target
func (e *AppError) As(target interface{}) bool {
	if target == nil {
		return false
	}

	if appErr, ok := target.(**AppError); ok {
		*appErr = e
		return true
	}

	return errors.As(e.Cause, target)
}

// WithContext adds context to the error
func (e *AppError) WithContext(key string, value interface{}) *AppError {
	if e.Context == nil {
		e.Context = make(map[string]interface{})
	}
	e.Context[key] = value
	return e
}

// WithDetails adds details to the error
func (e *AppError) WithDetails(details string) *AppError {
	e.Details = details
	return e
}

// WithStack adds stack trace to the error
func (e *AppError) WithStack() *AppError {
	e.Stack = captureStack()
	return e
}

// New creates a new AppError
func New(errorType ErrorType, message string) *AppError {
	return &AppError{
		Type:       errorType,
		Message:    message,
		StatusCode: typeToStatusCode(errorType),
		Context:    make(map[string]interface{}),
	}
}

// Wrap wraps an existing error with additional context
func Wrap(err error, errorType ErrorType, message string) *AppError {
	if err == nil {
		return nil
	}

	// If it's already an AppError, preserve the original type if not specified
	if appErr, ok := err.(*AppError); ok && errorType == "" {
		errorType = appErr.Type
	}

	return &AppError{
		Type:       errorType,
		Message:    message,
		StatusCode: typeToStatusCode(errorType),
		Cause:      err,
		Context:    make(map[string]interface{}),
	}
}

// WrapWithStack wraps an error and captures the stack trace
func WrapWithStack(err error, errorType ErrorType, message string) *AppError {
	appErr := Wrap(err, errorType, message)
	if appErr != nil {
		appErr.WithStack()
	}
	return appErr
}

// Validation creates a validation error
func Validation(message string) *AppError {
	return New(ErrorTypeValidation, message)
}

// ValidationWithDetails creates a validation error with details
func ValidationWithDetails(message, details string) *AppError {
	return New(ErrorTypeValidation, message).WithDetails(details)
}

// NotFound creates a not found error
func NotFound(message string) *AppError {
	return New(ErrorTypeNotFound, message)
}

// Unauthorized creates an unauthorized error
func Unauthorized(message string) *AppError {
	return New(ErrorTypeUnauthorized, message)
}

// Forbidden creates a forbidden error
func Forbidden(message string) *AppError {
	return New(ErrorTypeForbidden, message)
}

// Conflict creates a conflict error
func Conflict(message string) *AppError {
	return New(ErrorTypeConflict, message)
}

// Internal creates an internal server error
func Internal(message string) *AppError {
	return New(ErrorTypeInternal, message).WithStack()
}

// InternalWithCause creates an internal server error with a cause
func InternalWithCause(message string, cause error) *AppError {
	return Wrap(cause, ErrorTypeInternal, message).WithStack()
}

// BadRequest creates a bad request error
func BadRequest(message string) *AppError {
	return New(ErrorTypeBadRequest, message)
}

// Timeout creates a timeout error
func Timeout(message string) *AppError {
	return New(ErrorTypeTimeout, message)
}

// Database creates a database error
func Database(message string, cause error) *AppError {
	return Wrap(cause, ErrorTypeDatabase, message)
}

// External creates an external service error
func External(message string, cause error) *AppError {
	return Wrap(cause, ErrorTypeExternal, message)
}

// IsType checks if an error is of a specific type
func IsType(err error, errorType ErrorType) bool {
	if appErr, ok := err.(*AppError); ok {
		return appErr.Type == errorType
	}
	return false
}

// GetStatusCode returns the HTTP status code for an error
func GetStatusCode(err error) int {
	if appErr, ok := err.(*AppError); ok {
		return appErr.StatusCode
	}
	return http.StatusInternalServerError
}

// typeToStatusCode maps error types to HTTP status codes
func typeToStatusCode(errorType ErrorType) int {
	switch errorType {
	case ErrorTypeValidation:
		return http.StatusBadRequest
	case ErrorTypeNotFound:
		return http.StatusNotFound
	case ErrorTypeUnauthorized:
		return http.StatusUnauthorized
	case ErrorTypeForbidden:
		return http.StatusForbidden
	case ErrorTypeConflict:
		return http.StatusConflict
	case ErrorTypeBadRequest:
		return http.StatusBadRequest
	case ErrorTypeTimeout:
		return http.StatusRequestTimeout
	case ErrorTypeDatabase:
		return http.StatusInternalServerError
	case ErrorTypeExternal:
		return http.StatusBadGateway
	case ErrorTypeInternal:
		return http.StatusInternalServerError
	default:
		return http.StatusInternalServerError
	}
}

// captureStack captures the current stack trace
func captureStack() string {
	const depth = 32
	var pcs [depth]uintptr
	n := runtime.Callers(3, pcs[:]) // Skip runtime.Callers, captureStack, and the calling function

	frames := runtime.CallersFrames(pcs[:n])
	var stack strings.Builder

	for {
		frame, more := frames.Next()
		fmt.Fprintf(&stack, "%s\n\t%s:%d\n", frame.Function, frame.File, frame.Line)
		if !more {
			break
		}
	}

	return stack.String()
}

// ErrorList represents a collection of errors
type ErrorList struct {
	Errors []*AppError `json:"errors"`
}

// NewErrorList creates a new error list
func NewErrorList() *ErrorList {
	return &ErrorList{
		Errors: make([]*AppError, 0),
	}
}

// Add adds an error to the list
func (el *ErrorList) Add(err *AppError) {
	if err != nil {
		el.Errors = append(el.Errors, err)
	}
}

// AddValidation adds a validation error to the list
func (el *ErrorList) AddValidation(field, message string) {
	err := Validation(message).WithContext("field", field)
	el.Add(err)
}

// HasErrors returns true if there are errors in the list
func (el *ErrorList) HasErrors() bool {
	return len(el.Errors) > 0
}

// Error implements the error interface
func (el *ErrorList) Error() string {
	if len(el.Errors) == 0 {
		return "no errors"
	}

	if len(el.Errors) == 1 {
		return el.Errors[0].Error()
	}

	var messages []string
	for _, err := range el.Errors {
		messages = append(messages, err.Error())
	}

	return fmt.Sprintf("multiple errors: %s", strings.Join(messages, "; "))
}

// First returns the first error in the list
func (el *ErrorList) First() *AppError {
	if len(el.Errors) > 0 {
		return el.Errors[0]
	}
	return nil
}
