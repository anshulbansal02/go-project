package websockets

import (
	"anshulbansal02/scribbly/pkg/utils"

	"fmt"
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

var generateClientId = utils.NewRandomStringGenerator(utils.CHARSET_URL_SAFE, 8)

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
		writeChannel: make(chan OutgoingWebSocketMessage),
		manager:      m,
	}

	m.clientPool.add(client)

	go client.readLoop()
	go client.writeLoop()

	fmt.Println("Client Connected: ", client.ID)
}

// Process a single client websocket emitted message
func (m *WebSocketManager) processMessage(client *Client, _ int, msg []byte) {
	message, err := DecodeMessage(msg)
	if err != nil {
		fmt.Println("Error decoding message: ", err)
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
	return m.clientPool.get(clientId)
}

func (m *WebSocketManager) EmitTo(clientId string, message OutgoingWebSocketMessage) {
	client := m.GetClient(clientId)
	client.Emit(message)
}

func (m *WebSocketManager) Multicast(clientIds []string, exceptIds []string, message OutgoingWebSocketMessage) {
	for _, clientId := range clientIds {
		if slices.Contains(exceptIds, clientId) {
			continue
		}

		m.GetClient(clientId).Emit(message)
	}
}
