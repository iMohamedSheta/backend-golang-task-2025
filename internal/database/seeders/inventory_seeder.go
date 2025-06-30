package seeders

import (
	"fmt"
	"log"
	"taskgo/internal/database/models"
	"taskgo/internal/deps"
	pkgEnums "taskgo/pkg/enums"
	"time"
)

func SeedInventory() {
	db := deps.Gorm().DB

	var products []models.Product
	if err := db.Find(&products).Error; err != nil {
		fmt.Printf("Failed to load products: %v\n", err)
		return
	}

	for _, product := range products {
		inventory := models.Inventory{
			ProductID:     product.ID,
			Quantity:      100,
			ReorderPoint:  20,
			ReorderAmount: 50,
			Location:      "Main Warehouse",
			LastRestocked: time.Now().Format("2006-01-02"),
			UnitCost:      10.50,
		}

		if err := db.Create(&inventory).Error; err != nil {
			fmt.Printf("Failed to seed inventory for product %d: %v\n", product.ID, err)
		}
	}

	log.Println(pkgEnums.Green.Value() + "Seeding inventory completed successfully" + pkgEnums.Reset.Value())
}
