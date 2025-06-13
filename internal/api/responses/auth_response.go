package responses

import (
	"taskgo/internal/database/models"
	"taskgo/internal/enums"
	"taskgo/pkg/response"
	"time"

	"github.com/gin-gonic/gin"
)

type LoginResponse struct {
	Message string `json:"message" example:"User logged in successfully"`
	Data    struct {
		AccessToken  string `json:"access_token" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."`
		RefreshToken string `json:"refresh_token" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."`
		User         struct {
			ID          uint           `json:"id" example:"1"`
			Email       string         `json:"email" example:"user@example.com"`
			FirstName   string         `json:"first_name" example:"John"`
			LastName    string         `json:"last_name" example:"Doe"`
			Role        enums.UserRole `json:"role" example:"customer"`
			PhoneNumber string         `json:"phone_number" example:"+1234567890"`
			LastLoginAt time.Time      `json:"last_login_at" example:"2024-01-01T00:00:00Z"`
			IsActive    bool           `json:"is_active" example:"true"`
			IsAdmin     bool           `json:"is_admin" example:"false"`
			CreatedAt   time.Time      `json:"created_at" example:"2024-01-01T00:00:00Z"`
			UpdatedAt   time.Time      `json:"updated_at" example:"2024-01-01T00:00:00Z"`
		} `json:"user"`
	} `json:"data"`
}

// return login successful response
func (r *LoginResponse) Response(c *gin.Context, user *models.User, accessToken string, refreshToken string) {
	// message
	r.Message = "User logged in successfully"
	// tokens
	r.Data.AccessToken = accessToken
	r.Data.RefreshToken = refreshToken
	// user data
	r.Data.User.ID = user.ID
	r.Data.User.Email = user.Email
	r.Data.User.FirstName = user.FirstName
	r.Data.User.LastName = user.LastName
	r.Data.User.Role = user.Role
	r.Data.User.PhoneNumber = user.PhoneNumber
	r.Data.User.LastLoginAt = user.LastLoginAt
	r.Data.User.IsActive = user.IsActive
	r.Data.User.IsAdmin = user.Role == enums.RoleAdmin
	r.Data.User.CreatedAt = user.CreatedAt
	r.Data.User.UpdatedAt = user.UpdatedAt

	response.Json(c, r.Message, r.Data, 200)
}
