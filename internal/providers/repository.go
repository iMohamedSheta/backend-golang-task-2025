package providers

import (
	"taskgo/internal/deps"
	"taskgo/internal/repository"
	"taskgo/pkg/ioc"
)

// Register Repository in the container
func RegisterRepository(c *ioc.Container) {
	// Register Product Repository
	err := ioc.Bind(c, func(c *ioc.Container) (*repository.ProductRepository, error) {
		gormDB, err := ioc.Make[*deps.GormDB](c)
		if err != nil {
			return nil, err
		}
		return repository.NewProductRepository(
			gormDB,
		), nil
	})
	logBindErr("ProductRepository", err)

	// Register User Repository
	err = ioc.Bind(c, func(c *ioc.Container) (*repository.UserRepository, error) {
		gormDB, err := ioc.Make[*deps.GormDB](c)
		if err != nil {
			return nil, err
		}
		return repository.NewUserRepository(
			gormDB,
		), nil
	})
	logBindErr("UserRepository", err)

	// Register Order Repository
	err = ioc.Bind(c, func(c *ioc.Container) (*repository.OrderRepository, error) {
		gormDB, err := ioc.Make[*deps.GormDB](c)
		if err != nil {
			return nil, err
		}
		return repository.NewOrderRepository(
			gormDB,
		), nil
	})
	logBindErr("OrderRepository", err)

	// Register Inventory Repository
	err = ioc.Bind(c, func(c *ioc.Container) (*repository.InventoryRepository, error) {
		gormDB, err := ioc.Make[*deps.GormDB](c)
		if err != nil {
			return nil, err
		}
		return repository.NewInventoryRepository(
			gormDB,
		), nil
	})
	logBindErr("InventoryRepository", err)
}
