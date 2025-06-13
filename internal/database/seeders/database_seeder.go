package seeders

// SeedDatabase seeds the database
func SeedDatabase() {
	// Seed admin user
	SeedAdminUser()

	// Seed products
	SeedProducts()

	// Seed Inventory
	SeedInventory()
}
