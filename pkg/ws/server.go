package ws

import (
	"github.com/gorilla/websocket"
)

type Server struct {
	Hub      *Hub
	Upgrader websocket.Upgrader
}

func NewWsServer(hub *Hub, upgrader websocket.Upgrader) *Server {
	return &Server{
		Hub:      hub,
		Upgrader: upgrader,
	}
}
