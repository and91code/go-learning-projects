// Package routes compõe os grupos de rotas da API usando Gin.
package routes

import (
	"database/sql"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/seu-usuario/taskflow-backend/internal/logger"
	"github.com/seu-usuario/taskflow-backend/internal/module"
)

// RouteRegistrar registra as rotas de um módulo no grupo versionado da API.
type RouteRegistrar func(*gin.RouterGroup)

// SetupRouter cria o engine Gin, injeta as dependências e registra as rotas.
func SetupRouter(db *sql.DB) *gin.Engine {
	router := gin.New()
	router.Use(
		logger.GinMiddleware(),
		gin.Recovery(),
		corsMiddleware(),
		jsonContentTypeMiddleware(),
	)
	router.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "pong"})
	})

	api := router.Group("/api/v1")
	TaskRoutes(module.BuildTaskModule(db))(api)

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
