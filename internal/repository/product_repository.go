package repository

import (
	"errors"
	"taskgo/internal/database/models"
	"taskgo/internal/enums"
	"taskgo/internal/filters"
	"taskgo/pkg/database"
	"taskgo/pkg/utils"

	"gorm.io/gorm"
)

type ProductRepository struct {
	db *gorm.DB
}

func NewProductRepository() *ProductRepository {
	db := database.GetDB()
	return &ProductRepository{
		db: db,
	}
}

// Create a new product
func (r *ProductRepository) Create(product *models.Product) error {
	return r.db.Create(product).Error
}

// Update a product by id
func (r *ProductRepository) UpdateById(id string, data map[string]interface{}) error {
	if id == "" {
		return errors.New("id is required")
	}

	if err := r.db.Model(&models.Product{}).Where("id = ?", id).Updates(data).Error; err != nil {
		return err
	}

	return nil
}

// Update a product and return the updated product
func (r *ProductRepository) UpdateByIdAndGet(id string, data map[string]interface{}) (*models.Product, error) {
	err := r.UpdateById(id, data)
	if err != nil {
		return nil, err
	}

	// Fetch the updated user to return the latest state
	var product models.Product
	if err := r.db.Where("id = ?", id).First(&product).Error; err != nil {
		return nil, err
	}

	return &product, nil
}

// Get a product by id
func (r *ProductRepository) FindById(id string) (*models.Product, error) {
	if id == "" {
		return nil, errors.New("id is required")
	}

	product := models.Product{}
	if err := r.db.Where("id = ?", id).First(&product).Error; err != nil {
		return nil, err
	}

	return &product, nil
}

// CheckIDsExist checks a list of ids exists in db and returns the ids that don't exist
func (r *ProductRepository) CheckIDsExist(ids []uint) ([]uint, error) {
	if len(ids) == 0 {
		return []uint{}, nil
	}

	var existingIDs []uint
	err := r.db.Model(&models.Product{}).
		Where("id IN ?", ids).
		Pluck("id", &existingIDs).Error

	if err != nil {
		return nil, err
	}

	// Find missing IDs by comparing input with existing
	existingMap := make(map[uint]bool)
	for _, id := range existingIDs {
		existingMap[id] = true
	}

	var missingIDs []uint
	for _, id := range ids {
		if !existingMap[id] {
			missingIDs = append(missingIDs, id)
		}
	}

	return missingIDs, nil
}

// Get a list of products by list of ids
func (r *ProductRepository) FindByIDs(ids []uint) ([]models.Product, error) {
	var products []models.Product
	err := r.db.Where("id IN ?", ids).Find(&products).Error
	return products, err
}

// Get a list of products with inventory by list of ids
func (r *ProductRepository) FindByIDsWithInventory(ids []uint, inventoryColumns ...string) ([]models.Product, error) {
	var products []models.Product

	validMap := utils.GetJSONKeys(models.Inventory{})

	// Use "*" if no columns passed
	if len(inventoryColumns) == 0 {
		inventoryColumns = []string{"*"}
	}

	// Filter valid columns
	filtered := make([]string, 0, len(inventoryColumns))
	hasWildcard := false
	for _, col := range inventoryColumns {
		if col == "*" {
			hasWildcard = true
			break
		}
		if validMap[col] {
			filtered = append(filtered, col)
		}
	}

	// If "*" provided or no columns provided, use "*"
	if hasWildcard || len(filtered) == 0 {
		filtered = []string{"*"}
	} else {
		// Ensure "product_id" is included for join
		hasProductID := false
		for _, col := range filtered {
			if col == "product_id" {
				hasProductID = true
				break
			}
		}
		if !hasProductID {
			filtered = append(filtered, "product_id")
		}
	}

	err := r.db.
		Preload("Inventory", func(db *gorm.DB) *gorm.DB {
			return db.Select(filtered)
		}).
		Where("id IN ?", ids).
		Find(&products).Error

	return products, err
}

// Paginate products with filters
func (r *ProductRepository) Paginate(productFilters *filters.ProductFilters) ([]*models.Product, int64, error) {
	var products []*models.Product
	var total int64

	db := r.db.Model(&models.Product{})

	db = r.applyFilters(db, productFilters)

	// Get total count before pagination
	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	db = r.applySorting(db, productFilters)

	// Set default values
	if productFilters.Page <= 0 {
		productFilters.Page = 1
	}

	if productFilters.PerPage <= 0 {
		productFilters.PerPage = 10
	}

	// Apply pagination
	offset := (productFilters.Page - 1) * productFilters.PerPage
	_ = db.Offset(offset).Limit(productFilters.PerPage)

	if err := db.Find(&products).Error; err != nil {
		return nil, 0, err
	}

	return products, total, nil
}

// applyFilters applies all the filters to the query
func (r *ProductRepository) applyFilters(db *gorm.DB, filters *filters.ProductFilters) *gorm.DB {
	if filters.Name != "" {
		db = db.Where("name ILIKE ?", "%"+filters.Name+"%")
	}

	if filters.SKU != "" {
		db = db.Where("sku = ?", filters.SKU)
	}

	if filters.Brand != "" {
		db = db.Where("brand ILIKE ?", "%"+filters.Brand+"%")
	}

	if filters.Category != "" {
		db = db.Where("category ILIKE ?", filters.Category)
	}

	if filters.Search != "" {
		searchFields := filters.GetSearchFields()
		conditions := ""
		args := []interface{}{}
		searchTerm := "%" + filters.Search + "%"

		for i, field := range searchFields {
			if i > 0 {
				conditions += " OR "
			}
			conditions += field + " ILIKE ?"
			args = append(args, searchTerm)
		}

		db = db.Where(conditions, args...)
	}

	// Price range filters
	if filters.MinPrice != nil && *filters.MinPrice > 0 {
		db = db.Where("price >= ?", *filters.MinPrice)
	}

	if filters.MaxPrice != nil && *filters.MaxPrice > 0 {
		db = db.Where("price <= ?", *filters.MaxPrice)
	}

	if filters.Status != nil && enums.IsValidProductStatus(*filters.Status) {
		db = db.Where("status = ?", *filters.Status)
	}

	if filters.CreatedAfter != nil && !filters.CreatedAfter.IsZero() {
		db = db.Where("created_at >= ?", *filters.CreatedAfter)
	}

	if filters.CreatedBefore != nil && !filters.CreatedBefore.IsZero() {
		db = db.Where("created_at <= ?", *filters.CreatedBefore)
	}

	return db
}

// applySorting applies sorting to the query
func (r *ProductRepository) applySorting(db *gorm.DB, f *filters.ProductFilters) *gorm.DB {
	sortBy := f.SortBy
	if sortBy == "" {
		sortBy = "created_at"
	}

	sortOrder := f.SortOrder
	if sortOrder == "" {
		sortOrder = "desc"
	}

	// Validate sort fields to prevent SQL injection
	sortFields := f.GetSortFields()
	validSortFields := make(map[string]bool)
	for _, field := range sortFields {
		validSortFields[field] = true
	}

	if !validSortFields[sortBy] {
		sortBy = "created_at"
	}

	if sortOrder != "asc" && sortOrder != "desc" {
		sortOrder = "desc"
	}

	return db.Order(sortBy + " " + sortOrder)
}
