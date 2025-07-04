package models

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"taskgo/internal/deps"
	"taskgo/internal/enums"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ProductAttributes map[string]any

// Value implements the driver.Valuer interface for ProductAttributes
func (a ProductAttributes) Value() (driver.Value, error) {
	return json.Marshal(a)
}

// Scan implements the sql.Scanner interface for ProductAttributes
func (a *ProductAttributes) Scan(value any) error {
	bytes, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}
	return json.Unmarshal(bytes, &a)
}

type Product struct {
	Base
	Name        string              `gorm:"not null" json:"name"`
	Description string              `gorm:"type:text" json:"description"`
	SKU         string              `gorm:"uniqueIndex;not null" json:"sku"` // index - stock keeping unit (unique identifier for the product)
	Price       float64             `gorm:"type:decimal(10,2);not null" json:"price"`
	Status      enums.ProductStatus `gorm:"type:varchar(20);not null;default:'available'" json:"status"`
	Attributes  ProductAttributes   `gorm:"type:jsonb" json:"attributes"`
	Category    string              `gorm:"index" json:"category"` // index
	Brand       string              `gorm:"index" json:"brand"`    // index
	Weight      float64             `gorm:"type:decimal(10,2)" json:"weight"`
	WeightUnit  string              `gorm:"size:10" json:"weight_unit"`
	Inventories []Inventory         `gorm:"foreignKey:ProductID" json:"inventory"` // relationship product (1) to inventory (m)
}

func (p *Product) LoadInventories() error {
	return deps.Gorm().DB.Model(p).Association("Inventories").Find(&p.Inventories)
}

// GenerateSKU generates a unique SKU for the product/
func (p *Product) GenerateSKU(prefix string) string {
	if prefix == "" {
		prefix = "SKU"
	}
	// Stock keeping unit logic will be added here later based on the way we want to track the product
	return strings.ToUpper(fmt.Sprintf("%s_%s", prefix, uuid.New().String()[:8]))
}

// before create hook
func (p *Product) BeforeCreate(tx *gorm.DB) error {
	if p.SKU != "" {
		return nil
	}
	p.SKU = p.GenerateSKU("SKU")
	return nil
}

func (p *Product) GetInventoryCacheKeyForAuthUser(authUserID uint) string {
	// Should do actual logic to get the inventory cache key based on the auth user location or other logic
	// But i will just get first inventory for now
	var firstInventory Inventory
	deps.Gorm().DB.Model(&Inventory{}).
		Where("product_id = ?", p.ID).
		Select("id, product_id, quantity").
		First(&firstInventory)

	return firstInventory.GetInventoryCacheKey()
}
