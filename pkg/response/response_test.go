package response

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	pkgErrors "template-golang/pkg/errors"
)

func setupGin() (*gin.Engine, *httptest.ResponseRecorder) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	w := httptest.NewRecorder()
	return router, w
}

func TestDefaultPagination(t *testing.T) {
	pagination := DefaultPagination()
	assert.Equal(t, 1, pagination.Page)
	assert.Equal(t, 10, pagination.Limit)
}

func TestPaginationRequest_Offset(t *testing.T) {
	tests := []struct {
		name     string
		page     int
		limit    int
		expected int
	}{
		{"first page", 1, 10, 0},
		{"second page", 2, 10, 10},
		{"third page", 3, 5, 10},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := PaginationRequest{Page: tt.page, Limit: tt.limit}
			assert.Equal(t, tt.expected, p.Offset())
		})
	}
}

func TestPaginationRequest_CalculateTotalPages(t *testing.T) {
	tests := []struct {
		name     string
		limit    int
		total    int
		expected int
	}{
		{"exact division", 10, 100, 10},
		{"with remainder", 10, 95, 10},
		{"less than limit", 10, 5, 1},
		{"zero limit", 0, 100, 0},
		{"zero total", 10, 0, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := PaginationRequest{Limit: tt.limit}
			assert.Equal(t, tt.expected, p.CalculateTotalPages(tt.total))
		})
	}
}

func TestPaginationRequest_ValidateAndDefault(t *testing.T) {
	tests := []struct {
		name          string
		inputPage     int
		inputLimit    int
		expectedPage  int
		expectedLimit int
	}{
		{"negative page", -1, 20, 1, 20},
		{"zero page", 0, 20, 1, 20},
		{"negative limit", 5, -1, 5, 10},
		{"zero limit", 5, 0, 5, 10},
		{"limit too high", 5, 200, 5, 100},
		{"valid values", 3, 25, 3, 25},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := PaginationRequest{Page: tt.inputPage, Limit: tt.inputLimit}
			p.ValidateAndDefault()
			assert.Equal(t, tt.expectedPage, p.Page)
			assert.Equal(t, tt.expectedLimit, p.Limit)
		})
	}
}

func TestJSON(t *testing.T) {
	router, w := setupGin()

	router.GET("/test", func(c *gin.Context) {
		JSON(c, http.StatusOK, map[string]string{"message": "test"})
	})

	req := httptest.NewRequest("GET", "/test", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response Response
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.True(t, response.Success)
	assert.NotNil(t, response.Data)
}

func TestSuccess(t *testing.T) {
	router, w := setupGin()

	router.GET("/test", func(c *gin.Context) {
		Success(c, map[string]string{"key": "value"})
	})

	req := httptest.NewRequest("GET", "/test", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response Response
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.True(t, response.Success)

	data, ok := response.Data.(map[string]interface{})
	assert.True(t, ok)
	assert.Equal(t, "value", data["key"])
}

func TestSuccessWithMessage(t *testing.T) {
	router, w := setupGin()

	router.GET("/test", func(c *gin.Context) {
		SuccessWithMessage(c, "Operation successful", map[string]string{"key": "value"})
	})

	req := httptest.NewRequest("GET", "/test", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response Response
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.True(t, response.Success)
	assert.Equal(t, "Operation successful", response.Message)
}

func TestCreated(t *testing.T) {
	router, w := setupGin()

	router.POST("/test", func(c *gin.Context) {
		Created(c, map[string]string{"id": "123"})
	})

	req := httptest.NewRequest("POST", "/test", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)

	var response Response
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.True(t, response.Success)
}

func TestCreatedWithMessage(t *testing.T) {
	router, w := setupGin()

	router.POST("/test", func(c *gin.Context) {
		CreatedWithMessage(c, "Resource created", map[string]string{"id": "123"})
	})

	req := httptest.NewRequest("POST", "/test", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)

	var response Response
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.True(t, response.Success)
	assert.Equal(t, "Resource created", response.Message)
}

func TestNoContent(t *testing.T) {
	router, w := setupGin()

	router.DELETE("/test", func(c *gin.Context) {
		NoContent(c)
	})

	req := httptest.NewRequest("DELETE", "/test", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNoContent, w.Code)
	assert.Empty(t, w.Body.String())
}

func TestError_WithAppError(t *testing.T) {
	router, w := setupGin()

	router.GET("/test", func(c *gin.Context) {
		err := pkgErrors.Validation("Invalid input").WithDetails("Field is required")
		Error(c, err)
	})

	req := httptest.NewRequest("GET", "/test", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	var response Response
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.False(t, response.Success)
	assert.NotNil(t, response.Error)
	assert.Equal(t, "validation", response.Error.Type)
	assert.Equal(t, "Invalid input", response.Error.Message)
	assert.Equal(t, "Field is required", response.Error.Details)
}

func TestError_WithRegularError(t *testing.T) {
	router, w := setupGin()

	router.GET("/test", func(c *gin.Context) {
		err := errors.New("something went wrong")
		Error(c, err)
	})

	req := httptest.NewRequest("GET", "/test", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)

	var response Response
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.False(t, response.Success)
	assert.NotNil(t, response.Error)
	assert.Equal(t, "internal", response.Error.Type)
	assert.Equal(t, "something went wrong", response.Error.Message)
}

func TestErrorWithCode(t *testing.T) {
	router, w := setupGin()

	router.GET("/test", func(c *gin.Context) {
		err := pkgErrors.Validation("Invalid input")
		ErrorWithCode(c, err, "VALIDATION_001")
	})

	req := httptest.NewRequest("GET", "/test", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	var response Response
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.False(t, response.Success)
	assert.NotNil(t, response.Error)
	assert.Equal(t, "VALIDATION_001", response.Error.Code)
}

func TestConvenienceErrorFunctions(t *testing.T) {
	tests := []struct {
		name       string
		fn         func(*gin.Context)
		statusCode int
	}{
		{
			name:       "BadRequest",
			fn:         func(c *gin.Context) { BadRequest(c, "Bad request") },
			statusCode: http.StatusBadRequest,
		},
		{
			name:       "Unauthorized",
			fn:         func(c *gin.Context) { Unauthorized(c, "Unauthorized") },
			statusCode: http.StatusUnauthorized,
		},
		{
			name:       "Forbidden",
			fn:         func(c *gin.Context) { Forbidden(c, "Forbidden") },
			statusCode: http.StatusForbidden,
		},
		{
			name:       "NotFound",
			fn:         func(c *gin.Context) { NotFound(c, "Not found") },
			statusCode: http.StatusNotFound,
		},
		{
			name:       "Conflict",
			fn:         func(c *gin.Context) { Conflict(c, "Conflict") },
			statusCode: http.StatusConflict,
		},
		{
			name:       "InternalServerError",
			fn:         func(c *gin.Context) { InternalServerError(c, "Internal error") },
			statusCode: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			router, w := setupGin()

			router.GET("/test", tt.fn)

			req := httptest.NewRequest("GET", "/test", nil)
			router.ServeHTTP(w, req)

			assert.Equal(t, tt.statusCode, w.Code)

			var response Response
			err := json.Unmarshal(w.Body.Bytes(), &response)
			assert.NoError(t, err)
			assert.False(t, response.Success)
			assert.NotNil(t, response.Error)
		})
	}
}

func TestValidationError(t *testing.T) {
	router, w := setupGin()

	router.POST("/test", func(c *gin.Context) {
		errors := pkgErrors.NewErrorList()
		errors.AddValidation("name", "is required")
		errors.AddValidation("email", "is invalid")
		ValidationError(c, errors)
	})

	req := httptest.NewRequest("POST", "/test", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	var response Response
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.False(t, response.Success)
	assert.NotNil(t, response.Error)
	assert.Equal(t, "validation", response.Error.Type)
	assert.Equal(t, "Validation failed", response.Error.Message)
}

func TestPaginated(t *testing.T) {
	router, w := setupGin()

	router.GET("/test", func(c *gin.Context) {
		data := []string{"item1", "item2", "item3"}
		pagination := PaginationRequest{Page: 1, Limit: 10}
		total := 50
		Paginated(c, data, pagination, total)
	})

	req := httptest.NewRequest("GET", "/test", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response Response
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.True(t, response.Success)
	assert.NotNil(t, response.Meta)
	assert.Equal(t, 1, response.Meta.Page)
	assert.Equal(t, 10, response.Meta.Limit)
	assert.Equal(t, 50, response.Meta.Total)
	assert.Equal(t, 5, response.Meta.TotalPages)
}

func TestPaginatedWithMessage(t *testing.T) {
	router, w := setupGin()

	router.GET("/test", func(c *gin.Context) {
		data := []string{"item1", "item2"}
		pagination := PaginationRequest{Page: 2, Limit: 5}
		total := 25
		PaginatedWithMessage(c, "Items retrieved", data, pagination, total)
	})

	req := httptest.NewRequest("GET", "/test", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response Response
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.True(t, response.Success)
	assert.Equal(t, "Items retrieved", response.Message)
	assert.NotNil(t, response.Meta)
	assert.Equal(t, 5, response.Meta.TotalPages)
}

func TestGetPaginationFromContext(t *testing.T) {
	router, w := setupGin()

	router.GET("/test", func(c *gin.Context) {
		pagination := GetPaginationFromContext(c)
		JSON(c, http.StatusOK, pagination)
	})

	req := httptest.NewRequest("GET", "/test?page=2&limit=20", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response Response
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)

	data, ok := response.Data.(map[string]interface{})
	assert.True(t, ok)
	assert.Equal(t, float64(2), data["page"])  // JSON numbers are float64
	assert.Equal(t, float64(20), data["limit"])
}

func TestBindAndValidate(t *testing.T) {
	type TestRequest struct {
		Name string `json:"name"`
		Age  int    `json:"age"`
	}

	router, w := setupGin()

	router.POST("/test", func(c *gin.Context) {
		var req TestRequest
		if err := BindAndValidate(c, &req); err != nil {
			Error(c, err)
			return
		}
		Success(c, req)
	})

	body := `{"name":"John","age":30}`
	req := httptest.NewRequest("POST", "/test", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestBindAndValidate_InvalidJSON(t *testing.T) {
	type TestRequest struct {
		Name string `json:"name"`
	}

	router, w := setupGin()

	router.POST("/test", func(c *gin.Context) {
		var req TestRequest
		if err := BindAndValidate(c, &req); err != nil {
			Error(c, err)
			return
		}
		Success(c, req)
	})

	body := `{"name":}`  // Invalid JSON
	req := httptest.NewRequest("POST", "/test", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestBindQueryAndValidate(t *testing.T) {
	type QueryRequest struct {
		Search string `form:"search"`
		Sort   string `form:"sort"`
	}

	router, w := setupGin()

	router.GET("/test", func(c *gin.Context) {
		var req QueryRequest
		if err := BindQueryAndValidate(c, &req); err != nil {
			Error(c, err)
			return
		}
		Success(c, req)
	})

	req := httptest.NewRequest("GET", "/test?search=john&sort=name", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestErrorHandler(t *testing.T) {
	router, w := setupGin()
	router.Use(ErrorHandler())

	router.GET("/test", func(c *gin.Context) {
		c.Error(errors.New("test error"))
	})

	req := httptest.NewRequest("GET", "/test", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)

	var response Response
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.False(t, response.Success)
	assert.NotNil(t, response.Error)
}

func TestCORS(t *testing.T) {
	router, w := setupGin()
	router.Use(CORS())

	router.GET("/test", func(c *gin.Context) {
		Success(c, "test")
	})

	req := httptest.NewRequest("GET", "/test", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, "*", w.Header().Get("Access-Control-Allow-Origin"))
	assert.Contains(t, w.Header().Get("Access-Control-Allow-Methods"), "GET")
	assert.Contains(t, w.Header().Get("Access-Control-Allow-Headers"), "Authorization")
}

func TestCORS_OPTIONS(t *testing.T) {
	router, w := setupGin()
	router.Use(CORS())

	router.OPTIONS("/test", func(c *gin.Context) {
		Success(c, "should not reach here")
	})

	req := httptest.NewRequest("OPTIONS", "/test", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNoContent, w.Code)
	assert.Empty(t, w.Body.String())
}
