package websockets

import (
	"anshulbansal02/scribbly/pkg/utils"
	"fmt"
	"strings"
	"sync"
)

type Observer func(m WebSocketMessage, c *Client)

type observerWithId struct {
	id       string
	observer Observer
}

type Hub struct {
	mutex         sync.Mutex
	observers     map[string][]observerWithId
	observerAtKey map[string]string
}

func NewHub() *Hub {
	return &Hub{
		observers:     make(map[string][]observerWithId),
		observerAtKey: make(map[string]string),
	}
}

var generateObserverId = utils.NewRandomStringGenerator(nil, 6)

// Creates a mapping key for given client and event to store observers in a flat map
func getObserverSlotKey(client *Client, event *Event) string {
	var key1, key2 string

	if client == nil {
		key1 = "*"
	}
	if event == nil {
		key2 = "*"
	}

	return fmt.Sprintf("%v$%v", key1, key2)
}

// Registers an observer for given client and event and returns an observerId for later removal
func (h *Hub) AddObserver(client *Client, event *Event, observer Observer) string {
	h.mutex.Lock()
	defer h.mutex.Unlock()

	newObserver := observerWithId{
		id:       generateObserverId(),
		observer: observer,
	}

	key := getObserverSlotKey(client, event)

	h.observerAtKey[newObserver.id] = key
	h.observers[key] = append(h.observers[key], newObserver)

	return newObserver.id
}

// Removes an observer using its observerId and returns a bool indicating if the observer was found or not
func (h *Hub) RemoveObserver(observerId string) bool {
	h.mutex.Lock()
	defer h.mutex.Unlock()

	key, ok := h.observerAtKey[observerId]
	if !ok {
		return false
	}

	observers := h.observers[key]

	var removalIndex int = -1
	for i, observer := range observers {
		if observer.id == observerId {
			removalIndex = i
			break
		}
	}
	if removalIndex == -1 {
		return false
	}

	observers[removalIndex] = observers[len(observers)-1]
	observers = observers[:len(observers)-1]

	return true
}

// Removes all observers for given client
func (h *Hub) RemoveObserversForClient(client *Client) {
	h.mutex.Lock()
	defer h.mutex.Unlock()

	for key, _ := range h.observers {
		if strings.HasPrefix(key, client.ID) {
			delete(h.observers, key)
		}
	}

	for observerId, key := range h.observerAtKey {
		if strings.HasPrefix(key, client.ID) {
			delete(h.observerAtKey, observerId)
		}
	}

}

// Removes all observers for given event
func (h *Hub) RemoveObserversForEvent(event Event) {
	h.mutex.Lock()
	defer h.mutex.Unlock()

	for key, _ := range h.observers {
		if strings.HasPrefix(key, string(event)) {
			delete(h.observers, key)
		}
	}

	for observerId, key := range h.observerAtKey {
		if strings.HasPrefix(key, string(event)) {
			delete(h.observerAtKey, observerId)
		}
	}
}

// Dispatches incoming message from a client to all registered observers
func (h *Hub) DispatchMessage(client *Client, message WebSocketMessage) {
	event := message.EventName

	// For observers bound to client and event both
	for _, handle := range h.observers[getObserverSlotKey(client, &event)] {
		go handle.observer(message, client)
	}
	// For observers bound to all clients
	for _, handle := range h.observers[getObserverSlotKey(nil, &event)] {
		go handle.observer(message, client)
	}
	// For observers bound to all events
	for _, handle := range h.observers[getObserverSlotKey(client, nil)] {
		go handle.observer(message, client)
	}
	// For observers bound to all
	for _, handle := range h.observers[getObserverSlotKey(nil, nil)] {
		go handle.observer(message, client)
	}

}
