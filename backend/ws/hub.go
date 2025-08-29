package ws

import (
	"log/slog"
	"sync"
)

type Hub struct {
	clients    map[*Client]bool
	register   chan *Client
	unregister chan *Client
	broadcast  chan broadcastMessage
	mu         sync.RWMutex
	logger     *slog.Logger
}

type broadcastMessage struct {
	Topic string
	Data  []byte
}

func NewHub(logger *slog.Logger) *Hub {
	return &Hub{
		clients:    make(map[*Client]bool),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		broadcast:  make(chan broadcastMessage, 1024),
		logger:     logger,
	}
}

func (h *Hub) Run() {
	for {
		select {
		case c := <-h.register:
			h.mu.Lock()
			h.clients[c] = true
			h.mu.Unlock()
			h.logger.Info("Client registered", "remote", c.conn.RemoteAddr())

		case c := <-h.unregister:
			h.mu.Lock()
			if _, ok := h.clients[c]; ok {
				delete(h.clients, c)
				close(c.send)
				h.logger.Info("Client unregistered", "remote", c.conn.RemoteAddr())
			}
			h.mu.Unlock()

		case msg := <-h.broadcast:
			h.mu.RLock()
			for c := range h.clients {
				if c.topics[msg.Topic] {
					select {
					case c.send <- msg.Data:
					default:
						close(c.send)
						delete(h.clients, c)
						h.logger.Warn("Dropped slow client", "remote", c.conn.RemoteAddr())
					}
				}
			}
			h.mu.RUnlock()
		}
	}
}

func (h *Hub) RegisterClient(c *Client) {
	h.register <- c
}
