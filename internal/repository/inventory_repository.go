package repository

import (
	"taskgo/internal/database/models"
	"taskgo/internal/deps"
)

type InventoryRepository struct {
	db *deps.GormDB
}

func NewInventoryRepository(db *deps.GormDB) *InventoryRepository {
	return &InventoryRepository{
		db: db,
	}
}

func (r *InventoryRepository) UpdateQuantity(inventoryId uint, quantity int) error {
	return r.db.DB.Model(&models.Inventory{}).Where("id = ?", inventoryId).Update("quantity", quantity).Error
}
