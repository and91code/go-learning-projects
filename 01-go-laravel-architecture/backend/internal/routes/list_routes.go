package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/seu-usuario/taskflow-backend/internal/handler"
)

// TaskRoutes registra os endpoints do módulo de tarefas.
func ListRoutes(h *handler.ListHandler) RouteRegistrar {
	return func(api *gin.RouterGroup) {
		list := api.Group("/list")
		list.GET("", h.Index)
		list.POST("", h.Store)
		list.GET("/:id", h.Show)
	}
}
