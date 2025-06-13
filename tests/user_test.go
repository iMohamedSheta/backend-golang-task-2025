package tests

import (
	"encoding/json"
	"net/http"
	"os"
	"taskgo/internal/api/handlers"
	"taskgo/internal/api/middleware"
	"taskgo/internal/api/requests"
	"taskgo/internal/database/models"
	"taskgo/pkg/database"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

// TestMain runs before all tests
func TestMain(m *testing.M) {
	TestSetup()
	code := m.Run()
	TestTeardown()
	os.Exit(code)
}

func TestUserHandler_CreateUser_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	handler := handlers.NewUserHandler()

	createUserRequest := requests.CreateUserRequest{
		FirstName:   "John",
		LastName:    "Doe",
		Email:       "john.doe@test.com",
		Password:    "password123",
		PhoneNumber: "01023456789",
	}

	w, c := createTestContext("POST", "/users", createUserRequest)

	wrappedCreateUser := middleware.HandleErrors(handler.CreateUser)
	wrappedCreateUser(c)

	assert.Equal(t, http.StatusCreated, w.Code)

	// user exists in db
	db := database.GetDB()
	var user models.User
	result := db.Where("email = ?", createUserRequest.Email).First(&user)
	assert.NoError(t, result.Error)
	assert.Equal(t, createUserRequest.Email, user.Email)
	assert.Equal(t, createUserRequest.FirstName, user.FirstName)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)

	if data, exists := response["data"]; exists {
		userData := data.(map[string]interface{})
		assert.Equal(t, createUserRequest.Email, userData["email"])
		assert.Equal(t, createUserRequest.FirstName, userData["first_name"])
		assert.Equal(t, createUserRequest.LastName, userData["last_name"])
		assert.Equal(t, createUserRequest.PhoneNumber, userData["phone_number"])
		assert.Equal(t, "customer", userData["role"])
		assert.Equal(t, user.ID, userData["id"])
	}

	truncateTables()
}

func TestUserHandler_CreateUser_DuplicateEmail(t *testing.T) {
	gin.SetMode(gin.TestMode)
	handler := handlers.NewUserHandler()

	db := database.GetDB()
	existingUser := models.User{
		FirstName:   "Existing",
		LastName:    "User",
		Email:       "duplicate@test.com",
		Password:    "password",
		PhoneNumber: "01012345678",
		Role:        "customer",
		IsActive:    true,
	}
	db.Create(&existingUser)

	createUserRequest := requests.CreateUserRequest{
		FirstName:   "John",
		LastName:    "Doe",
		Email:       "duplicate@test.com", // Use existing email
		Password:    "password123",
		PhoneNumber: "01023456789",
	}

	w, c := createTestContext("POST", "/users", createUserRequest)
	wrappedCreateUser := middleware.HandleErrors(handler.CreateUser)
	wrappedCreateUser(c)

	assertValidationError(t, w, map[string]string{
		"email": "Email is already taken",
	})

	truncateTables()
}

func TestUserHandler_CreateUser_InvalidData(t *testing.T) {
	gin.SetMode(gin.TestMode)
	handler := handlers.NewUserHandler()

	createUserRequest := requests.CreateUserRequest{
		FirstName: "John",
	}

	w, c := createTestContext("POST", "/users", createUserRequest)
	wrappedCreateUser := middleware.HandleErrors(handler.CreateUser)
	wrappedCreateUser(c)

	assertValidationError(t, w, map[string]string{
		"email":        "Email is required",
		"last_name":    "Last Name is required",
		"password":     "Password is required",
		"phone_number": "Phone Number is required",
	})
}
