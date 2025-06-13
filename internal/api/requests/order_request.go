package requests

type OrderItemRequest struct {
	ProductId uint `json:"product_id" validate:"required,gt=0"`
	Quantity  int  `json:"quantity" validate:"required,gt=0"`
}

type CreateOrderRequest struct {
	UserId          uint               `json:"user_id" validate:"required,gt=0"`
	Items           []OrderItemRequest `json:"items" validate:"required,dive"`
	ShippingAddress string             `json:"shipping_address" validate:"required,min=5,max=200"`
	BillingAddress  string             `json:"billing_address" validate:"required,min=5,max=200"`
	PaymentMethod   string             `json:"payment_method" validate:"required,oneof=credit_card paypal bank_transfer cash_on_delivery"`
	Notes           string             `json:"notes,omitempty" validate:"omitempty,max=500"`
	Request
}

func (r *CreateOrderRequest) Messages() map[string]string {
	return map[string]string{
		"items.required":            "At least one item is required",
		"items.dive":                "At least one item is required",
		"items.product_id.required": "Product ID is required",
		"items.product_id.gt":       "Product ID must be greater than 0",
		"items.quantity.required":   "Quantity is required",
		"items.quantity.gt":         "Quantity must be greater than 0",
		"user_id.required":          "User ID is required",
		"user_id.gt":                "User ID must be greater than 0",
		"shipping_address.required": "Shipping address is required",
		"shipping_address.min":      "Shipping address must be at least 5 characters",
		"shipping_address.max":      "Shipping address must be at most 200 characters",
		"billing_address.required":  "Billing address is required",
		"billing_address.min":       "Billing address must be at least 5 characters",
		"billing_address.max":       "Billing address must be at most 200 characters",
		"payment_method.required":   "Payment method is required",
		"payment_method.oneof":      "Payment method must be one of credit_card, paypal, bank_transfer, cash_on_delivery",
	}
}
