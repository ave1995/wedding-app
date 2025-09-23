package restapi

// RADEK: to řazení je hrozný :D
import (
	"log/slog"
	"net/http"
	"slices"
	"strings"

	"wedding-app/config"
	"wedding-app/ws"

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

	// RADEK: tohle není zase moc šikovně udělané protože celý ten tvuj "ws" systém je závyslí na websocketech - ted když tohle budeš chtít přepsat do GRPC tak to budeš muset celý překopat
	// websocket by měl zůstat v tomhle handleru a nikam jinam se už nepodívat protože to je jeho domena API/HTTP
	// to tvoje WS má implementovat connector a být ve složce connectors, jako implementace nějakého pub/sub, tedy connektor, který má funkce jako (register, unregister, send, receive)
	// Tady mají být ty loopy pro WS a jen iterovat s tím connektor interface -> pokud by si potom chtěl změnit technologii, tak použiješ stejný konnektor, jen ho budeš krmit jinudy
	// Také by ti to jednoduše povolilo používat oba prokoloy najednou nebo třeba HTTP napul s grpc
	client := ws.NewClient(conn, h.logger)

	// Parse topics from query param, e.g. /ws?topics=answers,questions
	// RADEK: tohle klidně mohl dostat jako mapu/pole/custom typ už v tom konstruktoru, a třeba by si nejdřív než zaregistruješ klienta mohl zkontrolovat zda vůbec nějaký topic přišel.
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

// RADEK: tohle je takzvaně funkce k kničemu :D
// Když bude prázdný string, tak to vrátí prázdné pole a ten range nic neudělá :|
func splitTopics(q string) []string {
	if q == "" {
		return nil
	}
	return strings.Split(q, ",")
}
