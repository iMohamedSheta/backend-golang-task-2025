package providers

import (
	"taskgo/internal/api/handlers"
	"taskgo/internal/policies"
	"taskgo/internal/services"
	"taskgo/pkg/ioc"
	"taskgo/pkg/utils"
	"taskgo/pkg/ws"
)

// Register Handlers in Container
func RegisterHandlers(c *ioc.Container) {
	// Register Product Handler
	err := ioc.Bind(c, func(c *ioc.Container) (*handlers.ProductHandler, error) {
		pService, err := ioc.Make[*services.ProductService](c)
		if err != nil {
			return nil, err
		}
		return handlers.NewProductHandler(
			pService,
			&policies.ProductPolicy{},
		), nil
	})
	logBindErr("ProductHandler", err)

	// Register Order handler
	err = ioc.Bind(c, func(c *ioc.Container) (*handlers.OrderHandler, error) {
		orderService, err := ioc.Make[*services.OrderService](c)
		if err != nil {
			return nil, err
		}
		return handlers.NewOrderHandler(
			orderService,
			&policies.OrderPolicy{},
		), nil
	})
	logBindErr("OrderHandler", err)

	// Register Auth handler
	err = ioc.Bind(c, func(c *ioc.Container) (*handlers.AuthHandler, error) {
		authService, err := ioc.Make[*services.AuthService](c)
		if err != nil {
			return nil, err
		}
		return handlers.NewAuthHandler(
			authService,
		), nil
	})
	logBindErr("AuthHandler", err)

	// Register User handler
	err = ioc.Bind(c, func(c *ioc.Container) (*handlers.UserHandler, error) {
		userService, err := ioc.Make[*services.UserService](c)
		if err != nil {
			return nil, err
		}
		return handlers.NewUserHandler(
			userService,
			&policies.UserPolicy{},
		), nil
	})
	logBindErr("UserHandler", err)

	// Register Notification handler
	err = ioc.Bind(c, func(c *ioc.Container) (*handlers.NotificationHandler, error) {
		wsServer, err := ioc.Make[*ws.Server](c)
		if err != nil {
			return nil, err
		}
		return handlers.NewNotificationHandler(
			wsServer,
		), nil
	})
	logBindErr("NotificationHandler", err)
}

func logBindErr(module string, err error) {
	if err != nil {
		utils.PrintErr("Failed to load " + module + " module in the ioc container : " + err.Error())
	}
}
