package handlers

import (
	"net/http"
	"strings"
	"taskgo/internal/api/requests"
	"taskgo/internal/api/responses"
	"taskgo/internal/repository"
	"taskgo/internal/services"
	pkgErrors "taskgo/pkg/errors"
	"taskgo/pkg/logger"
	"taskgo/pkg/response"
	"taskgo/pkg/utils"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
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
func (h *AuthHandler) Login(c *gin.Context) {
	var req requests.LoginRequest

	if err := utils.BindToRequestAndExtractFields(c, &req); err != nil {
		logger.Log().Error("Failed to bind login request", zap.Error(err))
		response.BadRequestBindingJson(c, err)
		return
	}

	token, refreshToken, user, err := h.authService.Login(&req)

	if err != nil {
		// Check if it's a validation error
		if valErr, ok := pkgErrors.AsValidationError(err); ok {
			response.ValidationErrorJson(c, valErr)
			return
		}

		// Check if it's a custom server error (with public message)
		if serverErr, ok := pkgErrors.AsServerError(err); ok {
			logger.Log().Error("Failed to login user", zap.Error(err))
			response.ServerErrorJson(c, serverErr)
			return
		}

		// Otherwise treat as unknown error
		logger.Log().Error("Failed to login user", zap.Error(err))
		response.ServerErrorJson(c, nil)
		return
	}

	var loginResponse responses.LoginResponse
	loginResponse.Response(c, user, token, refreshToken)
}

func (h *AuthHandler) RefreshAccessToken(c *gin.Context) {
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" || len(authHeader) < 8 || authHeader[:7] != "Bearer " {
		response.BadRequestJson(c, "Missing or invalid Authorization header")
		return
	}

	refreshToken := strings.TrimPrefix(authHeader, "Bearer ")

	accessToken, err := h.authService.RefreshAccessToken(refreshToken)

	if err != nil {
		if unAuthorizedErr, ok := pkgErrors.AsUnAuthorizedError(err); ok {
			response.UnauthorizedJson(c, unAuthorizedErr)
			return
		}

		// Check if it's a custom server error (with public message)
		if serverErr, ok := pkgErrors.AsServerError(err); ok {
			logger.Log().Error("Failed to refresh access token", zap.Error(err))
			response.ServerErrorJson(c, serverErr)
			return
		}

		// Unhandled error
		logger.Log().Error("Failed to refresh access token With unknown error", zap.Error(err))
		response.ServerErrorJson(c, nil)
		return
	}

	// Return new access token
	response.Json(c, "Access token refreshed successfully", map[string]any{
		"access_token": accessToken,
	}, http.StatusOK)
}
