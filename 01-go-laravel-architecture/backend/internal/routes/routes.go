// Package routes compõe os grupos de rotas da API usando Gin.
package routes

import (
	"database/sql"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/seu-usuario/taskflow-backend/internal/logger"
	"github.com/seu-usuario/taskflow-backend/internal/middleware"
	"github.com/seu-usuario/taskflow-backend/internal/module"
)

// RouteRegistrar registra as rotas de um módulo no grupo versionado da API.
type RouteRegistrar func(*gin.RouterGroup)

// SetupRouter cria o engine Gin, injeta as dependências e registra as rotas.
func SetupRouter(db *sql.DB) *gin.Engine {
	// Ocultar lista de rotas registradas
	gin.SetMode(gin.ReleaseMode)

	// 1. Inicializa o motor do Gin
	router := gin.New()

	// 2. Middlewares Globais (CORS, Recovery, etc.)
	router.Use(
		logger.GinMiddleware(),
		gin.Recovery(),
		corsMiddleware(),
		jsonContentTypeMiddleware(),
	)

	api := router.Group("/api/v1")
	modules := module.Build(db)
	{
		api.GET("/ping", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"status": "ok", "message": "pong"})
		})

		// Rotas Públicas (Sem necessidade de Token JWT)
		AuthRoutes(modules.Auth)(api)

		// Rotas Protegidas (Exigem o Token JWT)
		apiProtectedGroup := api.Group("")
		apiProtectedGroup.Use(middleware.AuthMiddleware())
		{
			ListRoutes(modules.Lists)(apiProtectedGroup)
			TaskRoutes(modules.Tasks)(apiProtectedGroup)
		}
	}

	return router
}

// corsMiddleware permite chamadas HTTP de origens diferentes à API.
func corsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Authorization, Content-Type")
		if c.Request.Method == http.MethodOptions {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}
		c.Next()
	}
}

// jsonContentTypeMiddleware define application/json nas respostas da API.
func jsonContentTypeMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Content-Type", "application/json")
		c.Next()
	}
}
