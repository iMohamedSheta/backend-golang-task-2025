package requests

type CreateUserRequest struct {
	Email       string `json:"email" validate:"required,email,unique_db=users-email"`
	Password    string `json:"password" validate:"required,min=8,max=30"`
	FirstName   string `json:"first_name" validate:"required,min=2,max=30"`
	LastName    string `json:"last_name" validate:"required,min=2,max=30"`
	PhoneNumber string `json:"phone_number" validate:"required,egyptian_phone"` // Egypt phone regex
	Request
}

func (r *CreateUserRequest) Messages() map[string]string {
	return map[string]string{
		"email.required":              "Email is required",
		"email.email":                 "Email is not valid",
		"email.unique_db":             "Email is already taken",
		"password.required":           "Password is required",
		"password.min":                "Password must be at least 8 characters",
		"password.max":                "Password must be at most 30 characters",
		"first_name.required":         "First Name is required",
		"first_name.min":              "First Name must be at least 2 characters",
		"first_name.max":              "First Name must be at most 30 characters",
		"last_name.required":          "Last Name is required",
		"last_name.min":               "Last Name must be at least 2 characters",
		"last_name.max":               "Last Name must be at most 30 characters",
		"phone_number.required":       "Phone Number is required",
		"phone_number.egyptian_phone": "Phone Number must be in Egyptian format",
	}
}

type UpdateUserRequest struct {
	Email       string `json:"email" validate:"required,email,unique_db=users-email"`
	Password    string `json:"password" validate:"required,min=8,max=30"`
	FirstName   string `json:"first_name" validate:"required,min=2,max=30"`
	LastName    string `json:"last_name" validate:"required,min=2,max=30"`
	PhoneNumber string `json:"phone_number" validate:"required,egyptian_phone"` // Egypt phone regex
	Request
}

func (r *UpdateUserRequest) Messages() map[string]string {
	return map[string]string{
		"email.required":              "Email is required",
		"email.email":                 "Email is not valid",
		"email.unique_db":             "Email is already taken",
		"password.required":           "Password is required",
		"password.min":                "Password must be at least 8 characters",
		"password.max":                "Password must be at most 30 characters",
		"first_name.required":         "First Name is required",
		"first_name.min":              "First Name must be at least 2 characters",
		"first_name.max":              "First Name must be at most 30 characters",
		"last_name.required":          "Last Name is required",
		"last_name.min":               "Last Name must be at least 2 characters",
		"last_name.max":               "Last Name must be at most 30 characters",
		"phone_number.required":       "Phone Number is required",
		"phone_number.egyptian_phone": "Phone Number must be in Egyptian format",
	}
}
