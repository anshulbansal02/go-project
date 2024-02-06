package websockets

import "sync"

// Stores all the websocket connections
type clientPool struct {
	clients map[string]*Client
	mu      sync.RWMutex
}

// Adds a new client to clients pool
func (cp *clientPool) add(c *Client) {
	cp.mu.Lock()
	defer cp.mu.Unlock()
	cp.clients[c.ID] = c
}

// Removes client from clients pool
func (cp *clientPool) remove(c *Client) {
	cp.mu.Lock()
	defer cp.mu.Unlock()
	delete(cp.clients, c.ID)
}

func (cp *clientPool) get(cId string) *Client {
	cp.mu.RLock()
	defer cp.mu.RUnlock()

	client := cp.clients[cId]
	return client
}
