package policies

import (
	"taskgo/internal/database/models"
	"taskgo/internal/enums"
)

type ProductPolicy struct {
}

// Check if the user can view any product
func (p *ProductPolicy) CanViewAny(user *models.User) bool {
	return true
}

// Check if the user can view a product
func (p *ProductPolicy) CanView(user *models.User, productId string) bool {
	return true
}

// Check if the user can create a product
func (p *ProductPolicy) CanCreate(user *models.User) bool {
	return user.Role == enums.RoleAdmin
}

// Check if the user can update a product
func (p *ProductPolicy) CanUpdate(user *models.User, productId string) bool {
	return user.Role == enums.RoleAdmin
}
