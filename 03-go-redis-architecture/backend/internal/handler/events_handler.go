package handler

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/seu-usuario/taskflow-backend/internal/middleware"
	"github.com/seu-usuario/taskflow-backend/internal/transport/sse"
)

type EventsHandler struct{ broker *sse.Broker }

func NewEventsHandler(broker *sse.Broker) *EventsHandler {
	return &EventsHandler{broker: broker}
}

func (h *EventsHandler) Stream(c *gin.Context) {
	userID, err := middleware.GetUserIDFromContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}
	if h.broker == nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{"error": "events unavailable"})
		return
	}
	events, unsubscribe, err := h.broker.Subscribe(userID)
	if err != nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{"error": "events unavailable"})
		return
	}
	defer unsubscribe()
	c.Header("Content-Type", "text/event-stream")
	c.Header("Cache-Control", "no-cache")
	c.Header("Connection", "keep-alive")
	c.Header("X-Accel-Buffering", "no")
	c.Status(http.StatusOK)
	flusher, ok := c.Writer.(http.Flusher)
	if !ok {
		return
	}
	flusher.Flush()
	keepAlive := time.NewTicker(25 * time.Second)
	defer keepAlive.Stop()
	for {
		select {
		case <-c.Request.Context().Done():
			return
		case <-keepAlive.C:
			if _, err := fmt.Fprint(c.Writer, ": keep-alive\n\n"); err != nil {
				return
			}
			flusher.Flush()
		case event, ok := <-events:
			if !ok {
				return
			}
			if _, err := fmt.Fprintf(c.Writer, "event: taskflow\ndata: %s\n\n", event); err != nil {
				return
			}
			flusher.Flush()
		}
	}
}
