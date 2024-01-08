package websockets

// Stores all the websocket connections
type clientPool struct {
	clients map[string]*Client
}

// Adds a new client to clients pool
func (cp *clientPool) add(c *Client) {
	cp.clients[c.ID] = c
}

// Removes client from clients pool
func (cp *clientPool) remove(c *Client) {
	delete(cp.clients, c.ID)
}
