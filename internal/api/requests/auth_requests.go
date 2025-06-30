package requests

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email,exists_db=users-email" example:"test@test.com"`
	Password string `json:"password" validate:"required,min=8,max=30" example:"123456789"`
	Request
}

func (r *LoginRequest) Messages() map[string]string {
	return map[string]string{
		"email.required":    "Email is required",
		"email.email":       "Invalid email format",
		"email.exists_db":   "User does not exist",
		"password.required": "Password is required",
		"password.min":      "Password must be at least 8 characters",
		"password.max":      "Password must be at most 30 characters",
	}
}
