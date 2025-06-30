package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"taskgo/internal/database/models"
	"taskgo/internal/deps"
	"taskgo/pkg/notify"
)

func DatabaseChannelHandler(ctx context.Context, task *notify.NotificationTask) error {
	db := deps.Gorm().DB

	// Marshal data map to JSON string
	dataBytes, err := json.Marshal(task.Data)
	if err != nil {
		return fmt.Errorf("failed to marshal data: %w", err)
	}

	notification := &models.Notification{
		Type:           task.NotificationType,
		Data:           string(dataBytes),
		NotifiableID:   task.NotifiableID,
		NotifiableType: task.NotifiableType,
		ReadAt:         nil,
	}

	if err := db.Create(notification).Error; err != nil {
		return fmt.Errorf("failed to create notification: %w", err)
	}
	return nil
}
