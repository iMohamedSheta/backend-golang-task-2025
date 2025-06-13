package types

import (
	"taskgo/internal/enums"

	"github.com/golang-jwt/jwt/v5"
)

type MyClaims struct {
	Role      enums.UserRole `json:"role"`
	TokenType string         `json:"token_type"`
	jwt.RegisteredClaims
}

func (c *MyClaims) GetRole() (enums.UserRole, error) {
	return c.Role, nil
}
