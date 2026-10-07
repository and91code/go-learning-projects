package handler

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/seu-usuario/taskflow-backend/internal/dto"
	"github.com/seu-usuario/taskflow-backend/internal/middleware"
	"github.com/seu-usuario/taskflow-backend/internal/repository"
	"github.com/seu-usuario/taskflow-backend/internal/service"
)

type ListHandler struct {
	listService service.ListService
}

func NewListHandler(s service.ListService) *ListHandler {
	return &ListHandler{listService: s}
}

func (h *ListHandler) Store(c *gin.Context) {
	userID := c.MustGet(middleware.UserIDKey).(int64)

	var req dto.CreateListRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": "Dados inválidos", "details": err.Error()})
		return
	}

	res, err := h.listService.CreateList(c.Request.Context(), userID, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, res)
}

func (h *ListHandler) Index(c *gin.Context) {
	userID := c.MustGet(middleware.UserIDKey).(int64)

	res, err := h.listService.GetUserLists(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, res)
}

func (h *ListHandler) Show(c *gin.Context) {
	userID := c.MustGet(middleware.UserIDKey).(int64)
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID de lista inválido"})
		return
	}

	res, err := h.listService.GetListByID(c.Request.Context(), id, userID)
	if err != nil {
		if errors.Is(err, repository.ErrListNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, res)
}
