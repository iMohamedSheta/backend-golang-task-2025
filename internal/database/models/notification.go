package models

import (
	"taskgo/internal/enums"
	"time"
)

type Notification struct {
	Base
	UserID         uint                      `gorm:"index;not null" json:"user_id"`
	Type           enums.NotificationType    `gorm:"type:varchar(20);not null" json:"type"`
	Status         enums.NotificationStatus  `gorm:"type:varchar(20);not null;default:'pending'" json:"status"`
	Channel        enums.NotificationChannel `gorm:"type:varchar(20);not null" json:"channel"`
	Title          string                    `gorm:"size:200;not null" json:"title"`
	Message        string                    `gorm:"type:text;not null" json:"message"`
	Data           string                    `gorm:"type:jsonb" json:"data"` // Additional notification data (nullable)
	ReadAt         time.Time                 `json:"read_at,omitempty"`
	SentAt         time.Time                 `json:"sent_at,omitempty"`
	FailureReason  string                    `gorm:"type:text" json:"failure_reason,omitempty"`
	RetryCount     int                       `gorm:"default:0" json:"retry_count"`
	NextRetryAt    time.Time                 `json:"next_retry_at,omitempty"`
	NotifiableID   uint                      `gorm:"index" json:"notifiable_id"` // morph relation (user, product, order, etc) the notification is related to
	NotifiableType string                    `gorm:"size:50" json:"notifiable_type"`
	User           User                      `gorm:"foreignKey:UserID" json:"user"` // relationship to user
}
