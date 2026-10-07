package routes

import "github.com/gin-gonic/gin"

// TaskHandler descreve os handlers necessários para registrar as rotas de tarefas.
type TaskHandler interface {
	Store(*gin.Context)
}

// TaskRoutes registra os endpoints do módulo de tarefas.
func TaskRoutes(h TaskHandler) RouteRegistrar {
	return func(api *gin.RouterGroup) {
		tasks := api.Group("/tasks")
		tasks.POST("", h.Store)
	}
}
