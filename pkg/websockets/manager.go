package websockets

import (
	"fmt"
	"net/http"

	"github.com/gorilla/websocket"
)

type WebSocketEvent struct {
	Name    string
	Message string
}

/* --------------- CLIENT ---------------- */

// Represents a websocket connection
type Client struct {
	ID        string
	socket    *websocket.Conn
	receivers []func(event WebSocketEvent)
}

// Send string message over to websocket client
func (c *Client) Send(message string) error {
	return c.socket.WriteMessage(websocket.TextMessage, []byte(message))
}

// Register a functional handler for an event on the client
func (c *Client) On(event string, handler func(event WebSocketEvent)) {
	c.receivers = append(c.receivers, handler)
}

// Process a single client websocket emitted message
func (c *Client) processMessage(mt int, msg []byte) {
	for _, handler := range c.receivers {
		handler(WebSocketEvent{
			Name:    "Message",
			Message: string(msg),
		})
	}
}

/* --------------- CLIENT POOL ---------------- */

// Stores all the websocket connections
type clientPool struct {
	clients map[string]*Client
}

// Add a new client to clients pool
func (cp *clientPool) addNew(c *Client) {
	cp.clients[c.ID] = c
}

// Close websocket client connection and remove from clients pool
func (cp *clientPool) closeAndRemove(c *Client) {
	delete(cp.clients, c.ID)
	c.socket.Close()
}

/* --------------- WEBSOCKET MANAGER ---------------- */

type webSocketManager struct {
	clientPool   clientPool
	connUpgrader websocket.Upgrader
}

func NewWebSocketManager() *webSocketManager {
	return &webSocketManager{
		clientPool: clientPool{
			clients: make(map[string]*Client),
		},
		connUpgrader: websocket.Upgrader{
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,

			CheckOrigin: func(r *http.Request) bool {
				return true
			},
		},
	}
}

// Upgrade an HTTP connection to WebSocket one
func (m *webSocketManager) HandleWSConnection(w http.ResponseWriter, r *http.Request) {
	conn, err := m.connUpgrader.Upgrade(w, r, nil)

	if err != nil {
		// Send http error
		return
	}

	client := &Client{ID: "", socket: conn}
	m.clientPool.addNew(client)

	go m.readLoop(client)

}

// Poll for messages from websocket clients
func (m *webSocketManager) readLoop(client *Client) {
	defer m.clientPool.closeAndRemove(client)
	for {
		mt, msgBytes, err := client.socket.ReadMessage()

		if err != nil {
			fmt.Println("ReadLoop Error: ", err)
		}

		go client.processMessage(mt, msgBytes)
	}
}
