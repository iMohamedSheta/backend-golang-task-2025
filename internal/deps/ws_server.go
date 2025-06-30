package deps

import (
	"fmt"
	"taskgo/pkg/ioc"
	"taskgo/pkg/ws"
)

func WS() *ws.Server {
	ws, err := ioc.AppMake[*ws.Server]()
	if err != nil {
		Log().Log().Error(fmt.Sprintf("WS Dependency Container Error: %s", err.Error()))
		return nil
	}

	return ws
}
