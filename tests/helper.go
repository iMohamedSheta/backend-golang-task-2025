package tests

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

// Helper function to create test context with response writer
func createTestContext(method, url string, body interface{}) (*httptest.ResponseRecorder, *gin.Context) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()

	var req *http.Request
	if body != nil {
		jsonData, _ := json.Marshal(body)
		req, _ = http.NewRequest(method, url, bytes.NewBuffer(jsonData))
		req.Header.Set("Content-Type", "application/json")
	} else {
		req, _ = http.NewRequest(method, url, nil)
	}

	c, _ := gin.CreateTestContext(w)
	c.Request = req

	return w, c
}

// Helper function to assert validation error response
func assertValidationError(t *testing.T, w *httptest.ResponseRecorder, expectedErrors map[string]string) {
	assert.Equal(t, http.StatusUnprocessableEntity, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)

	assert.Equal(t, "VALIDATION_ERROR", response["error_code"])
	assert.Contains(t, response, "data")

	data := response["data"].(map[string]interface{})
	for field, expectedMessage := range expectedErrors {
		assert.Contains(t, data, field)
		assert.Equal(t, expectedMessage, data[field])
	}
}
