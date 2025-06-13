package handlers

import (
	"taskgo/internal/api/requests"
	"taskgo/internal/helpers"
	"taskgo/internal/policies"
	"taskgo/internal/repository"
	"taskgo/internal/services"
	"taskgo/pkg/errors"
	"taskgo/pkg/response"
	"taskgo/pkg/utils"
	"taskgo/pkg/validate"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	userService *services.UserService
	userPolicy  *policies.UserPolicy
}

func NewUserHandler() *UserHandler {
	userRepo := repository.NewUserRepository()
	return &UserHandler{
		userService: services.NewUserService(userRepo),
		userPolicy:  &policies.UserPolicy{},
	}
}

// Create a new customer (register)
func (h *UserHandler) CreateUser(c *gin.Context) error {
	var req requests.CreateUserRequest

	// Policy check
	if !h.userPolicy.CanCreate(nil) {
		return errors.NewUnAuthorizedError("Unauthorized", "You are not allowed to view this user", nil)
	}

	if err := utils.BindToRequestAndExtractFields(c, &req); err != nil {
		return errors.NewBadRequestError("", "BadRequestError: Failed to bind request body to request object", err)
	}

	valid, validErrorsList := validate.ValidateRequest(&req)
	if !valid {
		return errors.NewValidationError(validErrorsList)
	}

	_, err := h.userService.CreateUser(&req)

	if err != nil {
		// Check if it's a custom server error (with public message)
		if serverErr, ok := errors.AsServerError(err); ok {
			return serverErr
		}

		return errors.NewServerError("Failed to create the user", "Failed to create the user", err)
	}

	response.Json(c, "User created successfully", nil, 201)

	return nil
}

// Get user details by ID
func (h *UserHandler) GetUser(c *gin.Context) error {
	userId := c.Param("id")
	// Get auth user
	authUser, authorizeErr := helpers.GetAuthUser(c)
	if authorizeErr != nil {
		return authorizeErr
	}

	// Policy check
	if !h.userPolicy.CanView(authUser, userId) {
		return errors.NewUnAuthorizedError("Unauthorized", "You are not allowed to view this user", nil)
	}

	// Get target user
	targetUser, err := h.userService.GetUserById(userId)
	if err != nil {
		if notFoundErr, ok := errors.AsNotFoundError(err); ok {
			return notFoundErr
		}

		return errors.NewServerError("", "Internal Server Error: Failed to get the user by id using UserService", err)
	}

	// Example response
	response.Json(c, "User details retrieved successfully", map[string]any{
		"user": map[string]any{
			"id":         targetUser.ID,
			"first_name": targetUser.FirstName,
			"last_name":  targetUser.LastName,
			"email":      targetUser.Email,
			"created_at": targetUser.CreatedAt,
			"updated_at": targetUser.UpdatedAt,
		},
	}, 200)

	return nil
}

// Update user details by ID
func (h *UserHandler) UpdateUser(c *gin.Context) error {
	var req requests.UpdateUserRequest

	if err := utils.BindToRequestAndExtractFields(c, &req); err != nil {
		return errors.NewBadRequestBindingError("", "Failed to bind request body to UpdateUserRequest", err)
	}

	userId := c.Param("id")
	authUser, authorizeErr := helpers.GetAuthUser(c)
	if authorizeErr != nil {
		return authorizeErr
	}

	// Policy check
	if !h.userPolicy.CanUpdate(authUser, userId) {
		return errors.NewUnAuthorizedError("Unauthorized", "You are not allowed to view this user", nil)
	}

	// Validate request
	valid, validErrorsList := validate.ValidateRequest(&req)
	if !valid {
		return errors.NewValidationError(validErrorsList)
	}

	// Update user details
	targetUser, err := h.userService.UpdateUserAndGet(userId, &req)
	if err != nil {
		if notfoundErr, ok := errors.AsNotFoundError(err); ok {
			return notfoundErr
		}
		return errors.NewServerError("", "Failed to update user details", err)
	}

	response.Json(c, "User details updated successfully", map[string]any{
		"user": map[string]any{
			"id":         targetUser.ID,
			"first_name": targetUser.FirstName,
			"last_name":  targetUser.LastName,
			"email":      targetUser.Email,
			"phone":      targetUser.PhoneNumber,
			"created_at": targetUser.CreatedAt,
			"updated_at": targetUser.UpdatedAt,
		},
	}, 200)

	return nil
}
