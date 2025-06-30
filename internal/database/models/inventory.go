package models

import "fmt"

type Inventory struct {
	Base
	ProductID     uint    `gorm:"uniqueIndex;not null" json:"product_id"` // foreign key - index
	Quantity      int     `gorm:"not null;default:0" json:"quantity"`
	ReorderPoint  int     `gorm:"not null;default:0" json:"reorder_point"`  // minimum quantity to reorder
	ReorderAmount int     `gorm:"not null;default:0" json:"reorder_amount"` // quantity to reorder
	Location      string  `gorm:"size:100" json:"location"`
	LastRestocked string  `gorm:"size:100" json:"last_restocked"`
	UnitCost      float64 `gorm:"type:decimal(10,2)" json:"unit_cost"`
}

// TODO: admin should be able to set (minimum quantity to reorder) and should implement notification for that
func (i *Inventory) NeedsRestock() bool {
	return i.Quantity <= i.ReorderPoint
}

func (i *Inventory) GetInventoryCacheKey() string {
	return fmt.Sprintf("product:%d:inventory:%d", i.ProductID, i.ID)
}

func (i *Inventory) GetTotalQuantity() int {
	return 500
}
