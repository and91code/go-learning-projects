package handler

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/seu-usuario/taskflow-backend/internal/dto"
	"github.com/seu-usuario/taskflow-backend/internal/service"
)

type TaskHandler struct {
	taskService service.TaskService
}

func NewTaskHandler(s service.TaskService) *TaskHandler {
	return &TaskHandler{taskService: s}
}

// Store cria uma nova tarefa
func (h *TaskHandler) Store(c *gin.Context) {
	var req dto.CreateTaskRequest

	// Decodifica JSON e valida as tags de entrada declaradas no DTO.
	if err := c.ShouldBindJSON(&req); err != nil {
		status := http.StatusBadRequest
		var validationErrors validator.ValidationErrors
		if errors.As(err, &validationErrors) {
			status = http.StatusUnprocessableEntity
		}
		c.JSON(status, gin.H{
			"error":   "Dados de entrada inválidos",
			"details": err.Error(),
		})
		return
	}

	// Erros de validação da regra de negócio são retornados como HTTP 422.
	res, err := h.taskService.CreateTask(c.Request.Context(), req)
	if err != nil {
		if errors.Is(err, service.ErrInvalidTask) {
			c.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Retorna HTTP 201 Created com o DTO formatado
	c.JSON(http.StatusCreated, res)
}
