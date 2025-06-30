package deps

import (
	"fmt"
	"taskgo/pkg/ioc"
)

func App[T any]() T {
	service, err := ioc.AppMake[T]()
	if err != nil {
		Log().Log().Error(fmt.Sprintf("AppMake[T] error: %v", err))
	}
	return service
}
