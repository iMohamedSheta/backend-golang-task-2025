package enums

// Maximum length of a product Status = 20 characters
type ProductStatus string

const (
	ProductStatusAvailable   ProductStatus = "available"
	ProductStatusUnavailable ProductStatus = "unavailable"
	ProductStatusArchived    ProductStatus = "archived"
	ProductStatusDeleted     ProductStatus = "deleted"
)

func IsValidProductStatus(s string) bool {
	switch ProductStatus(s) {
	case ProductStatusAvailable, ProductStatusUnavailable, ProductStatusArchived, ProductStatusDeleted:
		return true
	default:
		return false
	}
}
