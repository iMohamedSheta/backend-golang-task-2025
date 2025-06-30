package providers

import (
	"taskgo/internal/deps"
	"taskgo/internal/repository"
	"taskgo/internal/services"
	"taskgo/internal/tasks"
	"taskgo/pkg/ioc"
	"taskgo/pkg/notify"
)

func RegisterTaskHandlers(c *ioc.Container) {
	// Register inventory check task handler
	err := ioc.Bind(c, func(c *ioc.Container) (*tasks.InventoryCheckHandler, error) {
		inventoryService, err := ioc.Make[*services.InventoryService](c)
		if err != nil {
			return nil, err
		}

		orderRepo, err := ioc.Make[*repository.OrderRepository](c)
		if err != nil {
			return nil, err
		}

		return tasks.NewInventoryCheckHandler(
			inventoryService,
			orderRepo,
		), nil
	})
	logBindErr("InventoryCheckHandler", err)

	//  Register ProcessPayment task handler
	err = ioc.Bind(c, func(c *ioc.Container) (*tasks.ProcessPaymentHandler, error) {
		return tasks.NewProcessPaymentHandler(), nil
	})
	logBindErr("ProcessPaymentHandler", err)

	// Register SendNotification task handler
	err = ioc.Bind(c, func(c *ioc.Container) (*notify.NotificationHandler, error) {
		return notify.NewNotificationHandler(
			deps.Notify(),
		), nil
	})
	logBindErr("notify.NotificationHandler", err)

}
