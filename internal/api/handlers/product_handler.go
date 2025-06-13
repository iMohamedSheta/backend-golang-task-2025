package handlers

import (
	"taskgo/internal/api/requests"
	"taskgo/internal/filters"
	"taskgo/internal/helpers"
	"taskgo/internal/policies"
	"taskgo/internal/repository"
	"taskgo/internal/services"
	"taskgo/pkg/errors"
	"taskgo/pkg/response"
	"taskgo/pkg/utils"
	"taskgo/pkg/validate"

	"github.com/gin-gonic/gin"
)

type ProductHandler struct {
	productService *services.ProductService
	productPolicy  *policies.ProductPolicy
}

func NewProductHandlerWithDeps(productService *services.ProductService, productPolicy *policies.ProductPolicy) *ProductHandler {
	return &ProductHandler{
		productService: productService,
		productPolicy:  productPolicy,
	}
}

func NewProductHandler() *ProductHandler {
	productRepo := repository.NewProductRepository()
	productService := services.NewProductService(productRepo)
	return &ProductHandler{
		productService: productService,
		productPolicy:  &policies.ProductPolicy{},
	}
}

func (h *ProductHandler) ListProducts(c *gin.Context) error {
	var ProductFilters filters.ProductFilters

	if !h.productPolicy.CanViewAny(nil) {
		return errors.NewUnAuthorizedError("UnauthorizedError: You are not allowed to view this product", "user is not allowed to view this product", nil)
	}

	// Bind URL query parameters to filters struct
	if err := c.ShouldBindQuery(&ProductFilters); err != nil {
		return errors.NewBadRequestError("", "BadRequestError: Failed to bind URL query parameters to filters struct", err)
	}

	products, total, err := h.productService.GetPaginatedProducts(&ProductFilters)
	if err != nil {
		return errors.NewServerError("internal server error", "Err: Failed to get paginated products using productService", err)
	}

	var productResources []map[string]any

	for _, product := range products {
		productResource := map[string]any{
			"id":          product.ID,
			"name":        product.Name,
			"description": product.Description,
			"price":       product.Price,
			"category":    product.Category,
			"status":      product.Status,
			"sku":         product.SKU,
			"attributes":  product.Attributes,
			"brand":       product.Brand,
			"weight":      product.Weight,
			"weight_unit": product.WeightUnit,
			"created_at":  product.CreatedAt,
			"updated_at":  product.UpdatedAt,
		}
		productResources = append(productResources, productResource)
	}

	response.Json(c, "Products retrieved successfully", map[string]any{
		"products": productResources,
		"meta": map[string]any{
			"total": total,
		},
	}, 200)

	return nil
}

func (h *ProductHandler) GetProduct(c *gin.Context) error {
	productId := c.Param("id")

	// Policy check
	if !h.productPolicy.CanView(nil, productId) {
		return errors.NewUnAuthorizedError("UnauthorizedError: You are not allowed to view this product", "UnauthorizedError: user is not allowed to view this product", nil)
	}

	targetProduct, err := h.productService.GetProductById(productId)
	if err != nil {
		if notFoundErr, ok := errors.AsNotFoundError(err); ok {
			return notFoundErr
		}

		return errors.NewServerError("internal server error", "Err: Failed to get product by id", err)
	}

	response.Json(c, "Product retrieved successfully", map[string]any{
		"product": map[string]any{
			"id":          targetProduct.ID,
			"name":        targetProduct.Name,
			"description": targetProduct.Description,
			"price":       targetProduct.Price,
			"category":    targetProduct.Category,
			"status":      targetProduct.Status,
			"sku":         targetProduct.SKU,
			"weight":      targetProduct.Weight,
			"weight_unit": targetProduct.WeightUnit,
			"attributes":  targetProduct.Attributes,

			"created_at": targetProduct.CreatedAt,
			"updated_at": targetProduct.UpdatedAt,
		},
	}, 200)

	return nil
}

// Create New Product (only admins)
func (h *ProductHandler) CreateProduct(c *gin.Context) error {
	var req requests.CreateProductRequest

	if err := utils.BindToRequestAndExtractFields(c, &req); err != nil {
		return errors.NewBadRequestBindingError("", "BadRequestBindingError: Failed to bind request body to request struct", err)
	}

	authUser, authorizeErr := helpers.GetAuthUser(c)
	if authorizeErr != nil {
		return authorizeErr
	}

	if !h.productPolicy.CanCreate(authUser) {
		return errors.NewUnAuthorizedError("Unauthorized", "You are not allowed to create a product", nil)
	}

	valid, validErrorsList := validate.ValidateRequest(&req)
	if !valid {
		return errors.NewValidationError(validErrorsList)
	}

	product, err := h.productService.CreateProduct(&req)
	if err != nil {
		return errors.NewServerError("Failed to create product", "Failed to create product", err)
	}

	response.Json(c, "Product created successfully", map[string]any{
		"product": map[string]any{
			"id":         product.ID,
			"name":       product.Name,
			"price":      product.Price,
			"sku":        product.SKU,
			"brand":      product.Brand,
			"created_at": product.CreatedAt,
			"updated_at": product.UpdatedAt,
		},
	}, 201)

	return nil
}

func (h *ProductHandler) UpdateProduct(c *gin.Context) error {
	var req requests.UpdateProductRequest

	// This binds and stores the body so it can be reused
	if err := utils.BindToRequestAndExtractFields(c, &req); err != nil {
		return errors.NewBadRequestBindingError("", "BadRequestBindingError: Failed to bind request body to request struct", err)
	}

	productId := c.Param("id")
	authUser, authorizeErr := helpers.GetAuthUser(c)
	if authorizeErr != nil {
		return authorizeErr
	}

	if !h.productPolicy.CanUpdate(authUser, productId) {
		return errors.NewUnAuthorizedError("Unauthorized", "You are not allowed to update the product", nil)
	}

	valid, validErrorsList := validate.ValidateRequest(&req)
	if !valid {
		return errors.NewValidationError(validErrorsList)
	}

	product, err := h.productService.UpdateProduct(productId, &req)
	if err != nil {
		if notfoundErr, ok := errors.AsNotFoundError(err); ok {
			return notfoundErr
		}

		return errors.NewServerError("Failed to update product", "Failed to update product", err)
	}

	response.Json(c, "Product updated successfully", map[string]any{
		"product": map[string]any{
			"id":         product.ID,
			"name":       product.Name,
			"price":      product.Price,
			"sku":        product.SKU,
			"brand":      product.Brand,
			"created_at": product.CreatedAt,
			"updated_at": product.UpdatedAt,
		},
	}, 201)

	return nil
}

// TODO: add inventory check implementation
func (h *ProductHandler) CheckInventory(c *gin.Context) {
	response.Json(c, "Inventory checked successfully", gin.H{
		"product_id": c.Param("id"),
		"stock":      100,
	}, 200)
}
