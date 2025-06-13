package policies

import (
	"taskgo/internal/database/models"
)

type OrderPolicy struct {
}

// Check if the user can create a product
func (p *OrderPolicy) CanCreate(user *models.User) bool {
	if user != nil { // if user is authenticated
		return true
	}
	return false
}
