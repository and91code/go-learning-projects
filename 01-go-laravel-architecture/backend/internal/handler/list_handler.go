package handler

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/seu-usuario/taskflow-backend/internal/dto"
	"github.com/seu-usuario/taskflow-backend/internal/middleware"
	"github.com/seu-usuario/taskflow-backend/internal/types"
)

type ListHandler struct {
	listService types.ListService
}

func NewListHandler(s types.ListService) *ListHandler {
	return &ListHandler{listService: s}
}

func (h *ListHandler) Store(c *gin.Context) {
	userID := c.MustGet(middleware.UserIDKey).(int64)

	var req dto.CreateListRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": "Dados inválidos", "details": err.Error()})
		return
	}

	list, err := h.listService.CreateList(c.Request.Context(), userID, req.ToDomain())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, dto.ListResponseFromDomain(list))
}

func (h *ListHandler) Index(c *gin.Context) {
	userID := c.MustGet(middleware.UserIDKey).(int64)

	lists, err := h.listService.GetUserLists(c.Request.Context(), userID)
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

	list, err := h.listService.GetListByID(c.Request.Context(), id, userID)
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
