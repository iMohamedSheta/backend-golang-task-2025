package services

import (
	"context"
	"strconv"
	"taskgo/internal/api/requests"
	"taskgo/internal/database/models"
	"taskgo/internal/deps"
	"taskgo/internal/enums"
	"taskgo/internal/repository"
	pkgErrors "taskgo/pkg/errors"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

type AuthService struct {
	userRepository *repository.UserRepository
	jwtService     *JwtService
}

func NewAuthService(userRepo *repository.UserRepository, jwtService *JwtService) *AuthService {
	return &AuthService{
		userRepository: userRepo,
		jwtService:     jwtService,
	}
}

// Login a user into the system using email and password
func (s *AuthService) Login(ctx context.Context, req *requests.LoginRequest) (string, string, *models.User, error) {
	log := deps.Log()
	user, err := s.userRepository.FindByEmail(req.Email)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return "", "", nil, pkgErrors.NewValidationError(map[string]any{"email": "User not found"})
		}
		log.Log().Error("Error finding user", zap.Error(err))
		return "", "", nil, pkgErrors.NewServerError("", "Failed to find user by email", err)
	}

	if !user.CheckPassword(req.Password) {
		return "", "", nil, pkgErrors.NewValidationError(map[string]any{"password": "Incorrect password"})
	}

	// Update last login timestamp
	if err := s.userRepository.UpdateLastLogin(user); err != nil {
		log.Log().Error("Failed to update last login", zap.Error(err))
	}

	token, err := s.jwtService.GenerateAccessToken(user.ID, user.Role)
	if err != nil {
		return "", "", nil, pkgErrors.NewServerError("", "Failed to generate access tokens", err)
	}

	refreshToken, err := s.jwtService.GenerateRefreshToken(user.ID, user.Role)
	if err != nil {
		return "", "", nil, pkgErrors.NewServerError("", "Failed to generate refresh tokens", err)
	}

	return token, refreshToken, user, nil
}

// Refresh access token if refresh token is valid
func (s *AuthService) RefreshAccessToken(ctx context.Context, refreshToken string) (string, error) {

	claims, err := s.jwtService.ValidateAuthToken(enums.RefreshToken, refreshToken)
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

	newAccessToken, err := s.jwtService.GenerateAccessToken(uint(userID), userRole)

	if err != nil {
		return "", pkgErrors.NewServerError("", "Failed to generate new access token", err)
	}

	return newAccessToken, nil
}
