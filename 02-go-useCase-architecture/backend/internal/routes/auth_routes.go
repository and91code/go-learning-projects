package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/seu-usuario/taskflow-backend/internal/handler"
)

// AuthRoutes registra os endpoints públicos de autenticação.
func AuthRoutes(h *handler.AuthHandler) RouteRegistrar {
	return func(api *gin.RouterGroup) {
		auth := api.Group("/auth")
		auth.POST("/register", h.Register)
		auth.POST("/login", h.Login)
	}
}

// UserRoutes registra endpoints autenticados de gerenciamento do usuário.
func UserRoutes(h *handler.AuthHandler) RouteRegistrar {
	return func(api *gin.RouterGroup) {
		user := api.Group("/user")
		user.PUT("/profile", h.UpdateProfile)
		user.PUT("/password", h.UpdatePassword)
	}
}
