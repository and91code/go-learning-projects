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

type TaskHandler struct {
	taskService types.TaskService
}

func NewTaskHandler(s types.TaskService) *TaskHandler {
	return &TaskHandler{taskService: s}
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

	task, err := h.taskService.CreateTask(c.Request.Context(), userID, req.ToDomain())
	if err != nil {
		if errors.Is(err, types.ErrListNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Lista associada não encontrada"})
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

	tasks, err := h.taskService.GetTasksByList(c.Request.Context(), listID, userID)
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

	task, err := h.taskService.GetTaskByID(c.Request.Context(), taskID, userID)
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

	task, err := h.taskService.UpdateTask(c.Request.Context(), taskID, userID, req.ToDomain())
	if err != nil {
		if errors.Is(err, types.ErrTaskNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
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

	if err := h.taskService.DeleteTask(c.Request.Context(), taskID, userID); err != nil {
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
