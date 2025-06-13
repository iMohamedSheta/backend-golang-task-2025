package policies

import (
	"fmt"
	"taskgo/internal/database/models"
	"taskgo/internal/enums"
)

type UserPolicy struct {
}

// Check if user can create another user
func (p *UserPolicy) CanCreate(user *models.User) bool {
	return true
}

// Check if user can view another user
func (p *UserPolicy) CanView(user *models.User, userId string) bool {
	return user.Role == enums.RoleAdmin || fmt.Sprint(user.ID) == userId
}

// Check if user can update another user
func (p *UserPolicy) CanUpdate(user *models.User, userId string) bool {
	return user.Role == enums.RoleAdmin || fmt.Sprint(user.ID) == userId
}
