package handlers

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"template-golang/modules/cockroach/models"
	"template-golang/modules/cockroach/usecases/mocks"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestDetectCockroach(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name           string
		requestBody    interface{}
		mockError      error
		expectedStatus int
		expectedBody   map[string]interface{}
		skipSetupMock  bool
	}{
		{
			name: "Success",
			requestBody: models.AddCockroachData{
				Amount: 3,
			},
			mockError:      nil,
			expectedStatus: http.StatusOK,
			expectedBody: map[string]interface{}{
				"message": "Success 🪳🪳🪳",
			},
		},
		{
			name: "Invalid request body - zero amount",
			requestBody: models.AddCockroachData{
				Amount: 0,
			},
			mockError:      nil,
			expectedStatus: http.StatusBadRequest,
			expectedBody: map[string]interface{}{
				"message": "Key: 'AddCockroachData.Amount' Error:Field validation for 'Amount' failed on the 'required' tag",
			},
			skipSetupMock: true,
		},
		{
			name:           "Invalid request body",
			requestBody:    map[string]interface{}{},
			mockError:      nil,
			expectedStatus: http.StatusBadRequest,
			expectedBody: map[string]interface{}{
				"message": "Key: 'AddCockroachData.Amount' Error:Field validation for 'Amount' failed on the 'required' tag",
			},
			skipSetupMock: true,
		},
		{
			name: "Invalid request body - negative amount",
			requestBody: map[string]interface{}{
				"amount": -1,
			},
			mockError:      nil,
			expectedStatus: http.StatusBadRequest,
			expectedBody: map[string]interface{}{
				"message": "json: cannot unmarshal number -1 into Go struct field AddCockroachData.amount of type uint32",
			},
			skipSetupMock: true,
		},
		{
			name: "Invalid request body - string amount",
			requestBody: map[string]interface{}{
				"amount": "three",
			},
			mockError:      nil,
			expectedStatus: http.StatusBadRequest,
			expectedBody: map[string]interface{}{
				"message": "json: cannot unmarshal string into Go struct field AddCockroachData.amount of type uint32",
			},
			skipSetupMock: true,
		},
		{
			name:           "Empty request body",
			requestBody:    nil,
			mockError:      nil,
			expectedStatus: http.StatusBadRequest,
			expectedBody: map[string]interface{}{
				"message": "Key: 'AddCockroachData.Amount' Error:Field validation for 'Amount' failed on the 'required' tag",
			},
			skipSetupMock: true,
		},
		{
			name: "Processing error",
			requestBody: models.AddCockroachData{
				Amount: 2,
			},
			mockError:      errors.New("processing error"),
			expectedStatus: http.StatusInternalServerError,
			expectedBody: map[string]interface{}{
				"message": "Processing data failed",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockUsecase := mocks.NewMockCockroachUsecase(t)
			// Convert request body to JSON
			jsonBody, _ := json.Marshal(tt.requestBody)

			// Create request
			req := httptest.NewRequest(http.MethodPost, "/detect-cockroach", bytes.NewBuffer(jsonBody))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			// Setup router
			r := gin.New()
			handler := NewCockroachHttpHandler(mockUsecase)
			r.POST("/detect-cockroach", handler.DetectCockroach)

			if !tt.skipSetupMock {
				// Set mock expectations - use testify mock.Anything for any argument
				mockUsecase.On("ProcessData", mock.Anything).Return(tt.mockError)
			}

			// Perform request
			r.ServeHTTP(w, req)

			// Assert response
			assert.Equal(t, tt.expectedStatus, w.Code)

			var responseBody map[string]interface{}
			_ = json.Unmarshal(w.Body.Bytes(), &responseBody)
			assert.Equal(t, tt.expectedBody, responseBody)
		})
	}
}
