package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/seu-usuario/taskflow-backend/internal/handler"
)

func TagRoutes(h *handler.TagHandler) RouteRegistrar {
	return func(api *gin.RouterGroup) {
		tags := api.Group("/tags")
		tags.GET("", h.IndexByList)
		tags.POST("", h.Store)
		tags.GET("/:id", h.Show)
		tags.PUT("/:id", h.Update)
		tags.DELETE("/:id", h.Delete)
	}
}
