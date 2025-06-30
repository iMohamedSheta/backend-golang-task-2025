package seeders

import (
	"fmt"
	"log"
	"taskgo/internal/database/models"
	"taskgo/internal/deps"
	pkgEnums "taskgo/pkg/enums"
)

func SeedProducts() {
	db := deps.Gorm().DB

	products := []models.Product{
		{
			Name:        "Wireless Mouse",
			Description: "Ergonomic wireless mouse with USB receiver",
			Price:       25.99,
			Category:    "Electronics",
			Brand:       "LogiTech",
			Weight:      0.2,
			WeightUnit:  "kg",
		},
		{
			Name:        "Mechanical Keyboard",
			Description: "RGB backlit mechanical keyboard with blue switches",
			Price:       89.99,
			Category:    "Electronics",
			Brand:       "KeyPro",
			Weight:      0.9,
			WeightUnit:  "kg",
		},
		{
			Name:        "USB-C Charger",
			Description: "Fast charging USB-C power adapter 65W",
			Price:       39.99,
			Category:    "Accessories",
			Brand:       "ChargeX",
			Weight:      0.15,
			WeightUnit:  "kg",
		},
	}

	for _, product := range products {
		if err := db.Create(&product).Error; err != nil {
			fmt.Printf("Failed to seed product '%s': %v\n", product.Name, err)
		}
	}

	log.Println(pkgEnums.Green.Value() + "Seeding products completed successfully" + pkgEnums.Reset.Value())
}
