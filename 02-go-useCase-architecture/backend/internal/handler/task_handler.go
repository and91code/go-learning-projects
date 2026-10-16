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

type TaskHandler struct {
	create interface {
		Execute(context.Context, int64, types.Task) (types.Task, error)
	}
	update interface {
		Execute(context.Context, int64, int64, types.Task) (types.Task, error)
	}
	delete interface {
		Execute(context.Context, int64, int64) error
	}
	byList interface {
		Execute(context.Context, int64, int64) ([]types.Task, error)
	}
	byID interface {
		Execute(context.Context, int64, int64) (types.Task, error)
	}
}

func NewTaskHandler(create interface {
	Execute(context.Context, int64, types.Task) (types.Task, error)
}, update interface {
	Execute(context.Context, int64, int64, types.Task) (types.Task, error)
}, delete interface {
	Execute(context.Context, int64, int64) error
}, byList interface {
	Execute(context.Context, int64, int64) ([]types.Task, error)
}, byID interface {
	Execute(context.Context, int64, int64) (types.Task, error)
}) *TaskHandler {
	return &TaskHandler{create: create, update: update, delete: delete, byList: byList, byID: byID}
}

func (h *TaskHandler) Store(c *gin.Context) {
	userID, err := middleware.GetUserIDFromContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	var req dto.CreateTaskRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": "Dados inválidos", "details": err.Error()})
		return
	}

	task, err := h.create.Execute(c.Request.Context(), userID, req.ToDomain())
	if err != nil {
		if errors.Is(err, types.ErrListNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Lista associada não encontrada"})
			return
		}
		if errors.Is(err, types.ErrTagNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Uma ou mais tags não pertencem à lista"})
			return
		}
		if errors.Is(err, types.ErrInvalidTask) {
			c.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, dto.TaskResponseFromDomain(task))
}

func (h *TaskHandler) IndexByList(c *gin.Context) {
	userID, err := middleware.GetUserIDFromContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	listID, err := strconv.ParseInt(c.Query("list_id"), 10, 64)
	if err != nil || listID <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID de lista inválido"})
		return
	}

	tasks, err := h.byList.Execute(c.Request.Context(), listID, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, dto.TaskResponsesFromDomain(tasks))
}

func (h *TaskHandler) Show(c *gin.Context) {
	userID, err := middleware.GetUserIDFromContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}
	taskID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil || taskID <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID de tarefa inválido"})
		return
	}

	task, err := h.byID.Execute(c.Request.Context(), taskID, userID)
	if err != nil {
		if errors.Is(err, types.ErrTaskNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, dto.TaskResponseFromDomain(task))
}

func (h *TaskHandler) Update(c *gin.Context) {
	userID, err := middleware.GetUserIDFromContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	taskID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil || taskID <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID de tarefa inválido"})
		return
	}

	var req dto.UpdateTaskRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": "Dados inválidos", "details": err.Error()})
		return
	}

	task, err := h.update.Execute(c.Request.Context(), taskID, userID, req.ToDomain())
	if err != nil {
		if errors.Is(err, types.ErrTaskNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		if errors.Is(err, types.ErrTagNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Uma ou mais tags não pertencem à lista"})
			return
		}
		if errors.Is(err, types.ErrInvalidTask) {
			c.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, dto.TaskResponseFromDomain(task))
}

func (h *TaskHandler) Delete(c *gin.Context) {
	userID, err := middleware.GetUserIDFromContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	taskID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil || taskID <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID de tarefa inválido"})
		return
	}

	if err := h.delete.Execute(c.Request.Context(), taskID, userID); err != nil {
		if errors.Is(err, types.ErrTaskNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		if errors.Is(err, types.ErrInvalidTask) {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}
