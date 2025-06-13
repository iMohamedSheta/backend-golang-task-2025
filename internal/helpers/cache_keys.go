package helpers

import "fmt"

func GetInventoryCacheKey(InventoryID uint) string {
	return fmt.Sprintf("inventory:quantity:%d", InventoryID)
}
