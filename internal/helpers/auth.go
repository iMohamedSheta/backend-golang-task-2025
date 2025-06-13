package helpers

import (
	"fmt"
	"strconv"
	"taskgo/internal/config"
	"taskgo/internal/database/models"
	"taskgo/internal/enums"
	"taskgo/internal/repository"
	"taskgo/internal/types"
	"taskgo/pkg/errors"
	"taskgo/pkg/logger"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"go.uber.org/zap"
)

// GenerateAccessToken generates a JWT access token for the given user ID and role.
func GenerateAccessToken(userID uint, role enums.UserRole) (string, error) {
	hmacSecret, err := GetAppSecret()
	if err != nil {
		logger.Log().Error("Failed to get app secret", zap.Error(err))
		return "", fmt.Errorf("internal server error")
	}

	accessTokenExpiration := config.App.GetDuration("auth.jwt.access_token.expiry", 30*time.Minute)

	claims := types.MyClaims{
		Role:      role,
		TokenType: string(enums.AccessToken),
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(accessTokenExpiration)),
			Issuer:    config.App.GetString("auth.jwt.issuer", "TaskGo"),
			Audience:  []string{config.App.GetString("auth.jwt.audience", "TaskGoAudience")},
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Subject:   strconv.FormatUint(uint64(userID), 10),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(hmacSecret)
}

// GenerateRefreshToken generates a JWT refresh token for the given user ID and role.
func GenerateRefreshToken(userID uint, role enums.UserRole) (string, error) {
	hmacSecret, err := GetAppSecret()
	if err != nil {
		logger.Log().Error("Failed to get app secret", zap.Error(err))
		return "", fmt.Errorf("internal server error")
	}

	refreshTokenExpiration := config.App.GetDuration("auth.jwt.refresh_token.expiry", 168*time.Hour)

	claims := types.MyClaims{
		Role:      role,
		TokenType: string(enums.RefreshToken),
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(refreshTokenExpiration)),
			Issuer:    config.App.GetString("auth.jwt.issuer", "TaskGo"),
			Audience:  []string{config.App.GetString("auth.jwt.audience", "TaskGoAudience")},
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Subject:   strconv.FormatUint(uint64(userID), 10),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(hmacSecret)
}

// ValidateAuthToken validates a JWT Auth token (access token, refresh token) and returns the claims if valid.
func ValidateAuthToken(tokenType enums.JwtTokenType, tokenStr string) (*types.MyClaims, error) {
	hmacSecret, err := GetAppSecret()
	if err != nil {
		logger.Log().Error("Failed to get app secret", zap.Error(err))
		return nil, fmt.Errorf("internal server error")
	}

	expectedIssuer := config.App.GetString("auth.jwt.issuer", "TaskGo")
	expectedAudience := config.App.GetString("auth.jwt.audience", "TaskGoAudience")

	claims := &types.MyClaims{}

	// Create parser with validation options
	parser := jwt.NewParser(
		jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}),
		jwt.WithIssuer(expectedIssuer),     // Validate issuer
		jwt.WithAudience(expectedAudience), // Validate audience
		jwt.WithLeeway(5*time.Second),      // Allow 5 seconds leeway for time drift
	)

	token, err := parser.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrSignatureInvalid
		}
		return hmacSecret, nil
	})

	if err != nil {
		logger.Log().Warn("JWT parsing failed", zap.Error(err))
		return nil, fmt.Errorf("invalid or expired token")
	}

	if !token.Valid {
		return nil, fmt.Errorf("invalid or expired token")
	}

	// validate token type
	if !IsValidTokenType(tokenType, claims.TokenType) {
		return nil, fmt.Errorf("invalid token type")
	}

	return claims, nil
}

// Checks if the JWT token is expired based on the "exp"
func IsTokenExpired(claims *types.MyClaims) bool {
	if claims.ExpiresAt == nil {
		return true
	}
	return time.Now().After(claims.ExpiresAt.Time)
}

// Checks if the token type in claims matches the expected token type
func IsValidTokenType(tokenType enums.JwtTokenType, claimType interface{}) bool {
	claimTypeStr, ok := claimType.(string)
	if !ok {
		return false
	}
	return claimTypeStr == string(tokenType)
}

// Checks if the user role is valid
func IsValidUserRole(role string) bool {
	switch enums.UserRole(role) {
	case enums.RoleAdmin, enums.RoleCustomer:
		return true
	default:
		return false
	}
}

func GetAuthId(ctx *gin.Context) (string, *errors.UnAuthorizedError) {
	val, exists := ctx.Get(string(enums.ContextKeyAuthId))
	if !exists {
		return "", errors.NewUnAuthorizedError("unauthenticated", "Authentication ID not found in context", nil)
	}

	authId, ok := val.(string)
	if !ok || authId == "" {
		return "", errors.NewUnAuthorizedError("unauthenticated", "Invalid authentication ID in context", nil)
	}

	return authId, nil
}

func GetAuthUser(ctx *gin.Context) (*models.User, *errors.UnAuthorizedError) {
	val, unAuthorizedErr := GetAuthId(ctx)

	if unAuthorizedErr != nil {
		return nil, unAuthorizedErr
	}

	userRepo := repository.NewUserRepository()
	user, err := userRepo.FindById(val)

	if err != nil {
		return nil, errors.NewUnAuthorizedError("unauthenticated", "User not found", err)
	}

	return user, nil
}
