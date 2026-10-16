package handler

import (
	"context"
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/seu-usuario/taskflow-backend/internal/dto"
	"github.com/seu-usuario/taskflow-backend/internal/middleware"
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
}

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
