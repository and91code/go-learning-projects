package routes

import "github.com/gin-gonic/gin"

// ListHandler descreve os handlers necessários para registrar as rotas de tarefas.
type ListHandler interface {
	Index(*gin.Context)
	Store(*gin.Context)
	Show(*gin.Context)
}

// TaskRoutes registra os endpoints do módulo de tarefas.
func ListRoutes(h ListHandler) RouteRegistrar {
	return func(api *gin.RouterGroup) {
		list := api.Group("/list")
		list.GET("", h.Index)
		list.POST("", h.Store)
	}
}
