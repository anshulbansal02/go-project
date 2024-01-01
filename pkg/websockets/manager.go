package websockets

import (
	"anshulbansal02/scribbly/pkg/utils"
	"fmt"
	"net/http"
	"slices"

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

type WebSocketManager struct {
	clientPool   clientPool
	connUpgrader websocket.Upgrader
	listeners    []func(client *Client, event WebSocketEvent)
}

func NewWebSocketManager() *WebSocketManager {
	return &WebSocketManager{
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

var generateClientId = utils.NewRandomStringGenerator(nil, 8)

// Upgrade an HTTP connection to WebSocket one
func (m *WebSocketManager) HandleWSConnection(w http.ResponseWriter, r *http.Request) {
	conn, err := m.connUpgrader.Upgrade(w, r, nil)

	if err != nil {
		// Send http error
		return
	}

	client := &Client{ID: generateClientId(), socket: conn}
	m.clientPool.addNew(client)

	go m.readLoop(client)

}

// Poll for messages from websocket clients
func (m *WebSocketManager) readLoop(client *Client) {
	defer m.clientPool.closeAndRemove(client)
	for {
		mt, msgBytes, err := client.socket.ReadMessage()

		if err != nil {
			fmt.Println("ReadLoop Error: ", err)
		}

		go m.processMessage(client, mt, msgBytes)
	}
}

// Process a single client websocket emitted message
func (m *WebSocketManager) processMessage(client *Client, mt int, msg []byte) {
	// Extract event name from msg

	for _, handler := range m.listeners {
		handler(client, WebSocketEvent{
			Name:    "Message",
			Message: string(msg),
		})
	}

	for _, handler := range client.listeners {
		handler(WebSocketEvent{
			Name:    "Message",
			Message: string(msg),
		})
	}

}

func (m *WebSocketManager) GetClient(clientId string) *Client {
	return m.clientPool.clients[clientId]
}

func (m *WebSocketManager) Send(clientId, message string) {
	m.GetClient(clientId).Send(message)
}

// Register a functional handler for an event on the manager for all clients
func (m *WebSocketManager) OnEvent(event string, handler func(client *Client, event WebSocketEvent)) {
	m.listeners = append(m.listeners, handler)
}

func (m *WebSocketManager) Multicast(clientIds []string, exceptIds []string, message string) {
	for _, clientId := range clientIds {
		if slices.Contains(exceptIds, clientId) {
			continue
		}

		m.GetClient(clientId).Send(message)

	}
}
