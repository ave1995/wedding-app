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
		broadcast:  make(chan broadcastMessage, 1024), // RADEK: to je celkem velký channel :D
		logger:     logger,
	}
}

func (h *Hub) Run() {
	for {
		select {
		// RADEK: register/unregistr by zmizel byly by to separe metody, které by zavírali mutech
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
			// RADEK: chybí ti tu validní uzavření toho hubu
			// RADEK: tohle bych dal do anonymní funkce abych mohl použít na otevření defer, což je vždycky lepší pokud uděláš lock a nemůžeš hned napsat defer unlock tak je něco špatně pravděpodobně....
			h.mu.RLock()
			for c := range h.clients {
				// RADEK: tenhle filtr bych napsal jako separe funkci, nebo funkci nad clientem => otestuješ to jednoduše
				if c.topics[msg.Topic] {
					select {
					case c.send <- msg.Data:
					default: // RADEK: když použiješ klasickej mutex a ne RW mutex tak se ti tohle nikdy nestane :D předpokládal bych, že registrace/deregistrace bude extrémně rychlá akce - tohle je trochu overengineering
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
	// RADEK: stejně jako unregistr, registr stačí mutex
	h.register <- c
}
