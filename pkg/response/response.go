package response

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	pkgErrors "template-golang/pkg/errors"
)

// Response represents a standardized API response
type Response struct {
	Success   bool        `json:"success"`
	Message   string      `json:"message,omitempty"`
	Data      interface{} `json:"data,omitempty"`
	Error     *ErrorInfo  `json:"error,omitempty"`
	Meta      *Meta       `json:"meta,omitempty"`
	Timestamp time.Time   `json:"timestamp"`
}

// ErrorInfo represents error information in the response
type ErrorInfo struct {
	Type    string                 `json:"type"`
	Message string                 `json:"message"`
	Details string                 `json:"details,omitempty"`
	Context map[string]interface{} `json:"context,omitempty"`
	Code    string                 `json:"code,omitempty"`
}

// Meta represents metadata for responses (pagination, etc.)
type Meta struct {
	Page       int `json:"page,omitempty"`
	Limit      int `json:"limit,omitempty"`
	Total      int `json:"total,omitempty"`
	TotalPages int `json:"total_pages,omitempty"`
}

// PaginationRequest represents pagination parameters
type PaginationRequest struct {
	Page  int `json:"page" form:"page" validate:"min=1"`
	Limit int `json:"limit" form:"limit" validate:"min=1,max=100"`
}

// DefaultPagination returns default pagination values
func DefaultPagination() PaginationRequest {
	return PaginationRequest{
		Page:  1,
		Limit: 10,
	}
}

// Offset calculates the offset for database queries
func (p PaginationRequest) Offset() int {
	return (p.Page - 1) * p.Limit
}

// CalculateTotalPages calculates total pages based on total items
func (p PaginationRequest) CalculateTotalPages(total int) int {
	if p.Limit <= 0 {
		return 0
	}
	return (total + p.Limit - 1) / p.Limit
}

// ValidateAndDefault validates pagination request and sets defaults
func (p *PaginationRequest) ValidateAndDefault() {
	if p.Page <= 0 {
		p.Page = 1
	}
	if p.Limit <= 0 {
		p.Limit = 10
	}
	if p.Limit > 100 {
		p.Limit = 100
	}
}

// JSON sends a JSON response
func JSON(c *gin.Context, statusCode int, data interface{}) {
	response := Response{
		Success:   statusCode < 400,
		Data:      data,
		Timestamp: time.Now(),
	}

	c.JSON(statusCode, response)
}

// Success sends a successful response
func Success(c *gin.Context, data interface{}) {
	JSON(c, http.StatusOK, data)
}

// SuccessWithMessage sends a successful response with a message
func SuccessWithMessage(c *gin.Context, message string, data interface{}) {
	response := Response{
		Success:   true,
		Message:   message,
		Data:      data,
		Timestamp: time.Now(),
	}

	c.JSON(http.StatusOK, response)
}

// Created sends a 201 Created response
func Created(c *gin.Context, data interface{}) {
	JSON(c, http.StatusCreated, data)
}

// CreatedWithMessage sends a 201 Created response with a message
func CreatedWithMessage(c *gin.Context, message string, data interface{}) {
	response := Response{
		Success:   true,
		Message:   message,
		Data:      data,
		Timestamp: time.Now(),
	}

	c.JSON(http.StatusCreated, response)
}

// NoContent sends a 204 No Content response
func NoContent(c *gin.Context) {
	c.Status(http.StatusNoContent)
}

// Error sends an error response
func Error(c *gin.Context, err error) {
	var statusCode int
	var errorInfo *ErrorInfo

	// Handle AppError
	if appErr, ok := err.(*pkgErrors.AppError); ok {
		statusCode = appErr.StatusCode
		errorInfo = &ErrorInfo{
			Type:    string(appErr.Type),
			Message: appErr.Message,
			Details: appErr.Details,
			Context: appErr.Context,
		}
	} else {
		// Handle regular errors
		statusCode = http.StatusInternalServerError
		errorInfo = &ErrorInfo{
			Type:    "internal",
			Message: err.Error(),
		}
	}

	response := Response{
		Success:   false,
		Error:     errorInfo,
		Timestamp: time.Now(),
	}

	c.JSON(statusCode, response)
}

// ErrorWithCode sends an error response with a custom error code
func ErrorWithCode(c *gin.Context, err error, code string) {
	var statusCode int
	var errorInfo *ErrorInfo

	// Handle AppError
	if appErr, ok := err.(*pkgErrors.AppError); ok {
		statusCode = appErr.StatusCode
		errorInfo = &ErrorInfo{
			Type:    string(appErr.Type),
			Message: appErr.Message,
			Details: appErr.Details,
			Context: appErr.Context,
			Code:    code,
		}
	} else {
		// Handle regular errors
		statusCode = http.StatusInternalServerError
		errorInfo = &ErrorInfo{
			Type:    "internal",
			Message: err.Error(),
			Code:    code,
		}
	}

	response := Response{
		Success:   false,
		Error:     errorInfo,
		Timestamp: time.Now(),
	}

	c.JSON(statusCode, response)
}

// BadRequest sends a 400 Bad Request response
func BadRequest(c *gin.Context, message string) {
	err := pkgErrors.BadRequest(message)
	Error(c, err)
}

// Unauthorized sends a 401 Unauthorized response
func Unauthorized(c *gin.Context, message string) {
	err := pkgErrors.Unauthorized(message)
	Error(c, err)
}

// Forbidden sends a 403 Forbidden response
func Forbidden(c *gin.Context, message string) {
	err := pkgErrors.Forbidden(message)
	Error(c, err)
}

// NotFound sends a 404 Not Found response
func NotFound(c *gin.Context, message string) {
	err := pkgErrors.NotFound(message)
	Error(c, err)
}

// Conflict sends a 409 Conflict response
func Conflict(c *gin.Context, message string) {
	err := pkgErrors.Conflict(message)
	Error(c, err)
}

// InternalServerError sends a 500 Internal Server Error response
func InternalServerError(c *gin.Context, message string) {
	err := pkgErrors.Internal(message)
	Error(c, err)
}

// ValidationError sends a 400 Bad Request response for validation errors
func ValidationError(c *gin.Context, errors *pkgErrors.ErrorList) {
	var errorInfos []ErrorInfo
	for _, err := range errors.Errors {
		errorInfos = append(errorInfos, ErrorInfo{
			Type:    string(err.Type),
			Message: err.Message,
			Details: err.Details,
			Context: err.Context,
		})
	}

	response := Response{
		Success: false,
		Error: &ErrorInfo{
			Type:    "validation",
			Message: "Validation failed",
			Context: map[string]interface{}{
				"errors": errorInfos,
			},
		},
		Timestamp: time.Now(),
	}

	c.JSON(http.StatusBadRequest, response)
}

// Paginated sends a paginated response
func Paginated(c *gin.Context, data interface{}, pagination PaginationRequest, total int) {
	meta := &Meta{
		Page:       pagination.Page,
		Limit:      pagination.Limit,
		Total:      total,
		TotalPages: pagination.CalculateTotalPages(total),
	}

	response := Response{
		Success:   true,
		Data:      data,
		Meta:      meta,
		Timestamp: time.Now(),
	}

	c.JSON(http.StatusOK, response)
}

// PaginatedWithMessage sends a paginated response with a message
func PaginatedWithMessage(c *gin.Context, message string, data interface{}, pagination PaginationRequest, total int) {
	meta := &Meta{
		Page:       pagination.Page,
		Limit:      pagination.Limit,
		Total:      total,
		TotalPages: pagination.CalculateTotalPages(total),
	}

	response := Response{
		Success:   true,
		Message:   message,
		Data:      data,
		Meta:      meta,
		Timestamp: time.Now(),
	}

	c.JSON(http.StatusOK, response)
}

// GetPaginationFromContext extracts pagination parameters from Gin context
func GetPaginationFromContext(c *gin.Context) PaginationRequest {
	pagination := DefaultPagination()

	if err := c.ShouldBindQuery(&pagination); err == nil {
		pagination.ValidateAndDefault()
	}

	return pagination
}

// BindAndValidate binds request data and validates it
func BindAndValidate(c *gin.Context, obj interface{}) error {
	if err := c.ShouldBindJSON(obj); err != nil {
		return pkgErrors.BadRequest("Invalid request format: " + err.Error())
	}

	// Add validation logic here if needed
	return nil
}

// BindQueryAndValidate binds query parameters and validates them
func BindQueryAndValidate(c *gin.Context, obj interface{}) error {
	if err := c.ShouldBindQuery(obj); err != nil {
		return pkgErrors.BadRequest("Invalid query parameters: " + err.Error())
	}

	// Add validation logic here if needed
	return nil
}

// Middleware for error handling
func ErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		// Check if there are any errors
		if len(c.Errors) > 0 {
			// Get the last error
			err := c.Errors.Last().Err

			// If response hasn't been written yet, send error response
			if !c.Writer.Written() {
				Error(c, err)
			}
		}
	}
}

// CORS middleware
func CORS() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Origin, Content-Type, Accept, Authorization, X-Requested-With")
		c.Header("Access-Control-Expose-Headers", "Content-Length")
		c.Header("Access-Control-Allow-Credentials", "true")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	}
}
