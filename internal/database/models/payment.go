package models

import (
	"taskgo/internal/enums"
	"time"
)

type Payment struct {
	Base
	OrderID        uint                `gorm:"uniqueIndex;not null" json:"order_id"`
	Amount         float64             `gorm:"type:decimal(10,2);not null" json:"amount"`
	Currency       string              `gorm:"size:3;not null;default:'EGP'" json:"currency"`
	Status         enums.PaymentStatus `gorm:"type:varchar(20);not null;default:'pending'" json:"status"`
	Method         enums.PaymentMethod `gorm:"type:varchar(20);not null" json:"method"`
	TransactionID  string              `gorm:"size:100" json:"transaction_id"`
	PaymentDate    time.Time           `json:"payment_date"`
	FailureReason  string              `gorm:"type:text" json:"failure_reason,omitempty"`
	RefundAmount   float64             `gorm:"type:decimal(10,2);default:0" json:"refund_amount"`
	RefundDate     time.Time           `json:"refund_date,omitempty"`
	RefundReason   string              `gorm:"type:text" json:"refund_reason,omitempty"`
	PaymentDetails string              `gorm:"type:jsonb" json:"payment_details"`
}
