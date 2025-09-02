package restapi

import (
	"net/http"
	"slices"
	"strings"
	"wedding-app/config"
	"wedding-app/ws"

	"log/slog"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

type WSHandler struct {
	hub      *ws.Hub
	logger   *slog.Logger
	upgrader websocket.Upgrader
}

func NewWSHandler(hub *ws.Hub, logger *slog.Logger, config config.ServerConfig) *WSHandler {
	return &WSHandler{
		hub:    hub,
		logger: logger,
		upgrader: websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool {
				origin := r.Header.Get("Origin")
				return slices.Contains(config.Origins, origin)
			},
		}}
}

func (h *WSHandler) serveWS(c *gin.Context) {
	conn, err := h.upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		h.logger.Error("Failed to upgrade connection", "error", err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	client := ws.NewClient(conn, h.logger)

	// Parse topics from query param, e.g. /ws?topics=answers,questions
	q := c.Query("topics")
	for _, t := range splitTopics(q) {
		client.Subscribe(t)
	}

	h.hub.RegisterClient(client)

	// Start pumps
	go client.WritePump()
	go client.ReadPump(h.hub)

	h.logger.Info("WebSocket client connected", "remote", conn.RemoteAddr(), "topics", client.Topics())
}

func splitTopics(q string) []string {
	if q == "" {
		return nil
	}
	return strings.Split(q, ",")
}
