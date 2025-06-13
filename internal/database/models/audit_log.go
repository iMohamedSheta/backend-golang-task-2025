package models

import "taskgo/internal/enums"

type AuditLog struct {
	Base
	UserID      uint                 `gorm:"index" json:"user_id"` // foreign key user id
	Action      enums.AuditLogAction `gorm:"type:varchar(20);not null" json:"action"`
	ModelType   string               `gorm:"size:50;not null" json:"model_type"` // Morph relation
	ModelID     uint                 `gorm:"index;not null" json:"model_id"`
	OldValues   string               `gorm:"type:jsonb" json:"old_values,omitempty"`
	NewValues   string               `gorm:"type:jsonb" json:"new_values,omitempty"`
	IPAddress   string               `gorm:"size:45" json:"ip_address"`
	UserAgent   string               `gorm:"type:text" json:"user_agent"`
	RequestID   string               `gorm:"size:100" json:"request_id"`
	Description string               `gorm:"type:text" json:"description"`
	Metadata    string               `gorm:"type:jsonb" json:"metadata"`
	User        User                 `gorm:"foreignKey:UserID" json:"user"` // relationship to user (which made the action)
}
