package models

import (
	"time"
)

type Notification struct {
	Base
	Type           string     `gorm:"type:varchar(255);not null" json:"type"`
	Data           string     `gorm:"type:jsonb" json:"data"` // Additional notification data (nullable)
	ReadAt         *time.Time `json:"read_at,omitempty"`
	NotifiableID   uint       `gorm:"index" json:"notifiable_id"` // morph relation (user, product, order, etc) the notification is related to
	NotifiableType string     `gorm:"size:50" json:"notifiable_type"`
}
