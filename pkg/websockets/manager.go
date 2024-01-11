package websockets

import (
	"anshulbansal02/scribbly/pkg/utils"
	"net/http"
	"slices"

	"github.com/gorilla/websocket"
)

type WebSocketManager struct {
	clientPool   clientPool
	connUpgrader websocket.Upgrader
	hub          Hub
}

func NewWebSocketManager() *WebSocketManager {
	return &WebSocketManager{
		clientPool: clientPool{
			clients: make(map[string]*Client),
		},
		hub: *NewHub(),
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

// Upgrades an HTTP connection to WebSocket one
func (m *WebSocketManager) HandleWSConnection(w http.ResponseWriter, r *http.Request) {
	conn, err := m.connUpgrader.Upgrade(w, r, nil)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	client := &Client{
		ID:           generateClientId(),
		socket:       conn,
		writeChannel: make(chan WebSocketMessage),
		manager:      m,
	}

	m.clientPool.add(client)

	go client.readLoop()
	go client.writeLoop()
}

// Process a single client websocket emitted message
func (m *WebSocketManager) processMessage(client *Client, mt int, msg []byte) {
	message, err := DecodeMessage(msg)
	if err != nil {
		return
	}

	m.hub.DispatchMessage(client, *message)
}

func (m *WebSocketManager) AddObserver(event Event, observer Observer) {
	m.hub.AddObserver(nil, &event, observer)
}

func (m *WebSocketManager) RemoveObserver(observerId string) {
	m.hub.RemoveObserver(observerId)
}

func (m *WebSocketManager) GetClient(clientId string) *Client {
	return m.clientPool.clients[clientId]
}

func (m *WebSocketManager) EmitTo(clientId string, message WebSocketMessage) {

}

func (m *WebSocketManager) Multicast(clientIds []string, exceptIds []string, message WebSocketMessage) {
	for _, clientId := range clientIds {
		if slices.Contains(exceptIds, clientId) {
			continue
		}

		m.GetClient(clientId).Emit(message)
	}
}
