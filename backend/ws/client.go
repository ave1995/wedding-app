package ws

import (
	"log/slog"
	"maps"
	"time"

	"github.com/gorilla/websocket"
)

const (
	writeWait      = 10 * time.Second
	pongWait       = 60 * time.Second
	pingPeriod     = (pongWait * 9) / 10
	maxMessageSize = 512
)

type Client struct {
	conn   *websocket.Conn
	send   chan []byte
	topics map[string]bool
	logger *slog.Logger
}

func NewClient(conn *websocket.Conn, logger *slog.Logger) *Client {
	return &Client{
		conn:   conn,
		send:   make(chan []byte, 256),
		topics: make(map[string]bool),
		logger: logger,
	}
}

// ReadPump reads messages from the WebSocket and handles client disconnects
func (c *Client) ReadPump(h *Hub) {
	defer func() {
		h.unregister <- c
		c.conn.Close()
		c.logger.Info("Client disconnected", "remote", c.conn.RemoteAddr())
	}()

	c.conn.SetReadLimit(maxMessageSize)
	c.conn.SetReadDeadline(time.Now().Add(pongWait))
	c.conn.SetPongHandler(func(string) error {
		c.conn.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})

	for {
		_, msg, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				c.logger.Warn("Read error", "error", err, "remote", c.conn.RemoteAddr())
			}
			break
		}
		// Optional: handle messages sent by the client
		_ = msg
	}
}

// WritePump writes messages from the Send channel to the WebSocket
func (c *Client) WritePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.conn.Close()
		c.logger.Info("WritePump stopped", "remote", c.conn.RemoteAddr())
	}()

	for {
		select {
		case msg, ok := <-c.send:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				// Hub closed the channel
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				c.logger.Info("Channel closed by hub", "remote", c.conn.RemoteAddr())
				return
			}
			if err := c.conn.WriteMessage(websocket.TextMessage, msg); err != nil {
				c.logger.Warn("Write error", "error", err, "remote", c.conn.RemoteAddr())
				return
			}

		case <-ticker.C:
			// send ping to client
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				c.logger.Warn("Ping error", "error", err, "remote", c.conn.RemoteAddr())
				return
			}
		}
	}
}

func (c *Client) Topics() map[string]bool {
	copyTopics := make(map[string]bool, len(c.topics))
	maps.Copy(copyTopics, c.topics)
	return copyTopics
}

func (c *Client) Subscribe(topic string) {
	c.topics[topic] = true
	c.logger.Info("Subscribed to topic", "topic", topic)
}

func (c *Client) Unsubscribe(topic string) {
	delete(c.topics, topic)
	c.logger.Info("Unsubscribed from topic", "topic", topic)
}
