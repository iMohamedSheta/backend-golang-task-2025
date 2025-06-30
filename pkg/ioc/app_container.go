package ioc

import "sync"

var (
	appContainer *Container
	once         sync.Once
)

// SetAppContainer sets the global application container (only once).
func SetAppContainer(c *Container) {
	once.Do(func() {
		appContainer = c
	})
}

// AppContainer returns the global application container.
func AppContainer() *Container {
	if appContainer == nil {
		panic("AppContainer is not initialized. Did you forget to call SetAppContainer?")
	}
	return appContainer
}

// AppMake is a shortcut to resolve a service from the global container.
func AppMake[T any]() (T, error) {
	return Make[T](AppContainer())
}

type Shutdownable interface {
	Shutdown() error
}
