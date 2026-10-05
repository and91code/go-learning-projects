package routes

import "github.com/gin-gonic/gin"

// AuthHandler descreve os handlers necessários para autenticação.
type AuthHandler interface {
	Register(*gin.Context)
	Login(*gin.Context)
}

// AuthRoutes registra os endpoints públicos de autenticação.
func AuthRoutes(authHandler AuthHandler) RouteRegistrar {
	return func(api *gin.RouterGroup) {
		auth := api.Group("/auth")
		auth.POST("/register", authHandler.Register)
		auth.POST("/login", authHandler.Login)
	}
}
