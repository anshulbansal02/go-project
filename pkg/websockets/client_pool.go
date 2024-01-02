package websockets

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
