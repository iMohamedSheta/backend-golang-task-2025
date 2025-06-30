package responses

import (
	"net/http"
	"taskgo/internal/database/models"
	"taskgo/pkg/response"

	"github.com/gin-gonic/gin"
)

type ListProductsResponse struct {
	Message string `json:"message" example:"Products retrieved successfully"`
	Data    struct {
		Products []ProductData  `json:"products"`
		Meta     PaginationMeta `json:"meta"`
	} `json:"data"`
}

type ProductData struct {
	Id          int                      `json:"id" example:"1"`
	Name        string                   `json:"name" example:"Product 1"`
	Description string                   `json:"description" example:"Product 1 description"`
	Price       int                      `json:"price" example:"1000"`
	Status      string                   `json:"status" example:"available"`
	SKU         string                   `json:"sku" example:"SKU_1"`
	Attributes  models.ProductAttributes `json:"attributes"`
	Category    string                   `json:"category" example:"Category 1"`
	Brand       string                   `json:"brand" example:"brand 1"`
	Weight      float64                  `json:"weight" example:"10.5"`
	WeightUnit  string                   `json:"weight_unit" example:"kg"`
}

type PaginationMeta struct {
	Total      int64 `json:"total" example:"100"`
	Page       int   `json:"page" example:"2"`
	Limit      int   `json:"limit" example:"10"`
	TotalPages int   `json:"total_pages" example:"10"`
	NextPage   int   `json:"next_page" example:"3"`
	PrevPage   int   `json:"prev_page" example:"1"`
}

func SendListProductsResponse(gin *gin.Context, products []*models.Product, meta PaginationMeta) {
	r := &ListProductsResponse{}
	r.Message = "Products retrieved successfully"
	r.Data.Products = make([]ProductData, len(products))

	for i, product := range products {
		r.Data.Products[i].Id = int(product.ID)
		r.Data.Products[i].Name = product.Name
		r.Data.Products[i].Description = product.Description
		r.Data.Products[i].Price = int(product.Price)
		r.Data.Products[i].Status = string(product.Status)
		r.Data.Products[i].SKU = product.SKU
		r.Data.Products[i].Attributes = product.Attributes
		r.Data.Products[i].Category = product.Category
		r.Data.Products[i].Brand = product.Brand
		r.Data.Products[i].Weight = product.Weight
		r.Data.Products[i].WeightUnit = product.WeightUnit
	}

	response.Json(gin, r.Message, r.Data, http.StatusOK)
}

type GetProductResponse struct {
	Message string `json:"message" example:"Product details retrieved successfully"`
	Data    struct {
		Product ProductData `json:"product"`
	} `json:"data"`
}

func SendGetProductResponse(gin *gin.Context, product *models.Product) {
	r := &GetProductResponse{}
	r.Message = "Product details retrieved successfully"
	r.Data.Product.Id = int(product.ID)
	r.Data.Product.Name = product.Name
	r.Data.Product.Description = product.Description
	r.Data.Product.Price = int(product.Price)
	r.Data.Product.Status = string(product.Status)
	r.Data.Product.SKU = product.SKU
	r.Data.Product.Attributes = product.Attributes
	r.Data.Product.Category = product.Category
	r.Data.Product.Brand = product.Brand
	r.Data.Product.Weight = product.Weight
	r.Data.Product.WeightUnit = product.WeightUnit
	response.Json(gin, r.Message, r.Data, http.StatusOK)
}

type CreateProductResponse struct {
	Message string `json:"message" example:"Product created successfully"`
	Data    struct {
		Product ProductData `json:"product"`
	} `json:"data"`
}

func SendCreateProductResponse(gin *gin.Context, product *models.Product) {
	r := &CreateProductResponse{}
	r.Message = "Product created successfully"
	r.Data.Product.Id = int(product.ID)
	r.Data.Product.Name = product.Name
	r.Data.Product.Description = product.Description
	r.Data.Product.Price = int(product.Price)
	r.Data.Product.Status = string(product.Status)
	r.Data.Product.SKU = product.SKU
	r.Data.Product.Attributes = product.Attributes
	r.Data.Product.Category = product.Category
	r.Data.Product.Brand = product.Brand
	r.Data.Product.Weight = product.Weight
	r.Data.Product.WeightUnit = product.WeightUnit
	response.Json(gin, r.Message, r.Data, http.StatusCreated)
}

type UpdateProductResponse struct {
	Message string `json:"message" example:"Product updated successfully"`
	Data    struct {
		Product ProductData `json:"product"`
	} `json:"data"`
}

func SendUpdateProductResponse(gin *gin.Context, product *models.Product) {
	r := &UpdateProductResponse{}
	r.Message = "Product updated successfully"
	r.Data.Product.Id = int(product.ID)
	r.Data.Product.Name = product.Name
	r.Data.Product.Description = product.Description
	r.Data.Product.Price = int(product.Price)
	r.Data.Product.Status = string(product.Status)
	r.Data.Product.SKU = product.SKU
	r.Data.Product.Attributes = product.Attributes
	r.Data.Product.Category = product.Category
	r.Data.Product.Brand = product.Brand
	r.Data.Product.Weight = product.Weight
	r.Data.Product.WeightUnit = product.WeightUnit
	response.Json(gin, r.Message, r.Data, http.StatusOK)
}
