package ws

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"sync"

	"github.com/gorilla/websocket"
)

type Client struct {
	Conn     *websocket.Conn
	Send     chan []byte
	UserID   string
	Subs     map[string]bool // subscription to channels like "orders", "dashboard"
	subMutex sync.Mutex
}

type ChannelPolicy struct {
	Pattern  string
	CanRead  func(userID, channel string) bool
	CanWrite func(userID, channel string) bool
}

type Hub struct {
	clients  map[*Client]bool
	channels map[string]map[*Client]bool
	policies []*ChannelPolicy
	mutex    sync.RWMutex
}

type WSMessage struct {
	Type    string `json:"type"`
	Channel string `json:"channel"`
	From    string `json:"from"`
	Data    any    `json:"data"`
}

func NewHub() *Hub {
	return &Hub{
		clients:  make(map[*Client]bool),
		channels: make(map[string]map[*Client]bool),
	}
}

func (h *Hub) Register(c *Client) {
	h.mutex.Lock()
	h.clients[c] = true
	h.mutex.Unlock()
}

// Unregister - unregister client from the hub
func (h *Hub) Unregister(c *Client) {
	log.Println("Unregistering client")
	h.mutex.Lock()
	delete(h.clients, c)
	for ch := range c.Subs {
		delete(h.channels[ch], c)
		if len(h.channels[ch]) == 0 { // If there is user specific chan that would help like "user.{id}"
			delete(h.channels, ch)
		}
	}
	h.mutex.Unlock()
	close(c.Send)
}

// Unsubscribe - unsubscribe client from specific channel
func (h *Hub) Unsubscribe(c *Client, channel string) {
	h.mutex.Lock()

	if clients, ok := h.channels[channel]; ok {
		delete(clients, c)
		if len(clients) == 0 {
			delete(h.channels, channel)
		}
	}

	c.subMutex.Lock()
	delete(c.Subs, channel)
	c.subMutex.Unlock()
	h.mutex.Unlock()
}

// Subscribe - subscribe client to specific channel
func (h *Hub) Subscribe(c *Client, channel string) {
	h.mutex.Lock()
	if _, ok := h.channels[channel]; !ok {
		h.channels[channel] = make(map[*Client]bool)
	}
	h.channels[channel][c] = true
	c.subMutex.Lock()
	c.Subs[channel] = true
	c.subMutex.Unlock()
	h.mutex.Unlock()
}

// Broadcast  - broadcast message to specific channel
func (h *Hub) Broadcast(msg *WSMessage) {
	h.mutex.RLock()
	defer h.mutex.RUnlock()

	payload, err := json.Marshal(msg)
	if err != nil {
		log.Println("Failed to marshal message:", err)
		return
	}

	for registeredChannel, clients := range h.channels {
		if matchPattern(msg.Channel, registeredChannel) {
			for c := range clients {
				select {
				case c.Send <- payload:
				default:
					go h.Unregister(c)
				}
			}
		}
	}
}

// GetRegisteredChannels  - get list of registered channels
func (h *Hub) GetRegisteredChannels() []string {
	h.mutex.RLock()
	defer h.mutex.RUnlock()

	channels := make([]string, 0, len(h.channels))
	for ch := range h.channels {
		channels = append(channels, ch)
	}
	return channels
}

// getPolicy   - get policy for specific channel
func (h *Hub) getPolicy(channel string) *ChannelPolicy {
	h.mutex.RLock()
	defer h.mutex.RUnlock()
	for _, policy := range h.policies {
		if matchPattern(policy.Pattern, channel) {
			return policy
		}
	}
	return nil
}

// RegisterChannel register a new channel to subscribe on it using pattern and rules to allow access to it
func (h *Hub) RegisterChannel(channelPolicy *ChannelPolicy) {
	h.mutex.Lock()
	defer h.mutex.Unlock()
	h.policies = append(h.policies, channelPolicy)
}

// it just simple match pattern for wildcards only
func matchPattern(pattern, channel string) bool {
	if strings.HasSuffix(pattern, ".*") {
		prefix := strings.TrimSuffix(pattern, ".*")
		return strings.HasPrefix(channel, prefix+".")
	}
	return pattern == channel
}

// ReadPump - read messages from the websockets connection
func (c *Client) ReadPump(hub *Hub) {
	defer hub.Unregister(c)
	for {
		msgType, msg, err := c.Conn.ReadMessage()
		if err != nil {
			log.Printf("Read error: %v | msgType: %d | msg: %v", err, msgType, msg)
			break
		}

		if msgType != websocket.TextMessage {
			continue
		}

		var m map[string]string
		if err := json.Unmarshal(msg, &m); err != nil {
			continue
		}

		switch m["type"] {
		case "subscribe":
			channel := m["channel"]
			if channel != "" {
				policy := hub.getPolicy(channel)
				if policy != nil && policy.CanRead != nil && policy.CanRead(c.UserID, channel) {
					hub.Subscribe(c, channel)
				} else {
					log.Println("Unauthorized or unknown channel:", c.UserID, "→", channel)
				}
			}

		case "unsubscribe":
			channel := m["channel"]
			if channel != "" {
				hub.Unsubscribe(c, channel)
			}

		case "message":
			wsMessage, err := hub.receivedMessage(c, m)
			if err != nil {
				return
			}

			hub.Broadcast(wsMessage)

		default:
			log.Printf("Unknown message type: %s", m["type"])
		}
	}
}

func (h *Hub) receivedMessage(c *Client, m map[string]string) (*WSMessage, error) {
	channel := m["channel"]
	payload := m["data"]
	if channel == "" || !c.Subs[channel] {
		return nil, fmt.Errorf("unknown channel:  %s", channel)
	}

	policy := h.getPolicy(channel)
	if policy == nil || policy.CanWrite == nil || !policy.CanWrite(c.UserID, channel) {
		log.Println("unauthorized write attempt:", c.UserID, "→", channel)
		return nil, fmt.Errorf("unauthorized write attempt: %s → %s", c.UserID, channel)
	}

	return &WSMessage{
		Type:    "message",
		Channel: channel,
		From:    c.UserID,
		Data:    payload,
	}, nil
}

// WritePump - listen to send chan and writes messages to the websocket connection
func (c *Client) WritePump() {
	for msg := range c.Send {
		if err := c.Conn.WriteMessage(websocket.TextMessage, msg); err != nil {
			log.Println("Write error:", err)
			break
		}
	}
	c.Conn.Close()
}

// Listen - registers client to hub and starts WritePump and ReadPump[BLOCKING]
func (c *Client) Listen(hub *Hub) {
	hub.Register(c)
	go c.WritePump()
	c.ReadPump(hub)
}
