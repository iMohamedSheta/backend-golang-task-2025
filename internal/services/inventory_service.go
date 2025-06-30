package services

import (
	"context"
	"taskgo/internal/database/models"
	"taskgo/internal/repository"
)

type InventoryService struct {
	inventoryRepository *repository.InventoryRepository
	productRepository   *repository.ProductRepository
}

func NewInventoryService(inventoryRepository *repository.InventoryRepository, productRepository *repository.ProductRepository) *InventoryService {
	return &InventoryService{
		inventoryRepository: inventoryRepository,
		productRepository:   productRepository,
	}
}

// Reserve inventory for one product
// func (s *InventoryService) ReserveInventory(ctx *gin.Context, product models.Product, item requests.OrderItemRequest) error {
// 	cache, err := pkgRedis.Default()
// 	if err != nil {
// 		logger.Channel("inventory_log").Error("Redis connection failed: " + err.Error())
// 		return pkgErrors.NewServerError("Internal Server Error", "Internal Server Error: Failed to connect to cache", err)
// 	}

// 	ctx := context.Background()

// 	auth

// 	inventoryKey, err := product.GetInventoryCacheKey()
// 	if err != nil {
// 		logger.Channel("inventory_log").Error(fmt.Sprintf("Failed to get inventory cache key for product %s: %v", product.Name, err))
// 		return pkgErrors.NewServerError("Internal Server Error", "Internal Server Error: Failed to get inventory cache key", err)
// 	}

// 	logger.Channel("inventory_log").Info("Checking inventory for product: " + product.Name)

// 	_, err = cacher.Remember(ctx, inventoryKey, cacher.RememberForever, func() (any, error) {
// 		if err := product.LoadInventory(); err != nil {
// 			logger.Channel("inventory_log").Error(fmt.Sprintf("Failed to load inventory for product %s: %v", product.Name, err))
// 			return "", pkgErrors.NewServerError("Internal Server Error", "Internal Server Error: Failed to load inventory", err)
// 		}
// 		return product.Inventory.Quantity, nil
// 	})

// 	if err != nil {
// 		logger.Channel("inventory_log").Error(fmt.Sprintf("Failed to remember inventory cache for %s: %v", product.Name, err))
// 		return pkgErrors.NewServerError("Internal Server Error", "Internal Server Error: Failed to check inventory", err)
// 	}

// 	newStock, err := cache.DecrBy(ctx, inventoryKey, int64(item.Quantity)).Result()
// 	if err != nil {
// 		logger.Channel("inventory_log").Error(fmt.Sprintf("Redis DECRBY failed for %s: %v", inventoryKey, err))
// 		return pkgErrors.NewServerError("Internal Server Error", "Internal Server Error: Failed to reserve inventory", err)
// 	}

// 	logger.Channel("inventory_log").Info(fmt.Sprintf("Decremented inventory for %s, new stock: %d", inventoryKey, newStock))

// 	if newStock < 0 {
// 		_, rollbackErr := cache.IncrBy(ctx, inventoryKey, int64(item.Quantity)).Result()
// 		if rollbackErr != nil {
// 			logger.Channel("inventory_log").Error(fmt.Sprintf("Failed to rollback inventory for %s: %v", inventoryKey, rollbackErr))
// 		}
// 		logger.Channel("inventory_log").Warn(fmt.Sprintf("Reservation failed, stock too low for product %s", product.Name))

// 		return pkgErrors.NewValidationError(map[string]any{
// 			"items": fmt.Sprintf("Insufficient stock for product %s. Stock exhausted during reservation", product.Name),
// 		})
// 	}

// 	return nil
// }

// [TODO]: Finish GetInventoryCacheKey implementation to finish the reservation of the order
// ReserveInventoriesAtomic atomically reserves inventory for multiple products
// the only way to do this i found is lua script and MULTI in redis but lua script is better for conditional reservation
// Lua scripts are executed as a single atomic operation in redis, ensuring that no other commands will run in the middle of its execution
func (s *InventoryService) ReserveInventoriesAtomic(ctx context.Context, order *models.Order, orderItems []models.OrderItem) error {
	// cache := deps.Cache().Redis
	// log := deps.Log().Channel("inventory_log")
	// if cache == nil {
	// 	return errors.New("InventoryService: ReserveInventoriesAtomic redis cache connection failed")
	// }

	// // Extract unique product IDs from order items
	// productIDs := make([]uint, 0, len(orderItems))
	// for _, item := range orderItems {
	// 	productIDs = append(productIDs, item.ProductID)
	// }
	// uniqueProductIDs := utils.UniqueSliceUInts(productIDs)

	// // Fetch products with inventory data
	// products, err := s.productRepository.FindByIDsWithInventory(uniqueProductIDs, "id", "product_id", "quantity")
	// if err != nil {
	// 	log.Error("Failed to fetch products: " + err.Error())
	// 	return err
	// }

	// productMap := make(map[uint]models.Product, len(products))
	// for _, p := range products {
	// 	productMap[p.ID] = p
	// }

	// // Prepare Redis KEYS and ARGS
	// var redisKeys []string
	// var redisArgs []any

	// for _, item := range orderItems {
	// 	product, exists := productMap[item.ProductID]
	// 	if !exists {
	// 		return pkgErrors.NewServerError("Product not found", "Product not found in DB", nil)
	// 	}

	// 	key, err := product.GetInventoryCacheKey()
	// 	if err != nil {
	// 		logger.Channel("inventory_log").Error(fmt.Sprintf("Cache key error for product %s: %v", product.Name, err))
	// 		return pkgErrors.NewServerError("Redis error", "Failed to generate inventory key", err)
	// 	}

	// 	_, err = cacher.Remember(ctx, key, cacher.RememberForever, func() (any, error) {
	// 		return product.Inventory.Quantity, nil
	// 	})

	// 	if err != nil {
	// 		logger.Channel("inventory_log").Error(fmt.Sprintf("Failed to set inventory cache for %s: %v", product.Name, err))
	// 		return pkgErrors.NewServerError("Redis error", "Failed to set inventory in cache", err)
	// 	}

	// 	redisKeys = append(redisKeys, key)
	// 	redisArgs = append(redisArgs, item.Quantity)
	// }

	// // Lua script: check and reserve inventory atomically
	// luaScript := `
	// 	for i = 1, #KEYS do
	// 		local current = tonumber(redis.call("GET", KEYS[i]))
	// 		local quantity = tonumber(ARGV[i])
	// 		if not current or current < quantity then
	// 			return 0
	// 		end
	// 	end
	// 	for i = 1, #KEYS do
	// 		redis.call("DECRBY", KEYS[i], ARGV[i])
	// 	end
	// 	return 1
	// `

	// result, err := cache.Eval(ctx, luaScript, redisKeys, redisArgs...).Int()
	// if err != nil {
	// 	logger.Channel("inventory_log").Error("Redis Lua Eval failed: " + err.Error())
	// 	return pkgErrors.NewServerError("Reservation failed", "Lua script error", err)
	// }
	// if result != 1 {
	// 	logger.Channel("inventory_log").Warn("Insufficient inventory for one or more products")
	// 	return pkgErrors.NewValidationError(map[string]any{
	// 		"inventory": "Insufficient stock for one or more products",
	// 	})
	// }

	// logger.Channel("inventory_log").Info(fmt.Sprintf("Reserved inventory successfully for order ID %d", order.ID))
	return nil
}

// RestoreInventory restores inventory if order process fails
// func (s *InventoryService) RestoreInventory(ctx context.Context, product models.Product, quantity int) error {
// 	cache := deps.Cache().Redis
// 	if cache == nil {
// 		logger.Channel("inventory_log").Error("Redis connection failed (restore): " + err.Error())
// 		return pkgErrors.NewServerError("Internal Server Error: Failed to connect to cache", "Internal Server Error: Failed to connect to cache", err)
// 	}

// 	inventoryKey, err := product.GetInventoryCacheKey()
// 	if err != nil {
// 		logger.Channel("inventory_log").Error(fmt.Sprintf("Failed to get inventory key for  %s: %v", product.Name, err))
// 		return pkgErrors.NewServerError("Internal Server Error", "Internal Server Error: Failed to get inventory key", err)
// 	}

// 	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
// 	defer cancel()

// 	_, err = cache.IncrBy(ctx, inventoryKey, int64(quantity)).Result()
// 	if err != nil {
// 		logger.Channel("inventory_log").Error(fmt.Sprintf("Failed to restore inventory for %s: %v", inventoryKey, err))
// 		return pkgErrors.NewServerError("Internal Server Error: Failed to restore inventory", "Internal Server Error: Failed to restore inventory", err)
// 	}

// 	logger.Channel("inventory_log").Info(fmt.Sprintf("Restored %d items to inventory for %s", quantity, product.Name))
// 	return nil
// }

// // SyncInventoryToDB syncs inventory quantity from redis to database this should be scheduled or asynce
// func (s *InventoryService) SyncInventoryToDB(ctx context.Context, inventory *models.Inventory) error {
// 	cache, err := pkgRedis.Default()
// 	if err != nil {
// 		logger.Channel("inventory_log").Error("Redis connection failed (sync): " + err.Error())
// 		return err
// 	}

// 	inventoryKey := inventory.GetInventoryCacheKey()

// 	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
// 	defer cancel()

// 	currentQuantity, err := cache.Get(ctx, inventoryKey).Result()
// 	if err != nil {
// 		logger.Channel("inventory_log").Error(fmt.Sprintf("Failed to get Redis key %s: %v", inventoryKey, err))
// 		return err
// 	}

// 	quantity, err := strconv.Atoi(currentQuantity)
// 	if err != nil {
// 		logger.Channel("inventory_log").Error(fmt.Sprintf("Failed to convert Redis value to int for %s: %v", inventoryKey, err))
// 		return err
// 	}

// 	logger.Channel("inventory_log").Info(fmt.Sprintf("Syncing inventory to DB for %s: quantity=%d", inventoryKey, quantity))
// 	err = s.inventoryRepository.UpdateQuantity(inventory.ID, quantity)
// 	if err != nil {
// 		logger.Channel("inventory_log").Error(fmt.Sprintf("Failed to update inventory for  %s: %v", inventoryKey, err))
// 		return err
// 	}

// 	return nil
// }
