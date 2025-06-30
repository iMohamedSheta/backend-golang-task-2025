package handlers

import (
	"math"
	"taskgo/internal/api/requests"
	"taskgo/internal/api/responses"
	"taskgo/internal/database/models"
	"taskgo/internal/deps"
	"taskgo/internal/filters"
	"taskgo/internal/helpers"
	"taskgo/internal/policies"
	"taskgo/internal/services"
	"taskgo/pkg/errors"
	"taskgo/pkg/response"

	"github.com/gin-gonic/gin"
)

type ProductHandler struct {
	Handler
	productService *services.ProductService
	productPolicy  *policies.ProductPolicy
}

// NewProductHandler return a new ProductHandler
func NewProductHandler(productService *services.ProductService, productPolicy *policies.ProductPolicy) *ProductHandler {
	return &ProductHandler{
		productService: productService,
		productPolicy:  productPolicy,
	}
}

// @Summary     List products
// @Description Retrieves a paginated list of products.
// @Tags        Products
// @Accept      json
// @Produce     json
//
// @Param       request    query     filters.ProductFilters           true  "Filter and pagination"
//
// @Success     200        {object}  responses.ListProductsResponse   "Success"
// @Failure     400        {object}  response.BadRequestResponse      "Bad Request"
// @Failure     401        {object}  response.UnauthorizedResponse    "Unauthorized Action"
// @Failure     500        {object}  response.ServerErrorResponse     "Internal Server Error"
//
// @Router      /products [get]
func (h *ProductHandler) ListProducts(gin *gin.Context) error {
	var productFilters filters.ProductFilters

	if !h.productPolicy.CanViewAny(nil) {
		return errors.NewUnAuthorizedError("UnauthorizedError: You are not allowed to view this product", "user is not allowed to view this product", nil)
	}

	// Bind URL query parameters to filters struct
	if err := gin.ShouldBindQuery(&productFilters); err != nil {
		return errors.NewBadRequestError("", "BadRequestError: Failed to bind URL query parameters to filters struct", err)
	}

	products, total, err := h.productService.GetPaginatedProducts(gin.Request.Context(), &productFilters)
	if err != nil {
		return errors.NewServerError("internal server error", "Err: Failed to get paginated products using productService", err)
	}

	var totalPages int
	if productFilters.PerPage > 0 {
		totalPages = int(math.Ceil(float64(total) / float64(productFilters.PerPage)))
	}

	responses.SendListProductsResponse(gin, products, responses.PaginationMeta{
		Total:      total,
		Page:       productFilters.Page,
		Limit:      productFilters.PerPage,
		NextPage:   productFilters.Page + 1,
		PrevPage:   productFilters.Page - 1,
		TotalPages: totalPages,
	})

	return nil
}

// @Summary     Get product by ID
// @Description Retrieves a product by its ID.
// @Tags        Products
// @Accept      json
// @Produce     json
//
// @Param       id         path      string                             true  "Product ID"
//
// @Success     200        {object}  responses.GetProductResponse      "Product details retrieved successfully"
// @Failure     400        {object}  response.BadRequestResponse       "Bad Request"
// @Failure     404        {object}  response.NotFoundResponse         "Product not found"
// @Failure     401        {object}  response.UnauthorizedResponse     "Unauthorized Action"
// @Failure     500        {object}  response.ServerErrorResponse      "Internal Server Error"
//
// @Router      /products/{id} [get]
func (h *ProductHandler) GetProduct(gin *gin.Context) error {
	productId := gin.Param("id")

	if !h.productPolicy.CanView(nil, productId) {
		return errors.NewUnAuthorizedError("UnauthorizedError: You are not allowed to view this product", "UnauthorizedError: user is not allowed to view this product", nil)
	}

	targetProduct, err := h.productService.GetProductById(gin.Request.Context(), productId)
	if err != nil {
		return err
	}

	responses.SendGetProductResponse(gin, targetProduct)
	return nil
}

// @Summary     Create a new product
// @Description Creates a new product with the given details.
// @Tags        Products
// @Accept      json
// @Produce     json
//
// @Param       request  body      requests.CreateProductRequest       true  "Create product request body"
//
// @Success     201      {object}  responses.CreateProductResponse     "Product created successfully"
// @Failure     400      {object}  response.BadRequestResponse         "Bad Request"
// @Failure     401      {object}  response.UnauthorizedResponse       "Unauthorized Action"
// @Failure     422      {object}  response.ValidationErrorResponse    "Validation Error"
// @Failure     500      {object}  response.ServerErrorResponse        "Internal Server Error"
//
// @Router      /products [post]
func (h *ProductHandler) CreateProduct(gin *gin.Context) error {
	var req requests.CreateProductRequest
	var err error

	if err = h.BindBodyAndExtractToRequest(gin, &req); err != nil {
		return errors.NewBadRequestBindingError("", "BadRequestBindingError: Failed to bind request body to request struct", err)
	}

	authUser, authorizeErr := helpers.GetAuthUser(gin)
	if authorizeErr != nil {
		return authorizeErr
	}

	if !h.productPolicy.CanCreate(authUser) {
		return errors.NewUnAuthorizedError("Unauthorized", "You are not allowed to create a product", nil)
	}

	if err = deps.Validator().ValidateRequest(&req); err != nil {
		return err
	}

	var product *models.Product
	product, err = h.productService.CreateProduct(gin.Request.Context(), &req)
	if err != nil {
		return errors.NewServerError("Internal Server Error: Failed to create a new product", "Internal Server Error: Failed to create a new product", err)
	}

	responses.SendCreateProductResponse(gin, product)
	return nil
}

// @Summary     Update product by ID
// @Description Updates a product with the given ID and request body.
// @Tags        Products
// @Accept      json
// @Produce     json
//
// @Param       id       path      string                             true  "Product ID"
// @Param       request  body      requests.UpdateProductRequest      true  "Update product request body"
//
// @Success     200      {object}  responses.UpdateProductResponse    "Product updated successfully"
// @Failure     400      {object}  response.BadRequestResponse        "Bad Request"
// @Failure     401      {object}  response.UnauthorizedResponse      "Unauthorized Action"
// @Failure     404      {object}  response.NotFoundResponse          "Product Not Found"
// @Failure     422      {object}  response.ValidationErrorResponse   "Validation Error"
// @Failure     500      {object}  response.ServerErrorResponse       "Internal Server Error"
//
// @Router      /products/{id} [put]
func (h *ProductHandler) UpdateProduct(gin *gin.Context) error {
	var req requests.UpdateProductRequest
	var err error

	if err = h.BindBodyAndExtractToRequest(gin, &req); err != nil {
		return errors.NewBadRequestBindingError("", "BadRequestBindingError: Failed to bind request body to request struct", err)
	}

	productId := gin.Param("id")
	authUser, authorizeErr := helpers.GetAuthUser(gin)
	if authorizeErr != nil {
		return authorizeErr
	}

	if !h.productPolicy.CanUpdate(authUser, productId) {
		return errors.NewUnAuthorizedError("", "You are not allowed to update the product", nil)
	}

	if err = deps.Validator().ValidateRequest(&req); err != nil {
		return err
	}

	product, err := h.productService.UpdateProduct(gin.Request.Context(), productId, &req)
	if err != nil {
		return err
	}

	responses.SendUpdateProductResponse(gin, product)
	return nil
}

// TODO: add inventory check implementation
func (h *ProductHandler) CheckInventory(c *gin.Context) {
	response.Json(c, "Inventory checked successfully", gin.H{
		"product_id": c.Param("id"),
		"stock":      100,
	}, 200)
}
