package helpers

import (
	"taskgo/internal/database/models"
	"taskgo/internal/enums"
	"taskgo/internal/repository"
	"taskgo/pkg/errors"
	"taskgo/pkg/ioc"

	"github.com/gin-gonic/gin"
)

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

func GetAuthUser(ctx *gin.Context) (*models.User, error) {
	val, unAuthorizedErr := GetAuthId(ctx)

	if unAuthorizedErr != nil {
		return nil, unAuthorizedErr
	}

	userRepo, err := ioc.AppMake[*repository.UserRepository]()
	if err != nil {
		return nil, errors.NewServerError("internal server error", "internal server error: Failed to get user repository from ioc container", err)
	}

	user, err := userRepo.FindById(val)

	if err != nil {
		return nil, errors.NewUnAuthorizedError("unauthenticated", "User not found", err)
	}

	return user, nil
}
