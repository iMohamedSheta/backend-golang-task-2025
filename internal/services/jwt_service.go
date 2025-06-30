package services

import (
	"strconv"
	"taskgo/internal/config"
	"taskgo/internal/enums"
	"taskgo/internal/types"
	"taskgo/pkg/errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type JwtService struct {
	config *config.Config
}

func NewJwtService(cfg *config.Config) *JwtService {
	return &JwtService{config: cfg}
}

func (j *JwtService) GenerateAccessToken(userID uint, role enums.UserRole) (string, error) {
	expiry := j.config.GetDuration("auth.jwt.access_token.expiry", 30*time.Minute)
	claims := j.buildClaims(userID, role, enums.AccessToken, expiry)
	return j.signToken(claims)
}

func (j *JwtService) GenerateRefreshToken(userID uint, role enums.UserRole) (string, error) {
	expiry := j.config.GetDuration("auth.jwt.refresh_token.expiry", 168*time.Hour)
	claims := j.buildClaims(userID, role, enums.RefreshToken, expiry)
	return j.signToken(claims)
}

func (j *JwtService) ValidateAuthToken(tokenType enums.JwtTokenType, tokenStr string) (*types.MyClaims, error) {
	expectedIssuer := j.config.GetString("auth.jwt.issuer", "TaskGo")
	expectedAudience := j.config.GetString("auth.jwt.audience", "TaskGoAudience")

	claims := &types.MyClaims{}
	parser := jwt.NewParser(
		jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}),
		jwt.WithIssuer(expectedIssuer),
		jwt.WithAudience(expectedAudience),
		jwt.WithLeeway(5*time.Second),
	)

	token, err := parser.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrSignatureInvalid
		}
		secret, err := j.getSecret()
		if err != nil {
			return "", err
		}
		return secret, nil
	})

	if err != nil || !token.Valid {
		return nil, errors.NewUnAuthorizedError("unauthenticated", "invalid token", err)
	}

	if !j.isValidTokenType(tokenType, claims.TokenType) {
		return nil, errors.NewUnAuthorizedError("unauthenticated", "invalid token type", nil)
	}

	return claims, nil
}

func (j *JwtService) signToken(claims types.MyClaims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	secret, err := j.getSecret()
	if err != nil {
		return "", err
	}
	return token.SignedString(secret)
}

func (j *JwtService) buildClaims(userID uint, role enums.UserRole, tokenType enums.JwtTokenType, expiry time.Duration) types.MyClaims {
	return types.MyClaims{
		Role:      role,
		TokenType: string(tokenType),
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(expiry)),
			Issuer:    j.config.GetString("auth.jwt.issuer", "TaskGo"),
			Audience:  []string{j.config.GetString("auth.jwt.audience", "TaskGoAudience")},
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Subject:   strconv.FormatUint(uint64(userID), 10),
		},
	}
}

func (j *JwtService) isValidTokenType(expected enums.JwtTokenType, actual string) bool {
	return actual == string(expected)
}

func (j *JwtService) getSecret() ([]byte, error) {
	secret := j.config.GetString("app.secret", "")
	if secret == "" {
		return nil, errors.NewServerError("internal server error", "internal server error: Missing or invalid app secret in configuration", nil)
	}
	return []byte(secret), nil
}
