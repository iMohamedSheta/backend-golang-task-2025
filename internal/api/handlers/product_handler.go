package handlers

import (
	"taskgo/internal/api/requests"
	"taskgo/internal/filters"
	"taskgo/internal/helpers"
	"taskgo/internal/policies"
	"taskgo/internal/repository"
	"taskgo/internal/services"
	"taskgo/pkg/errors"
	"taskgo/pkg/logger"
	"taskgo/pkg/response"
	"taskgo/pkg/utils"
	"taskgo/pkg/validate"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
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

func (h *ProductHandler) ListProducts(c *gin.Context) {
	var ProductFilters filters.ProductFilters

	if !h.productPolicy.CanViewAny(nil) {
		response.UnauthorizedJson(c, errors.NewUnAuthorizedError("You are not allowed to view this product", "user is not allowed to view this product", nil))
		return
	}

	// Bind URL query parameters to filters struct
	if err := c.ShouldBindQuery(&ProductFilters); err != nil {
		logger.Log().Error("Failed to bind product filters", zap.Error(err))
		response.BadRequestBindingJson(c, err)
		return
	}

	products, total, err := h.productService.GetPaginatedProducts(&ProductFilters)
	if err != nil {
		response.ServerErrorJson(c, errors.NewServerError("Failed to get products", "failed to get products", err))
		return
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
}

func (h *ProductHandler) GetProduct(c *gin.Context) {
	productId := c.Param("id")

	// Policy check
	if !h.productPolicy.CanView(nil, productId) {
		response.UnauthorizedJson(c, errors.NewUnAuthorizedError("You are not allowed to view this product", "user is not allowed to view this product", nil))
		return
	}

	targetProduct, err := h.productService.GetProductById(productId)
	if err != nil {
		response.NotFoundJson(c, errors.NewNotFoundError("Product not found", "product with id "+productId+" not found", nil))
		return
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
}

// Create New Product (only admins)
func (h *ProductHandler) CreateProduct(c *gin.Context) {
	var req requests.CreateProductRequest

	if err := utils.BindToRequestAndExtractFields(c, &req); err != nil {
		logger.Log().Error("Failed to bind create product request", zap.Error(err))
		response.BadRequestBindingJson(c, err)
		return
	}

	authUser, authorizeErr := helpers.GetAuthUser(c)
	if authorizeErr != nil {
		response.UnauthorizedJson(c, authorizeErr)
		return
	}

	if !h.productPolicy.CanCreate(authUser) {
		response.UnauthorizedJson(c, errors.NewUnAuthorizedError("Unauthorized", "You are not allowed to create a product", nil))
		return
	}

	valid, validErrorsList := validate.ValidateRequest(&req)
	if !valid {
		validErrors := errors.NewValidationError(validErrorsList)
		response.ValidationErrorJson(c, validErrors)
		return
	}

	product, err := h.productService.CreateProduct(&req)
	if err != nil {
		logger.Log().Error("Failed to create product", zap.Error(err))
		response.ServerErrorJson(c, errors.NewServerError("Failed to create product", "Failed to create product", err))
		return
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
}

func (h *ProductHandler) UpdateProduct(c *gin.Context) {
	var req requests.UpdateProductRequest

	// This binds and stores the body so it can be reused
	if err := utils.BindToRequestAndExtractFields(c, &req); err != nil {
		logger.Log().Error("Failed to bind update product request", zap.Error(err))
		response.BadRequestBindingJson(c, err)
		return
	}

	productId := c.Param("id")
	authUser, authorizeErr := helpers.GetAuthUser(c)
	if authorizeErr != nil {
		response.UnauthorizedJson(c, authorizeErr)
		return
	}

	if !h.productPolicy.CanUpdate(authUser, productId) {
		response.UnauthorizedJson(c, errors.NewUnAuthorizedError("Unauthorized", "You are not allowed to update the product", nil))
		return
	}

	valid, validErrorsList := validate.ValidateRequest(&req)
	if !valid {
		validErrors := errors.NewValidationError(validErrorsList)
		response.ValidationErrorJson(c, validErrors)
		return
	}

	product, err := h.productService.UpdateProduct(productId, &req)
	if err != nil {
		if notfoundErr, ok := errors.AsNotFoundError(err); ok {
			response.NotFoundJson(c, notfoundErr)
			return
		}
		logger.Log().Error("Failed to update product", zap.Error(err))
		response.ServerErrorJson(c, errors.NewServerError("Failed to update product", "Failed to update product", err))
		return
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
}

// TODO: add inventory check implementation
func (h *ProductHandler) CheckInventory(c *gin.Context) {
	response.Json(c, "Inventory checked successfully", gin.H{
		"product_id": c.Param("id"),
		"stock":      100,
	}, 200)
}
