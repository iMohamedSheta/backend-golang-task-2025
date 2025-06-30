package responses

import (
	"net/http"
	"taskgo/internal/database/models"
	"taskgo/pkg/response"
	"time"

	"github.com/gin-gonic/gin"
)

type CreateUserResponse struct {
	Message string `json:"message" example:"User created successfully"`
	Data    struct {
		User UserData `json:"user"`
	} `json:"data"`
}

type UserData struct {
	ID          uint      `json:"id"`
	FirstName   string    `json:"first_name"`
	LastName    string    `json:"last_name"`
	Email       string    `json:"email"`
	PhoneNumber string    `json:"phone_number"`
	IsActive    bool      `json:"is_active"`
	LastLoginAt time.Time `json:"last_login_at"`
	CreatedAt   string    `json:"created_at"`
}

func SendCreateUserResponse(gin *gin.Context, user *models.User) {
	r := &CreateUserResponse{}
	r.Message = "User created successfully"
	r.Data.User.ID = user.ID
	r.Data.User.FirstName = user.FirstName
	r.Data.User.LastName = user.LastName
	r.Data.User.Email = user.Email
	r.Data.User.PhoneNumber = user.PhoneNumber
	r.Data.User.IsActive = user.IsActive
	r.Data.User.LastLoginAt = user.LastLoginAt
	r.Data.User.CreatedAt = user.CreatedAt.Format("2006-01-02 15:04:05")
	response.Json(gin, r.Message, r.Data, http.StatusCreated)
}

type GetUserResponse struct {
	Message string `json:"message" example:"User retrieved successfully"`
	Data    struct {
		User UserData `json:"user"`
	} `json:"data"`
}

func SendGetUserResponse(gin *gin.Context, user *models.User) {
	r := &GetUserResponse{}
	r.Message = "User retrieved successfully"
	r.Data.User.ID = user.ID
	r.Data.User.FirstName = user.FirstName
	r.Data.User.LastName = user.LastName
	r.Data.User.Email = user.Email
	r.Data.User.PhoneNumber = user.PhoneNumber
	r.Data.User.IsActive = user.IsActive
	r.Data.User.LastLoginAt = user.LastLoginAt
	r.Data.User.CreatedAt = user.CreatedAt.Format("2006-01-02 15:04:05")
	response.Json(gin, r.Message, r.Data, http.StatusOK)
}

type UpdateUserResponse struct {
	Message string `json:"message" example:"User updated successfully"`
	Data    struct {
		User UserData `json:"user"`
	} `json:"data"`
}

func SendUpdateUserResponse(gin *gin.Context, user *models.User) {
	r := &UpdateUserResponse{}
	r.Message = "User updated successfully"
	r.Data.User.ID = user.ID
	r.Data.User.FirstName = user.FirstName
	r.Data.User.LastName = user.LastName
	r.Data.User.Email = user.Email
	r.Data.User.PhoneNumber = user.PhoneNumber
	r.Data.User.IsActive = user.IsActive
	r.Data.User.LastLoginAt = user.LastLoginAt
	r.Data.User.CreatedAt = user.CreatedAt.Format("2006-01-02 15:04:05")
	response.Json(gin, r.Message, r.Data, http.StatusOK)
}
