package handlers

import (
	"strings"
	"taskgo/internal/api/requests"
	"taskgo/internal/api/responses"
	"taskgo/internal/deps"
	"taskgo/internal/services"
	pkgErrors "taskgo/pkg/errors"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	Handler
	authService *services.AuthService
}

// NewAuthHandler return a new AuthHandler
func NewAuthHandler(authService *services.AuthService) *AuthHandler {
	return &AuthHandler{
		authService: authService,
	}
}

// @Summary     User login
// @Description Authenticate user with email and password
// @Tags        Authentication
// @Accept      json
// @Produce     json
//
// @Param       request  body      requests.LoginRequest            true  "Login request"
//
// @Success     200      {object}  responses.LoginResponse          "Login successful"
// @Failure     400      {object}  response.BadRequestResponse      "Bad Request"
// @Failure     422      {object}  response.ValidationErrorResponse "Validation error"
// @Failure     500      {object}  response.ServerErrorResponse     "Internal server error"
//
// @Router      /login [post]
func (h *AuthHandler) Login(gin *gin.Context) error {
	var req requests.LoginRequest
	var err error

	if err = h.BindBodyAndExtractToRequest(gin, &req); err != nil {
		return pkgErrors.NewBadRequestBindingError("", "Failed to bind and extract request to LoginRequest", err)
	}

	if err = deps.Validator().ValidateRequest(&req); err != nil {
		return err
	}

	ctx := gin.Request.Context()
	token, refreshToken, user, err := h.authService.Login(ctx, &req)
	if err != nil {
		return err
	}

	responses.SendLoginResponse(gin, user, token, refreshToken)
	return nil
}

// @Summary     Refresh access token
// @Description Refresh access token using refresh token
// @Tags        Authentication
// @Accept      json
// @Produce     json
//
// @Param       request  body      requests.LoginRequest                 true  "Login request"
//
// @Success     200      {object}  responses.RefreshAccessTokenResponse  "Access token refreshed successfully"
// @Failure     400      {object}  response.BadRequestResponse           "Bad Request"
// @Failure     401      {object}  response.UnauthorizedResponse         "Unauthorized Action"
// @Failure     500      {object}  response.ServerErrorResponse          "Internal Server Error"
//
// @Router      /refresh-token [post]
func (h *AuthHandler) RefreshAccessToken(gin *gin.Context) error {
	authHeader := gin.GetHeader("Authorization")
	if authHeader == "" || len(authHeader) < 8 || authHeader[:7] != "Bearer " {
		return pkgErrors.NewUnAuthorizedError("", "No authorization header provided", nil)
	}

	refreshToken := strings.TrimPrefix(authHeader, "Bearer ")

	accessToken, err := h.authService.RefreshAccessToken(gin.Request.Context(), refreshToken)
	if err != nil {
		return err
	}

	responses.SendRefreshAccessTokenResponse(gin, accessToken)
	return nil
}
