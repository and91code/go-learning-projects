package handler

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/seu-usuario/taskflow-backend/internal/dto"
	"github.com/seu-usuario/taskflow-backend/internal/middleware"
	"github.com/seu-usuario/taskflow-backend/internal/queue"
	"github.com/seu-usuario/taskflow-backend/internal/types"
)

type ListHandler struct {
	create interface {
		Execute(context.Context, int64, types.List) (types.List, error)
	}
	get interface {
		Execute(context.Context, int64, int64) (types.List, error)
	}
	list interface {
		Execute(context.Context, int64) ([]types.List, error)
	}
	update interface {
		Execute(context.Context, int64, int64, types.List) (types.List, error)
	}
	delete interface {
		Execute(context.Context, int64, int64) error
	}
	jobs queue.Producer
}

func (h *ListHandler) SetQueue(q queue.Producer) { h.jobs = q }

func NewListHandler(create interface {
	Execute(context.Context, int64, types.List) (types.List, error)
}, get interface {
	Execute(context.Context, int64, int64) (types.List, error)
}, list interface {
	Execute(context.Context, int64) ([]types.List, error)
}, update interface {
	Execute(context.Context, int64, int64, types.List) (types.List, error)
}, del interface {
	Execute(context.Context, int64, int64) error
}) *ListHandler {
	return &ListHandler{create: create, get: get, list: list, update: update, delete: del}
}

func (h *ListHandler) Store(c *gin.Context) {
	userID := c.MustGet(middleware.UserIDKey).(int64)

	var req dto.CreateListRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": "Dados inválidos", "details": err.Error()})
		return
	}

	if h.jobs != nil {
		domain := req.ToDomain()
		if err := queue.ValidateList(userID, 0, domain); err != nil {
			c.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
			return
		}
		payload, err := json.Marshal(domain)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "unable to encode list"})
			return
		}
		if err := h.jobs.Enqueue(c.Request.Context(), queue.Job{Action: "LIST_CREATE", UserID: userID, Payload: payload}); err != nil {
			c.JSON(http.StatusServiceUnavailable, gin.H{"error": "unable to queue list"})
			return
		}
		c.JSON(http.StatusAccepted, gin.H{"status": "accepted"})
		return
	}
	list, err := h.create.Execute(c.Request.Context(), userID, req.ToDomain())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, dto.ListResponseFromDomain(list))
}

func (h *ListHandler) Index(c *gin.Context) {
	userID := c.MustGet(middleware.UserIDKey).(int64)

	lists, err := h.list.Execute(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, dto.ListResponsesFromDomain(lists))
}

func (h *ListHandler) Show(c *gin.Context) {
	userID := c.MustGet(middleware.UserIDKey).(int64)
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID de lista inválido"})
		return
	}

	list, err := h.get.Execute(c.Request.Context(), id, userID)
	if err != nil {
		if errors.Is(err, types.ErrListNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, dto.ListResponseFromDomain(list))
}

func (h *ListHandler) Update(c *gin.Context) {
	userID := c.MustGet(middleware.UserIDKey).(int64)
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil || id <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID de lista inválido"})
		return
	}

	var req dto.UpdateListRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": "Dados inválidos", "details": err.Error()})
		return
	}

	if h.jobs != nil {
		domain := req.ToDomain()
		if err := queue.ValidateList(userID, id, domain); err != nil {
			c.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
			return
		}
		payload, err := json.Marshal(domain)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "unable to encode list"})
			return
		}
		if err := h.jobs.Enqueue(c.Request.Context(), queue.Job{Action: "LIST_UPDATE", UserID: userID, EntityID: id, Payload: payload}); err != nil {
			c.JSON(http.StatusServiceUnavailable, gin.H{"error": "unable to queue list"})
			return
		}
		c.JSON(http.StatusAccepted, gin.H{"status": "accepted"})
		return
	}
	list, err := h.update.Execute(c.Request.Context(), id, userID, req.ToDomain())
	if err != nil {
		if errors.Is(err, types.ErrListNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, dto.ListResponseFromDomain(list))
}

func (h *ListHandler) Delete(c *gin.Context) {
	userID := c.MustGet(middleware.UserIDKey).(int64)
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil || id <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID de lista inválido"})
		return
	}

	if h.jobs != nil {
		if err := h.jobs.Enqueue(c.Request.Context(), queue.Job{Action: "LIST_DELETE", UserID: userID, EntityID: id}); err != nil {
			c.JSON(http.StatusServiceUnavailable, gin.H{"error": "unable to queue list"})
			return
		}
		c.JSON(http.StatusAccepted, gin.H{"status": "accepted"})
		return
	}
	if err := h.delete.Execute(c.Request.Context(), id, userID); err != nil {
		if errors.Is(err, types.ErrListNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}
