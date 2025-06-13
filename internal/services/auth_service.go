package services

import (
	"strconv"
	"taskgo/internal/api/requests"
	"taskgo/internal/database/models"
	"taskgo/internal/enums"
	"taskgo/internal/helpers"
	"taskgo/internal/repository"
	pkgErrors "taskgo/pkg/errors"
	"taskgo/pkg/logger"
	"taskgo/pkg/validate"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

type AuthService struct {
	userRepository *repository.UserRepository
}

func NewAuthService(userRepo *repository.UserRepository) *AuthService {
	return &AuthService{userRepository: userRepo}
}

func (s *AuthService) Login(req *requests.LoginRequest) (string, string, *models.User, error) {
	// Validate request
	valid, errors := validate.ValidateRequest(req)
	if !valid {
		return "", "", nil, pkgErrors.NewValidationError(errors)
	}

	user, err := s.userRepository.FindByEmail(req.Email)

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return "", "", nil, pkgErrors.NewValidationError(map[string]any{"email": "User not found"})
		}
		logger.Log().Error("Error finding user", zap.Error(err))
		return "", "", nil, pkgErrors.NewServerError("", "Failed to find user by email", err)
	}

	if !user.CheckPassword(req.Password) {
		return "", "", nil, pkgErrors.NewValidationError(map[string]any{"password": "Incorrect password"})
	}

	// Update last login timestamp
	if err := s.userRepository.UpdateLastLogin(user); err != nil {
		logger.Log().Error("Failed to update last login", zap.Error(err))
	}

	token, err := helpers.GenerateAccessToken(user.ID, user.Role)
	if err != nil {
		return "", "", nil, pkgErrors.NewServerError("", "Failed to generate access tokens", err)
	}

	refreshToken, err := helpers.GenerateRefreshToken(user.ID, user.Role)
	if err != nil {
		return "", "", nil, pkgErrors.NewServerError("", "Failed to generate refresh tokens", err)
	}

	return token, refreshToken, user, nil
}

// Refresh access token if refresh token is valid
func (s *AuthService) RefreshAccessToken(refreshToken string) (string, error) {

	claims, err := helpers.ValidateAuthToken(enums.RefreshToken, refreshToken)
	if err != nil {
		return "", pkgErrors.NewUnAuthorizedError("", "Invalid refresh token", err)
	}

	subject, err := claims.GetSubject()
	if err != nil {
		return "", pkgErrors.NewUnAuthorizedError("", "Invalid subject in token", err)
	}

	userID, err := strconv.ParseUint(subject, 10, 64)
	if err != nil {
		return "", pkgErrors.NewUnAuthorizedError("", "Invalid subject in token", err)
	}

	userRole, err := claims.GetRole()

	if err != nil {
		return "", pkgErrors.NewUnAuthorizedError("", "Invalid role in token", err)
	}

	newAccessToken, err := helpers.GenerateAccessToken(uint(userID), userRole)

	if err != nil {
		return "", pkgErrors.NewServerError("", "Failed to generate new access token", err)
	}

	return newAccessToken, nil
}
