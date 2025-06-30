package deps

import (
	"fmt"
	"taskgo/pkg/ioc"
	"taskgo/pkg/notify"
)

func Notify() *notify.Notify {
	n, err := ioc.AppMake[*notify.Notify]()
	if err != nil {
		Log().Log().Error(fmt.Sprintf("Notify can't be resolved: %s", err.Error()))
		return nil
	}
	return n
}
