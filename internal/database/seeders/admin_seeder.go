package seeders

import (
	"log"
	"taskgo/internal/config"
	"taskgo/internal/database/models"
	"taskgo/internal/enums"
	"taskgo/pkg/database"
	pkgEnums "taskgo/pkg/enums"

	"gorm.io/gorm"
)

func SeedAdminUser() {
	db := database.GetDB()

	adminData, err := config.App.Get("app.admins")
	if err != nil {
		log.Fatalf("Failed to get admin user data from config: %v", err)
		return
	}

	adminsSlice, ok := adminData.([]map[string]any)
	if !ok {
		log.Fatalf("Admin user data should be an array of map[string]any, got: %T", adminData)
		return
	}

	for _, adminMap := range adminsSlice {
		email, _ := adminMap["email"].(string)
		password, _ := adminMap["password"].(string)
		firstName, _ := adminMap["first_name"].(string)
		lastName, _ := adminMap["last_name"].(string)

		if email == "" {
			log.Println(pkgEnums.Yellow.Value() + "Skipping admin user seeding due to empty email" + pkgEnums.Reset.Value())
			continue
		}
		if password == "" {
			log.Println(pkgEnums.Yellow.Value() + "Skipping admin user seeding due to empty password" + pkgEnums.Reset.Value())
			continue
		}
		if firstName == "" {
			firstName = "Admin"
		}
		if lastName == "" {
			lastName = "User"
		}

		var existingUser models.User
		err := db.Where("email = ?", email).First(&existingUser).Error
		if err == nil {
			log.Println(pkgEnums.Yellow.Value() + "Admin user already exists: " + email + pkgEnums.Reset.Value())
			continue
		} else if err != gorm.ErrRecordNotFound {
			log.Fatalf("Failed to query admin user %s: %v", email, err)
			return
		}

		log.Println(pkgEnums.Green.Value() + "Seeding admin user: " + email + pkgEnums.Reset.Value())

		adminUser := models.User{
			Email:     email,
			Password:  password, // hashed in User model BeforeCreate hook
			FirstName: firstName,
			LastName:  lastName,
			Role:      enums.RoleAdmin,
			IsActive:  true,
		}

		if err := db.Create(&adminUser).Error; err != nil {
			log.Fatalf("Failed to seed admin user %s: %v", email, err)
		}

		log.Println(pkgEnums.Green.Value() + "Admin user seeded successfully: " + email + pkgEnums.Reset.Value())
	}
}
