package ioc_test

import (
	"sync"
	"testing"

	"taskgo/pkg/ioc"
)

type ServiceA struct {
	Value string
}

type ServiceB struct {
	Value int
}

type ShutdownMock struct {
	Called bool
	mu     sync.Mutex
}

func (s *ShutdownMock) Shutdown() error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.Called = true
	return nil
}

func TestSingletonResolution(t *testing.T) {
	c := ioc.NewContainer()

	err := ioc.Singleton(c, func(c *ioc.Container) (*ServiceA, error) {
		return &ServiceA{Value: "singleton"}, nil
	})
	if err != nil {
		t.Fatalf("failed to register singleton: %v", err)
	}

	c.Bootstrap()

	a1, _ := ioc.Resolve[*ServiceA](c)
	a2, _ := ioc.Resolve[*ServiceA](c)

	if a1 != a2 {
		t.Fatal("expected singleton instance")
	}
	if a1.Value != "singleton" {
		t.Fatal("unexpected value")
	}
}

func TestTransientResolution(t *testing.T) {
	c := ioc.NewContainer()

	err := ioc.Bind(c, func(c *ioc.Container) (*ServiceA, error) {
		return &ServiceA{Value: "transient"}, nil
	})
	if err != nil {
		t.Fatalf("failed to bind transient: %v", err)
	}

	a1, _ := ioc.Resolve[*ServiceA](c)
	a2, _ := ioc.Resolve[*ServiceA](c)

	if a1 == a2 {
		t.Fatal("expected different instances for transient")
	}
}

func TestShutdownIsCalled(t *testing.T) {
	c := ioc.NewContainer()

	mock := &ShutdownMock{}

	_ = ioc.Singleton(c, func(c *ioc.Container) (*ShutdownMock, error) {
		return mock, nil
	})

	c.Bootstrap()
	c.ShutdownAll()

	mock.mu.Lock()
	called := mock.Called
	mock.mu.Unlock()

	if !called {
		t.Fatal("expected shutdown to be called")
	}
}

func TestResolveNonExistent(t *testing.T) {
	c := ioc.NewContainer()

	_, err := ioc.Resolve[*ServiceB](c)
	if err == nil {
		t.Fatal("expected error when resolving unregistered service")
	}
}

func TestConcurrencySafety(t *testing.T) {
	c := ioc.NewContainer()

	err := ioc.Singleton(c, func(c *ioc.Container) (*ServiceA, error) {
		return &ServiceA{Value: "concurrent"}, nil
	})
	if err != nil {
		t.Fatalf("failed to register singleton: %v", err)
	}

	c.Bootstrap()

	var wg sync.WaitGroup
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			s, err := ioc.Resolve[*ServiceA](c)
			if err != nil || s.Value != "concurrent" {
				t.Errorf("concurrent resolve failed: %v", err)
			}
		}()
	}
	wg.Wait()
}
