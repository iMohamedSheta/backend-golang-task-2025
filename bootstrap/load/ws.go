package load

import (
	"net/http"
	"taskgo/internal/deps"
	"taskgo/pkg/ioc"
	"taskgo/pkg/utils"
	"taskgo/pkg/ws"

	"github.com/gorilla/websocket"
)

func InitWebsocketServer(c *ioc.Container) {
	err := ioc.Singleton(c, func(c *ioc.Container) (*ws.Server, error) {
		originRaw := deps.Config().GetWithDefault("cors.origin", nil)
		var origin []string
		if originRaw != nil {
			if val, ok := originRaw.([]string); ok {
				origin = val
			}
		}

		upgrader := websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool { return checkOrigin(r, origin) },
			// HandshakeTimeout: 1 * time.Second,
		}

		return ws.NewWsServer(ws.NewHub(), upgrader), nil
	})

	if err != nil {
		utils.PrintErr("Failed to load websocket server module in the ioc container : " + err.Error())
	}
}

func checkOrigin(r *http.Request, origin []string) bool {
	if len(origin) > 0 {
		if origin[0] == "*" {
			return true
		}

		for _, o := range origin {
			if r.Header.Get("Origin") == o {
				return true
			}
		}
	}

	return false
}
