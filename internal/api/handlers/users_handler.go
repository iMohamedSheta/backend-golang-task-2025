package handlers

import (
	"taskgo/internal/api/requests"
	"taskgo/internal/helpers"
	"taskgo/internal/policies"
	"taskgo/internal/repository"
	"taskgo/internal/services"
	"taskgo/pkg/errors"
	"taskgo/pkg/logger"
	"taskgo/pkg/response"
	"taskgo/pkg/utils"
	"taskgo/pkg/validate"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
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
func (h *UserHandler) CreateUser(c *gin.Context) {
	var req requests.CreateUserRequest

	// Policy check
	if !h.userPolicy.CanCreate(nil) {
		response.UnauthorizedJson(c, errors.NewUnAuthorizedError("Unauthorized", "You are not allowed to view this user", nil))
		return
	}

	if err := utils.BindToRequestAndExtractFields(c, &req); err != nil {
		logger.Log().Error("Failed to bind create user request", zap.Error(err))
		response.BadRequestBindingJson(c, err)
		return
	}

	valid, validErrorsList := validate.ValidateRequest(&req)
	if !valid {
		validErrors := errors.NewValidationError(validErrorsList)
		response.ValidationErrorJson(c, validErrors)
		return
	}

	_, err := h.userService.CreateUser(&req)

	if err != nil {
		// Check if it's a custom server error (with public message)
		if serverErr, ok := errors.AsServerError(err); ok {
			logger.Log().Error("Failed to create user", zap.Error(err))
			response.ServerErrorJson(c, serverErr)
			return
		}

		// Otherwise treat as unknown error
		logger.Log().Error("Failed to create user", zap.Error(err))
		response.ServerErrorJson(c, errors.NewServerError("Failed to create the user", "Failed to create the user", err))
		return
	}

	response.Json(c, "User created successfully", nil, 201)
}

// Get user details by ID
func (h *UserHandler) GetUser(c *gin.Context) {
	userId := c.Param("id")
	// Get auth user
	authUser, authorizeErr := helpers.GetAuthUser(c)
	if authorizeErr != nil {
		response.UnauthorizedJson(c, authorizeErr)
		return
	}

	// Policy check
	if !h.userPolicy.CanView(authUser, userId) {
		response.UnauthorizedJson(c, errors.NewUnAuthorizedError("Unauthorized", "You are not allowed to view this user", nil))
		return
	}

	// Get target user
	targetUser, err := h.userService.GetUserById(userId)
	if err != nil {
		response.ServerErrorJson(c, errors.NewServerError("", "Failed to get the user", err))
		return
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
}

// Update user details by ID
func (h *UserHandler) UpdateUser(c *gin.Context) {
	var req requests.UpdateUserRequest

	if err := utils.BindToRequestAndExtractFields(c, &req); err != nil {
		logger.Log().Error("Failed to bind create user request", zap.Error(err))
		response.BadRequestBindingJson(c, err)
		return
	}

	userId := c.Param("id")
	authUser, authorizeErr := helpers.GetAuthUser(c)
	if authorizeErr != nil {
		response.UnauthorizedJson(c, authorizeErr)
		return
	}

	// Policy check
	if !h.userPolicy.CanUpdate(authUser, userId) {
		response.UnauthorizedJson(c, errors.NewUnAuthorizedError("Unauthorized", "You are not allowed to view this user", nil))
		return
	}

	// Validate request
	valid, validErrorsList := validate.ValidateRequest(&req)
	if !valid {
		response.ValidationErrorJson(c, errors.NewValidationError(validErrorsList))
		return
	}

	// Update user details
	targetUser, err := h.userService.UpdateUserAndGet(userId, &req)
	if err != nil {
		if notfoundErr, ok := errors.AsNotFoundError(err); ok {
			response.NotFoundJson(c, errors.NewNotFoundError("Product not found", "Product not found", notfoundErr))
			return
		}
		logger.Log().Error("Failed to update user details", zap.Error(err))
		response.ServerErrorJson(c, errors.NewServerError("", "Failed to update user details", err))
		return
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
}
