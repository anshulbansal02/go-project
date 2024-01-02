package websockets

import "github.com/gorilla/websocket"

// Represents a websocket connection
type Client struct {
	ID        string
	socket    *websocket.Conn
	listeners []func(event WebSocketEvent)
}

// Send string message over to websocket client
func (c *Client) Send(message string) error {
	return c.socket.WriteMessage(websocket.TextMessage, []byte(message))
}

// Register a functional handler for an event on the client
func (c *Client) On(event string, handler func(event WebSocketEvent)) {
	c.listeners = append(c.listeners, handler)
}
