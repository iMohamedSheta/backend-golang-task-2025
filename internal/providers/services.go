package providers

import (
	"taskgo/internal/deps"
	"taskgo/internal/repository"
	"taskgo/internal/services"
	"taskgo/pkg/ioc"
)

// Register Services in Container
func RegisterServices(c *ioc.Container) {
	// Register Product Service
	err := ioc.Bind(c, func(c *ioc.Container) (*services.ProductService, error) {
		pRepo, err := ioc.Make[*repository.ProductRepository](c)
		if err != nil {
			return nil, err
		}

		return services.NewProductService(pRepo), nil
	})
	logBindErr("ProductService", err)

	// Register User Service
	err = ioc.Bind(c, func(c *ioc.Container) (*services.UserService, error) {
		uRepo, err := ioc.Make[*repository.UserRepository](c)
		if err != nil {
			return nil, err
		}

		return services.NewUserService(uRepo), nil
	})
	logBindErr("UserService", err)

	// Register Auth Service
	err = ioc.Bind(c, func(c *ioc.Container) (*services.AuthService, error) {
		uRepo, err := ioc.Make[*repository.UserRepository](c)
		if err != nil {
			return nil, err
		}

		jwtService, err := ioc.Make[*services.JwtService](c)
		if err != nil {
			return nil, err
		}

		return services.NewAuthService(uRepo, jwtService), nil
	})
	logBindErr("AuthService", err)

	// Register Inventory Service
	err = ioc.Bind(c, func(c *ioc.Container) (*services.InventoryService, error) {
		invRepo, err := ioc.Make[*repository.InventoryRepository](c)
		if err != nil {
			return nil, err
		}
		productRepo, err := ioc.Make[*repository.ProductRepository](c)
		if err != nil {
			return nil, err
		}

		return services.NewInventoryService(invRepo, productRepo), nil
	})
	logBindErr("InventoryService", err)

	// Register Order Service
	err = ioc.Bind(c, func(c *ioc.Container) (*services.OrderService, error) {
		invService, err := ioc.Make[*services.InventoryService](c)
		if err != nil {
			return nil, err
		}
		orderRepo, err := ioc.Make[*repository.OrderRepository](c)
		if err != nil {
			return nil, err
		}

		productRepo, err := ioc.Make[*repository.ProductRepository](c)
		if err != nil {
			return nil, err
		}

		return services.NewOrderService(invService, orderRepo, productRepo), nil
	})
	logBindErr("OrderService", err)

	// Register Product Service
	err = ioc.Bind(c, func(c *ioc.Container) (*services.JwtService, error) {
		return services.NewJwtService(deps.Config()), nil
	})
	logBindErr("JwtService", err)
}
