package ioc

import (
	"fmt"
	"reflect"
	"sync"
)

/*
	This package is meant to be used as service resolver using service container.
*/

// Container is the main services container.
type Container struct {
	mu            sync.RWMutex
	services      map[string]any // the services container should be associated with the service name
	bootstrappers []func()
	shutdownables []Shutdownable
}

type serviceProvider[T any] func(c *Container) (T, error)

type lifetime int

const (
	singleton lifetime = iota
	transient
)

// Service
type serviceContainer[T any] struct {
	mu       sync.Mutex
	instance T                  // the resolved service if it is resolved
	provider serviceProvider[T] // the service provider function that returns the service after it is resolved
	isMade   bool               // indicates whether the service is made
	life     lifetime           // the service lifetime (singleton or transient per request)
}

func (c *Container) GetRegisteredService() map[string]any {
	return c.services
}

func (s *serviceContainer[T]) Resolve(c *Container) (T, error) {
	// Always create a new instance for transient
	if s.life == transient {
		return s.provider(c)
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	// Otherwise, it's a singleton
	if s.isMade {
		return s.instance, nil
	}

	instance, err := s.provider(c)
	if err != nil {
		return instance, err
	}

	s.instance = instance
	s.isMade = true

	return instance, nil
}

// NewContainer creates a new services container.
func NewContainer() *Container {
	return &Container{
		services: make(map[string]any),
	}
}

// Bind registers a transient service, resolved fresh on each Resolve() call.
func Bind[T any](c *Container, provider serviceProvider[T]) error {
	if err := RegisterServiceProvider(c, provider, transient, false); err != nil {
		return err
	}
	return nil
}

// Singleton registers a singleton service, resolved eagerly at Bootstrap().
func Singleton[T any](c *Container, provider serviceProvider[T]) error {
	if err := RegisterServiceProvider(c, provider, singleton, true); err != nil {
		return err
	}
	return nil
}

// SingletonDeferred registers a singleton, but defers instantiation until used.
func SingletonDeferred[T any](c *Container, provider serviceProvider[T]) error {
	if err := RegisterServiceProvider(c, provider, singleton, false); err != nil {
		return err
	}
	return nil
}

func RegisterServiceProvider[T any](c *Container, provider serviceProvider[T], life lifetime, bootstrap bool) error {
	name, err := getServiceName[T]()
	if err != nil {
		return fmt.Errorf("register failed: %w", err)
	}

	c.mu.Lock()
	defer c.mu.Unlock()

	c.services[name] = &serviceContainer[T]{
		provider: provider,
		life:     life,
	}

	if bootstrap && life == singleton {
		c.bootstrappers = append(c.bootstrappers, func() {

			instance, err := Resolve[T](c)
			if err != nil {
				fmt.Printf("failed to bootstrap service[%s]: %s\n", reflect.TypeOf(instance), err)
				return
			}

			if s, ok := any(instance).(Shutdownable); ok {
				c.mu.Lock()
				c.shutdownables = append(c.shutdownables, s)
				c.mu.Unlock()
			}

		})
	}

	return nil
}

// Resolve tries to resolve a service. If the service is not found, it returns an error.
func Resolve[T any](c *Container) (T, error) {
	var zero T
	name, err := getServiceName[T]()
	if err != nil {
		return zero, err
	}

	c.mu.RLock()
	any, doesServiceExists := c.services[name]
	c.mu.RUnlock()

	if !doesServiceExists {
		return zero, fmt.Errorf("failed to resolve service[%s]: not found service is not registered in the container", name)
	}

	serviceContainer, ok := any.(*serviceContainer[T])
	if !ok {
		return zero, fmt.Errorf("failed to resolve service[%s]: type assertion failed  - the service container for  %s is not a service container", name, name)
	}

	return serviceContainer.Resolve(c)
}

// Alias for Resolve  - make a service if it is not made yet and return an error if it is not found.
func Make[T any](c *Container) (T, error) {
	return Resolve[T](c)
}

// Bootstrap initializes all eagerly registered services.
func (c *Container) Bootstrap() {
	for _, bootstrapper := range c.bootstrappers {
		bootstrapper()
	}

	c.bootstrappers = nil // Clear after bootstrapping to prevent bootstrapping again
}

func (c *Container) ShutdownAll() {
	c.mu.RLock()
	shutdownables := c.shutdownables
	c.mu.RUnlock()

	for _, s := range shutdownables {
		if err := s.Shutdown(); err != nil {
			fmt.Printf("Error shutting down service: %v\n", err)
		}
	}
}

func getServiceName[T any]() (string, error) {
	typeForT := reflect.TypeOf((*T)(nil)).Elem()

	if typeForT.Kind() != reflect.Pointer {
		return "", fmt.Errorf("Resolve (Make[*T]()): type must be a pointer, got %v", typeForT)
	}

	elem := typeForT.Elem()

	if elem.PkgPath() == "" || elem.Name() == "" {
		return "", fmt.Errorf("unsupported anonymous or builtin type: %v", typeForT)
	}

	return "*" + elem.PkgPath() + "." + elem.Name(), nil
}
