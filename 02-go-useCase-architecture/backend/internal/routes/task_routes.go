package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/seu-usuario/taskflow-backend/internal/handler"
)

// TaskRoutes registra os endpoints CRUD do módulo de tarefas.
func TaskRoutes(h *handler.TaskHandler) RouteRegistrar {
	return func(api *gin.RouterGroup) {
		tasks := api.Group("/tasks")
		tasks.GET("", h.IndexByList)
		tasks.POST("", h.Store)
		tasks.GET("/:id", h.Show)
		tasks.PUT("/:id", h.Update)
		tasks.DELETE("/:id", h.Delete)
	}
}
