package websockets

import (
	"fmt"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

// Represents a websocket connection
type Client struct {
	ID           string
	socket       *websocket.Conn
	writeChannel chan OutgoingWebSocketMessage
	manager      *WebSocketManager
	closeOnce    sync.Once
}

// Closes and cleans up client from all places
func (c *Client) closeAndCleanup() {
	c.closeOnce.Do(func() {
		// Remove client from client pool
		c.manager.clientPool.remove(c)
		// Remove all observers for client
		c.manager.hub.RemoveObserversForClient(c)
		// Close writeChannel
		close(c.writeChannel)
		// Close underlying socket connection after sending close control message
		c.socket.SetWriteDeadline(time.Now().Add(writeWait))
		c.socket.WriteMessage(websocket.CloseMessage, []byte{})
		c.socket.Close()

		fmt.Println("Client Disconnected: ", c.ID)
	})
}

// Reads incoming messages from websocket client and sends back to manager's processMessage
func (c *Client) readLoop() {
	defer c.closeAndCleanup()

	c.socket.SetReadLimit(maxMessageSize)
	c.socket.SetReadDeadline(time.Now().Add(pongWait))

	c.socket.SetPongHandler(func(string) error {
		c.socket.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})

	for {
		mt, msg, err := c.socket.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				fmt.Println("ReadLoop Error: ", err)
			}
			break
		}

		go c.manager.processMessage(c, mt, msg)
	}
}

// Central writter to write all buffered messages to websocket client
func (c *Client) writeLoop() {
	// Ticker for regular ping control messages
	ticker := time.NewTicker(pingPeriod)

	defer func() {
		ticker.Stop()
		c.closeAndCleanup()
	}()

	for {
		select {
		case message, ok := <-c.writeChannel:
			// The writeChannel was closed
			if !ok {
				return
			}

			c.socket.SetWriteDeadline(time.Now().Add(writeWait))

			w, err := c.socket.NextWriter(websocket.BinaryMessage)
			if err != nil {
				return
			}

			msg, err := EncodeMessage(&message)
			if err == nil {
				w.Write([]byte(msg))
			}

			// Add queued messages to the current websocket message.
			n := len(c.writeChannel)
			for i := 0; i < n; i++ {
				w.Write(newline)
				message, ok = <-c.writeChannel
				if !ok {
					return
				}
				msg, err = EncodeMessage(&message)
				if err == nil {
					w.Write([]byte(msg))
				}
			}

			if err := w.Close(); err != nil {
				return
			}
		case <-ticker.C:
			c.socket.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.socket.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

// Emits a websocket message to the client
func (c *Client) Emit(e OutgoingWebSocketMessage) {
	c.writeChannel <- e
}

// Adds an observer for the event from the client and returns an observer id for registered observer
func (c *Client) AddObserver(event Event, observer Observer) string {
	return c.manager.hub.AddObserver(c, &event, observer)
}
