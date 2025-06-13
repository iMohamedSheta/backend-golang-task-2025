package filters

import "time"

// ProductFilters struct for product filtering options
type ProductFilters struct {
	// Search filters
	Name   string `json:"name,omitempty" form:"name"`
	SKU    string `json:"sku,omitempty" form:"sku"`
	Brand  string `json:"brand,omitempty" form:"brand"`
	Search string `json:"search,omitempty" form:"search"`

	// Price filters
	MinPrice *float64 `json:"min_price,omitempty" form:"min_price"`
	MaxPrice *float64 `json:"max_price,omitempty" form:"max_price"`

	Status *string `json:"status,omitempty" form:"status"`

	CreatedAfter  *time.Time `json:"created_after,omitempty" form:"created_after"`
	CreatedBefore *time.Time `json:"created_before,omitempty" form:"created_before"`

	Category string `json:"category,omitempty" form:"category"`

	// Sorting
	SortBy    string `json:"sort_by,omitempty" form:"sort_by"`
	SortOrder string `json:"sort_order,omitempty" form:"sort_order" `

	// Pagination
	Page    int `json:"page,omitempty" form:"page"`
	PerPage int `json:"per_page,omitempty" form:"per_page"`
}

func (p *ProductFilters) GetSearchFields() []string {
	return []string{
		"name",
		"brand",
		"sku",
		"description",
	}
}

func (p *ProductFilters) GetSortFields() []string {
	return []string{
		"name",
		"sku",
		"brand",
		"price",
		"created_at",
		"updated_at",
	}
}
