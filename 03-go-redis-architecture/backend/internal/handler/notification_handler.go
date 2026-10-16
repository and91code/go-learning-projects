package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"github.com/seu-usuario/taskflow-backend/internal/middleware"
)

type NotificationHandler struct{ redis *redis.Client }

func NewNotificationHandler(client *redis.Client) *NotificationHandler {
	return &NotificationHandler{redis: client}
}

type notification struct {
	ID           string    `json:"id"`
	Message      string    `json:"message"`
	ResourceType string    `json:"resource_type"`
	ResourceID   int64     `json:"resource_id"`
	Action       string    `json:"action"`
	Timestamp    time.Time `json:"timestamp"`
	Unread       bool      `json:"unread"`
}

func (h *NotificationHandler) key(c *gin.Context) (string, error) {
	userID, err := middleware.GetUserIDFromContext(c)
	return fmt.Sprintf("taskflow:user:%d:notifications", userID), err
}

func (h *NotificationHandler) List(c *gin.Context) {
	key, err := h.key(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}
	if h.redis == nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{"error": "notifications unavailable"})
		return
	}
	ids, err := h.redis.ZRevRange(c.Request.Context(), key+":index", 0, -1).Result()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not load notifications"})
		return
	}
	if len(ids) == 0 {
		c.JSON(http.StatusOK, gin.H{"notifications": []notification{}})
		return
	}
	values, err := h.redis.HMGet(c.Request.Context(), key, ids...).Result()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not load notifications"})
		return
	}
	items := make([]notification, 0, len(values))
	for _, value := range values {
		raw, ok := value.(string)
		if !ok {
			continue
		}
		var item notification
		if json.Unmarshal([]byte(raw), &item) == nil {
			items = append(items, item)
		}
	}
	c.JSON(http.StatusOK, gin.H{"notifications": items})
}

func (h *NotificationHandler) MarkRead(c *gin.Context) {
	h.update(c, true)
}

func (h *NotificationHandler) Delete(c *gin.Context) {
	key, err := h.key(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}
	if h.redis == nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{"error": "notifications unavailable"})
		return
	}
	id := c.Param("id")
	if _, err := h.redis.TxPipelined(c.Request.Context(), func(pipe redis.Pipeliner) error {
		pipe.HDel(c.Request.Context(), key, id)
		pipe.ZRem(c.Request.Context(), key+":index", id)
		return nil
	}); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not delete notification"})
		return
	}
	c.Status(http.StatusNoContent)
}

func (h *NotificationHandler) update(c *gin.Context, read bool) {
	key, err := h.key(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}
	if h.redis == nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{"error": "notifications unavailable"})
		return
	}
	ctx := c.Request.Context()
	id := c.Param("id")
	if read {
		raw, err := h.redis.Eval(ctx, `
			local value = redis.call('HGET', KEYS[1], ARGV[1])
			if not value then return nil end
			local item = cjson.decode(value)
			item.unread = false
			local updated = cjson.encode(item)
			redis.call('HSET', KEYS[1], ARGV[1], updated)
			return updated
		`, []string{key}, id).Text()
		if err != nil {
			if err == redis.Nil {
				c.JSON(http.StatusNotFound, gin.H{"error": "notification not found"})
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{"error": "could not update notification"})
			return
		}
		var item notification
		if json.Unmarshal([]byte(raw), &item) != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "invalid notification"})
			return
		}
		c.JSON(http.StatusOK, item)
		return
	}
	raw, err := h.redis.HGet(ctx, key, id).Result()
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "notification not found"})
		return
	}
	var item notification
	if json.Unmarshal([]byte(raw), &item) != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "invalid notification"})
		return
	}
	item.Unread = false
	encoded, _ := json.Marshal(item)
	if err := h.redis.HSet(ctx, key, id, encoded).Err(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not update notification"})
		return
	}
	c.JSON(http.StatusOK, item)
}
