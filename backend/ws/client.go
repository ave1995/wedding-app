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
	pingPeriod     = 5 * time.Second
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
	// c.conn.SetPongHandler(func(appData string) error {
	// 	c.logger.Info("Received pong", "remote", c.conn.RemoteAddr())
	// 	c.conn.SetReadDeadline(time.Now().Add(pongWait)) // extend deadline
	// 	return nil
	// })

	for {
		_, msg, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsCloseError(err, websocket.CloseNormalClosure, websocket.CloseGoingAway) {
				c.logger.Info("Client closed connection normally", "error", err, "remote", c.conn.RemoteAddr())
			} else if websocket.IsUnexpectedCloseError(err, websocket.CloseAbnormalClosure) {
				c.logger.Warn("Unexpected close error", "error", err, "remote", c.conn.RemoteAddr())
			} else {
				c.logger.Warn("ReadMessage error", "error", err, "remote", c.conn.RemoteAddr())
			}
			break
		}

		env, err := unwrapEvent[any](msg) // assuming you have wrapEvent/unwrapEvent
		if err != nil {
			c.logger.Warn("Failed to parse message", "error", err)
			continue
		}

		switch env.Topic {
		case TopicHeartbeatEvent:
			c.logger.Info("Received pong", "remote", c.conn.RemoteAddr())
			c.conn.SetReadDeadline(time.Now().Add(pongWait)) // extend deadline
		}
	}
}

const TopicHeartbeatEvent = "heartbeat"

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

			data, err := wrapEvent(TopicHeartbeatEvent, "ping!")
			if err != nil {
				c.logger.Warn("failed wrap data for heartbeat", "error", err, "remote", c.conn.RemoteAddr())
				return
			}

			if err := c.conn.WriteMessage(websocket.TextMessage, data); err != nil {
				c.logger.Warn("Ping error", "error", err, "remote", c.conn.RemoteAddr())
				return
			}
			c.logger.Info("Sent ping", "remote", c.conn.RemoteAddr())
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
