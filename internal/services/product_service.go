package services

import (
	"errors"
	"taskgo/internal/api/requests"
	"taskgo/internal/database/models"
	"taskgo/internal/filters"
	"taskgo/internal/repository"
	pkgErrors "taskgo/pkg/errors"
	"taskgo/pkg/utils"

	"gorm.io/gorm"
)

type ProductService struct {
	productRepository *repository.ProductRepository
}

// Create a new product service
func NewProductService(userRepository *repository.ProductRepository) *ProductService {
	return &ProductService{productRepository: userRepository}
}

// Create a new product
func (s *ProductService) CreateProduct(req *requests.CreateProductRequest) (*models.Product, error) {
	product := &models.Product{
		Name:        req.Name,
		Price:       req.Price,
		Description: req.Description,
		Status:      req.Status,
		Category:    req.Category,
		Brand:       req.Brand,
		Weight:      req.Weight,
		WeightUnit:  req.WeightUnit,
		Attributes:  req.Attributes,
	}

	err := s.productRepository.Create(product)
	if err != nil {
		return nil, err
	}
	return product, nil
}

// Updates a product
func (s *ProductService) UpdateProduct(productId string, req *requests.UpdateProductRequest) (*models.Product, error) {

	updatedFields := make(map[string]any)
	sentFields := req.GetRequestSentFields()
	validKeys := utils.GetJSONKeys(req)

	for key, value := range sentFields {
		if validKeys[key] {
			updatedFields[key] = value
		}
	}

	product, err := s.productRepository.UpdateByIdAndGet(productId, updatedFields)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, pkgErrors.NewNotFoundError("product not found", "product not found", err)
		}
		return nil, err
	}
	return product, nil
}

// Get a product by id
func (s *ProductService) GetProductById(id string) (*models.Product, error) {
	product, err := s.productRepository.FindById(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, pkgErrors.NewNotFoundError("product not found", "product not found", err)
		}
		return nil, err
	}
	return product, nil
}

// Get paginated products
func (s *ProductService) GetPaginatedProducts(productFilters *filters.ProductFilters) ([]*models.Product, int64, error) {
	products, total, err := s.productRepository.Paginate(productFilters)
	if err != nil {
		return nil, 0, err
	}
	return products, total, nil
}
