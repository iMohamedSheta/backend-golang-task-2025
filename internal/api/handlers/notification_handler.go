package handlers

import (
	"strings"
	"taskgo/internal/deps"
	"taskgo/internal/helpers"
	"taskgo/internal/notification"
	"taskgo/pkg/errors"
	"taskgo/pkg/response"
	"taskgo/pkg/ws"

	"github.com/gin-gonic/gin"
)

type NotificationHandler struct {
	ws *ws.Server
}

// NewNotificationHandler creates a new notification handler
func NewNotificationHandler(wsServer *ws.Server) *NotificationHandler {
	return &NotificationHandler{
		ws: wsServer,
	}
}

func (h *NotificationHandler) WsUserNotificationHandler(c *gin.Context) {
	userId, authErr := helpers.GetAuthId(c)
	if authErr != nil {
		deps.Log().Log().Error(authErr.PrivateMessage + ": " + authErr.Error())
		response.UnauthorizedJson(c, authErr)
		return
	}

	// Echo back the subprotocols used in the request to the client (important)!!
	upgrader := h.ws.Upgrader
	subprotocolHeader := c.GetHeader("Sec-WebSocket-Protocol")
	protocols := strings.Split(subprotocolHeader, ",")
	for i := range protocols {
		protocols[i] = strings.TrimSpace(protocols[i])
	}
	upgrader.Subprotocols = protocols

	// Upgrade to websocket connection
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		deps.Log().Log().Error(err.Error())
		response.BadRequestErrorJson(c, errors.NewBadRequestError("", "Websocket upgrade failed to websocket connection", err))
		return
	}

	// Create new websocket client
	client := &ws.Client{
		Conn:   conn,
		UserID: userId,
		Send:   make(chan []byte, 256),
		Subs:   make(map[string]bool),
	}

	client.Listen(h.ws.Hub)
}

func (h *NotificationHandler) TestSendNotification(c *gin.Context) {
	// msg := c.Query("msg")
	// channel := c.Query("chan")
	// wsMsg := ws.WSMessage{
	// 	Channel: channel,
	// 	Data:    msg,
	// 	From:    "Server",
	// 	Type:    "notification",
	// }
	// h.ws.Hub.Broadcast(&wsMsg)
	authUser, err := helpers.GetAuthUser(c)
	if err != nil {
		response.UnauthorizedJson(c, errors.NewUnAuthorizedError("", "Unauthorized", err))
		return
	}

	err = deps.Notify().Send(notification.NewOrderCreatedNotification(123), authUser)
	if err != nil {
		response.BadRequestErrorJson(c, errors.NewBadRequestError("", "Notification failed to send", err))
		return
	}
	response.Json(c, "Notification sent successfully", nil, 200)
}
