package tests

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"taskgo/internal/api/handlers"
	"taskgo/internal/api/middleware"
	"taskgo/internal/api/requests"
	"taskgo/internal/database/models"
	"taskgo/pkg/database"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestAuthHandler_Login_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	handler := handlers.NewAuthHandler()

	db := database.GetDB()
	testUser := models.User{
		FirstName:   "Test",
		LastName:    "User",
		Email:       "test@example.com",
		Password:    "123456789",
		PhoneNumber: "010123456789",
		Role:        "customer",
		IsActive:    true,
	}
	result := db.Create(&testUser)
	assert.NoError(t, result.Error)

	loginRequest := requests.LoginRequest{
		Email:    testUser.Email,
		Password: "123456789",
	}

	w, c := createTestContext("POST", "/login", loginRequest)

	wrappedLogin := middleware.HandleErrors(handler.Login)
	wrappedLogin(c)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)

	assert.Contains(t, response["data"], "access_token")
	assert.Contains(t, response["data"], "refresh_token")

	truncateTables()
}

func TestAuthHandler_Login_InvalidCredentials(t *testing.T) {
	gin.SetMode(gin.TestMode)
	handler := handlers.NewAuthHandler()

	loginRequest := requests.LoginRequest{
		Email:    "nonexistent@example.com",
		Password: "wrongpassword",
	}

	w, c := createTestContext("POST", "/login", loginRequest)
	wrappedLogin := middleware.HandleErrors(handler.Login)
	wrappedLogin(c)

	assertValidationError(t, w, map[string]string{
		"email": "User not found",
	})
}

func TestAuthHandler_Login_InvalidJSON(t *testing.T) {
	gin.SetMode(gin.TestMode)
	handler := handlers.NewAuthHandler()

	req, _ := http.NewRequest("POST", "/login", bytes.NewBuffer([]byte("invalid json")))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req

	wrappedLogin := middleware.HandleErrors(handler.Login)
	wrappedLogin(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}
