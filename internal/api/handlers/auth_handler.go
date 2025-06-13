package handlers

import (
	"net/http"
	"strings"
	"taskgo/internal/api/requests"
	"taskgo/internal/api/responses"
	"taskgo/internal/repository"
	"taskgo/internal/services"
	pkgErrors "taskgo/pkg/errors"
	"taskgo/pkg/response"
	"taskgo/pkg/utils"
	"taskgo/pkg/validate"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	authService *services.AuthService
}

func NewAuthHandler() *AuthHandler {
	userRepo := repository.NewUserRepository()
	return &AuthHandler{
		authService: services.NewAuthService(userRepo),
	}
}

// AuthHandler examples
// User login
//
//		@Router			/login [post]
//
//		@Summary		User login
//		@Description	Authenticate user with email and password
//		@Tags			Authentication
//		@Accept			json
//		@Produce		json
//
//	 @Param			request	body	requests.LoginRequest	true	"Login request"
//
//		@Success		200		{object}	responses.LoginResponse		"Login successful"
//		@Failure		400		{object}	response.BadRequestResponse	"Bad request"
//		@Failure		422		{object}	response.ValidationErrorResponse "Validation error"
//		@Failure		500		{object}	response.ServerErrorResponse "Internal server error"
func (h *AuthHandler) Login(c *gin.Context) error {
	var req requests.LoginRequest

	if err := utils.BindToRequestAndExtractFields(c, &req); err != nil {
		return pkgErrors.NewBadRequestError("Bad request", "Failed to bind login request", err)
	}

	// Validate request
	valid, errors := validate.ValidateRequest(&req)
	if !valid {
		return pkgErrors.NewValidationError(errors)
	}

	// Login user
	token, refreshToken, user, err := h.authService.Login(&req)

	if err != nil {
		return err
	}

	var loginResponse responses.LoginResponse
	loginResponse.Response(c, user, token, refreshToken)
	return nil
}

func (h *AuthHandler) RefreshAccessToken(c *gin.Context) error {
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" || len(authHeader) < 8 || authHeader[:7] != "Bearer " {
		return pkgErrors.NewUnAuthorizedError("Unauthorized", "No authorization header provided", nil)
	}

	refreshToken := strings.TrimPrefix(authHeader, "Bearer ")

	accessToken, err := h.authService.RefreshAccessToken(refreshToken)

	if err != nil {
		if unAuthorizedErr, ok := pkgErrors.AsUnAuthorizedError(err); ok {
			return unAuthorizedErr
		}

		if serverErr, ok := pkgErrors.AsServerError(err); ok {
			return serverErr
		}

		return pkgErrors.NewServerError("Internal server error", "Failed to refresh access token", err)
	}

	// Return new access token
	response.Json(c, "Access token refreshed successfully", map[string]any{
		"access_token": accessToken,
	}, http.StatusOK)

	return nil
}
