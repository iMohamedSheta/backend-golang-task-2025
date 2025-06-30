package services

import (
	"context"
	"fmt"
	"taskgo/internal/api/requests"
	"taskgo/internal/database/models"
	"taskgo/internal/repository"
	pkgErrors "taskgo/pkg/errors"
)

type OrderService struct {
	inventoryService  *InventoryService
	orderRepository   *repository.OrderRepository
	productRepository *repository.ProductRepository
}

func NewOrderService(inventoryService *InventoryService, orderRepo *repository.OrderRepository, productRepo *repository.ProductRepository) *OrderService {
	return &OrderService{
		inventoryService:  inventoryService,
		orderRepository:   orderRepo,
		productRepository: productRepo,
	}
}

// TODO: Implement full order service from first try
// -------------------------------------------------------------------------------------------------------------------
// Check if products is in inventory and status available
// ummm... we need to update inventory quantity better to save it inside cache server (redis) i am not going to write it now
// we will just think it is done
// Also we need to lock the inventory quantity while we are updating it
// Future implementation:
// 1. Try to reserve inventory in Redis with expiration
// 2. Create order in DB
// 3. Confirm reservation in Redis
// 4. Create orderItems in DB
// 4. Update actual inventory in background job
// DECR/DECRBY and INCR/INCRBY in redis is the way
// ops i need to sync database now with redis (better than locking db while updating the inventory)
// --------------------------------------------------------------------------------------------------------------------
// Ok that's a problem reserve inventory should be for all the products in one transaction not for each product
// err := s.inventoryService.ReserveInventory(product, item)
func (s *OrderService) CreateOrder(ctx context.Context, req *requests.CreateOrderRequest) (*models.Order, error) {
	// Create base order
	order := &models.Order{
		UserID:          req.UserId,
		ShippingAddress: req.ShippingAddress,
		BillingAddress:  req.BillingAddress,
		Notes:           req.Notes,
	}

	// Extract order productsIDs so we can check them
	var productIDs []uint
	for _, item := range req.Items {
		productIDs = append(productIDs, item.ProductId)
	}

	// Fetch products with inventory needed columns
	products, err := s.productRepository.FindByIDsWithInventory(productIDs, "id", "product_id", "quantity")
	if err != nil {
		return nil, pkgErrors.NewServerError("Internal Server Error: Failed to fetch products", "Internal Server Error: Failed to fetch products", err)
	}

	// Validation for product
	productMap := make(map[uint]models.Product)
	for _, product := range products {
		productMap[product.ID] = product
	}
	orderItems := make([]*models.OrderItem, len(req.Items))

	// Order checks before order creation
	for i, item := range req.Items {

		// Validate products exists
		product, ok := productMap[item.ProductId]
		if !ok {
			return nil, pkgErrors.NewValidationError(map[string]any{
				"items": fmt.Sprintf("Product with ID  %d does not exist", item.ProductId),
			})
		}

		// Validate product quantity is available (no reservation yet)
		if product.Inventories[0].GetTotalQuantity() < item.Quantity {
			return nil, pkgErrors.NewValidationError(map[string]any{
				"items": fmt.Sprintf("Product with ID %d has insufficient stock", item.ProductId),
			})
		}

		orderItems[i] = mapOrderItemData(&item, &product)
	}

	// Create order with order items after reserving inventory
	err = s.orderRepository.CreateWithOrderItems(order, orderItems)
	if err != nil {
		return nil, pkgErrors.NewServerError("Internal Server Error: Failed to create order", "Internal Server Error: Failed to create order", err)
	}

	// Load order with order items so it can be returned in response
	order.OrderItems = make([]models.OrderItem, len(orderItems))
	for i, item := range orderItems {
		order.OrderItems[i] = *item
	}

	return order, nil
}

func mapOrderItemData(item *requests.OrderItemRequest, product *models.Product) *models.OrderItem {
	orderItem := &models.OrderItem{
		ProductID: item.ProductId,
		Quantity:  item.Quantity,
	}
	orderItem.CalculateUnitPrice(product.Price)
	orderItem.CalculateTotalPrice()
	return orderItem
}
