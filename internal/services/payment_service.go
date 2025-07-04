package services

import (
	"context"
	"fmt"
	"time"
)

type PaymentService struct{}

func NewPaymentService() *PaymentService {
	return &PaymentService{}
}

// Example to process payment logic
func (s *PaymentService) ProcessPayment(ctx context.Context, req any) error {
	fmt.Printf("🔄 Processing payment for request: %v\n", req)
	time.Sleep(5 * time.Second)
	fmt.Println("Payment processed successfully")
	return nil
}
