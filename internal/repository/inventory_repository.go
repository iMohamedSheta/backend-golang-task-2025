package repository

import (
	"taskgo/internal/database/models"
	"taskgo/pkg/database"

	"gorm.io/gorm"
)

type InventoryRepository struct {
	db *gorm.DB
}

func NewInventoryRepository() *InventoryRepository {
	db := database.GetDB()
	return &InventoryRepository{
		db: db,
	}
}

func (r *InventoryRepository) UpdateQuantity(inventoryId uint, quantity int) error {
	return r.db.Model(&models.Inventory{}).Where("id = ?", inventoryId).Update("quantity", quantity).Error
}
