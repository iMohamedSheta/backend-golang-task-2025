package requests

import "taskgo/internal/enums"

type CreateProductRequest struct {
	Name        string              `json:"name" validate:"required,min=3,max=100"`
	Description string              `json:"description" validate:"required,min=10,max=500"`
	Price       float64             `json:"price" validate:"required,gt=0,lte=999999.99"`
	Status      enums.ProductStatus `json:"status" validate:"required"`
	Category    string              `json:"category" validate:"required,min=2,max=50"`
	Brand       string              `json:"brand" validate:"required,min=2,max=50"`
	Weight      float64             `json:"weight,omitempty" validate:"omitempty,gt=0,lte=10000"`
	WeightUnit  string              `json:"weight_unit,omitempty" validate:"required_with=Weight,oneof=kg g lb"`
	Attributes  map[string]any      `json:"attributes,omitempty" validate:"omitempty,dive,keys,required,endkeys,required"` //  Keys and Values are required (dive,keys,required,endkeys,required)
	Request
}

func (r *CreateProductRequest) Messages() map[string]string {
	return map[string]string{
		"name.required":             "Name is required",
		"name.min":                  "Name must be at least 3 characters",
		"name.max":                  "Name must be at most 100 characters",
		"description.required":      "Description is required",
		"description.min":           "Description must be at least 10 characters",
		"description.max":           "Description must be at most 500 characters",
		"price.required":            "Price is required",
		"price.gt":                  "Price must be greater than 0",
		"price.lte":                 "Price must be less than or equal to 999999.99",
		"status.required":           "Status is required",
		"category.required":         "Category is required",
		"category.min":              "Category must be at least 2 characters",
		"category.max":              "Category must be at most 50 characters",
		"brand.required":            "Brand is required",
		"brand.min":                 "Brand must be at least 2 characters",
		"brand.max":                 "Brand must be at most 50 characters",
		"weight.gt":                 "Weight must be greater than 0",
		"weight.lte":                "Weight must be less than or equal to 10000",
		"weight_unit.required_with": "WeightUnit is required when Weight is present",
		"weight_unit.oneof":         "WeightUnit must be one of kg, g, lb",
		"attributes.dive":           "Attributes must be a map",
		"attributes.keys":           "Attributes must have keys",
		"attributes.required":       "Attributes must have values",
		"attributes.endkeys":        "Attributes must have values",
	}
}

type UpdateProductRequest struct {
	Name        string              `json:"name" validate:"required,min=3,max=100"`
	Description string              `json:"description" validate:"required,min=10,max=500"`
	Price       float64             `json:"price" validate:"required,gt=0,lte=999999.99"`
	Status      enums.ProductStatus `json:"status" validate:"required"`
	Category    string              `json:"category" validate:"required,min=2,max=50"`
	Brand       string              `json:"brand" validate:"required,min=2,max=50"`
	Weight      float64             `json:"weight,omitempty" validate:"omitempty,gt=0,lte=10000"`
	WeightUnit  string              `json:"weight_unit,omitempty" validate:"required_with=Weight,oneof=kg g lb"`
	Attributes  map[string]any      `json:"attributes,omitempty" validate:"omitempty,dive,keys,required,endkeys,required"` //  Keys and Values are required (dive,keys,required,endkeys,required)
	Request
}

func (r *UpdateProductRequest) Messages() map[string]string {
	return map[string]string{
		"name.required":             "Name is required",
		"name.min":                  "Name must be at least 3 characters",
		"name.max":                  "Name must be at most 100 characters",
		"description.required":      "Description is required",
		"description.min":           "Description must be at least 10 characters",
		"description.max":           "Description must be at most 500 characters",
		"price.required":            "Price is required",
		"price.gt":                  "Price must be greater than 0",
		"price.lte":                 "Price must be less than or equal to 999999.99",
		"status.required":           "Status is required",
		"category.required":         "Category is required",
		"category.min":              "Category must be at least 2 characters",
		"category.max":              "Category must be at most 50 characters",
		"brand.required":            "Brand is required",
		"brand.min":                 "Brand must be at least 2 characters",
		"brand.max":                 "Brand must be at most 50 characters",
		"weight.gt":                 "Weight must be greater than 0",
		"weight.lte":                "Weight must be less than or equal to 10000",
		"weight_unit.required_with": "WeightUnit is required when Weight is present",
		"weight_unit.oneof":         "WeightUnit must be one of kg, g, lb",
		"attributes.dive":           "Attributes must be a map",
		"attributes.keys":           "Attributes must have keys",
		"attributes.required":       "Attributes must have values",
		"attributes.endkeys":        "Attributes must have values",
	}
}
