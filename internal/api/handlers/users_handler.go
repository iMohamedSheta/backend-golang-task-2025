package handlers

import (
	"taskgo/internal/api/requests"
	"taskgo/internal/api/responses"
	"taskgo/internal/database/models"
	"taskgo/internal/deps"
	"taskgo/internal/helpers"
	"taskgo/internal/policies"
	"taskgo/internal/services"
	"taskgo/pkg/errors"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	Handler
	userService *services.UserService
	userPolicy  *policies.UserPolicy
}

// NewUserHandler return a new UserHandler
func NewUserHandler(userService *services.UserService, userPolicy *policies.UserPolicy) *UserHandler {
	return &UserHandler{
		userService: userService,
		userPolicy:  userPolicy,
	}
}

// @Summary     Create a new user
// @Description Creates a new user with the provided information.
// @Tags        Users
// @Accept      json
// @Produce     json
//
// @Param       request  body      requests.CreateUserRequest     true  "Create user request body"
//
// @Success     201      {object}  responses.CreateUserResponse    "User created successfully"
// @Failure     400      {object}  response.BadRequestResponse     "Bad Request"
// @Failure     401      {object}  response.UnauthorizedResponse   "Unauthorized Action"
// @Failure     422      {object}  response.ValidationErrorResponse "Validation Error"
// @Failure     500      {object}  response.ServerErrorResponse    "Internal Server Error"
//
// @Router      /users [post]
func (h *UserHandler) CreateUser(gin *gin.Context) error {
	var req requests.CreateUserRequest
	var err error

	if err = h.BindBodyAndExtractToRequest(gin, &req); err != nil {
		return errors.NewBadRequestError("", "BadRequestError: Failed to bind request body to request object", err)
	}

	if !h.userPolicy.CanCreate(nil) {
		return errors.NewUnAuthorizedError("Unauthorized", "You are not allowed to view this user", nil)
	}

	if err = deps.Validator().ValidateRequest(&req); err != nil {
		return err
	}

	var user *models.User
	user, err = h.userService.CreateUser(gin, &req)
	if err != nil {
		return err
	}

	responses.SendCreateUserResponse(gin, user)
	return nil
}

// @Summary     Get user by ID
// @Description Retrieves user details by their ID.
// @Tags        Users
// @Accept      json
// @Produce     json
//
// @Param       id       path      string                         true  "User ID"
//
// @Success     200      {object}  responses.GetUserResponse      "User retrieved successfully"
// @Failure     400      {object}  response.BadRequestResponse     "Bad Request"
// @Failure     401      {object}  response.UnauthorizedResponse   "Unauthorized Action"
// @Failure     404      {object}  response.NotFoundResponse       "User Not Found"
// @Failure     500      {object}  response.ServerErrorResponse    "Internal Server Error"
//
// @Router      /users/{id} [get]
func (h *UserHandler) GetUser(gin *gin.Context) error {
	userId := gin.Param("id")
	authUser, authorizeErr := helpers.GetAuthUser(gin)
	if authorizeErr != nil {
		return authorizeErr
	}

	if !h.userPolicy.CanView(authUser, userId) {
		return errors.NewUnAuthorizedError("Unauthorized", "You are not allowed to view this user", nil)
	}

	targetUser, err := h.userService.GetUserById(gin.Request.Context(), userId)
	if err != nil {
		return err
	}

	responses.SendGetUserResponse(gin, targetUser)
	return nil
}

// @Summary     Update user by ID
// @Description Updates an existing user with the provided ID and request body.
// @Tags        Users
// @Accept      json
// @Produce     json
//
// @Param       id       path      string                         true  "User ID"
// @Param       request  body      requests.UpdateUserRequest     true  "Update user request body"
//
// @Success     200      {object}  responses.UpdateUserResponse   "User updated successfully"
// @Failure     400      {object}  response.BadRequestResponse     "Bad Request"
// @Failure     401      {object}  response.UnauthorizedResponse   "Unauthorized Action"
// @Failure     404      {object}  response.NotFoundResponse       "User Not Found"
// @Failure     422      {object}  response.ValidationErrorResponse "Validation Error"
// @Failure     500      {object}  response.ServerErrorResponse    "Internal Server Error"
//
// @Router      /users/{id} [put]
func (h *UserHandler) UpdateUser(gin *gin.Context) error {
	var req requests.UpdateUserRequest
	var err error

	if err = h.BindBodyAndExtractToRequest(gin, &req); err != nil {
		return errors.NewBadRequestBindingError("", "Failed to bind request body to UpdateUserRequest", err)
	}

	userId := gin.Param("id")
	authUser, authorizeErr := helpers.GetAuthUser(gin)
	if authorizeErr != nil {
		return authorizeErr
	}

	if !h.userPolicy.CanUpdate(authUser, userId) {
		return errors.NewUnAuthorizedError("Unauthorized", "You are not allowed to view this user", nil)
	}

	if err = deps.Validator().ValidateRequest(&req); err != nil {
		return err
	}

	targetUser, err := h.userService.UpdateUserAndGet(gin.Request.Context(), userId, &req)
	if err != nil {
		return err
	}

	responses.SendUpdateUserResponse(gin, targetUser)
	return nil
}
