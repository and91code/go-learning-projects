package routes

import "github.com/gin-gonic/gin"

// ProductHandler descreve os handlers necessários para as rotas de produtos.
type ProductHandler interface {
	Index(*gin.Context)
	Store(*gin.Context)
	Show(*gin.Context)
	Update(*gin.Context)
	Destroy(*gin.Context)
}

// ProductRoutes registra os endpoints de produtos protegidos pelo middleware recebido.
func ProductRoutes(productHandler ProductHandler, authMiddleware gin.HandlerFunc) RouteRegistrar {
	return func(api *gin.RouterGroup) {
		products := api.Group("/products", authMiddleware)
		products.GET("", productHandler.Index)
		products.POST("", productHandler.Store)
		products.GET("/:id", productHandler.Show)
		products.PUT("/:id", productHandler.Update)
		products.DELETE("/:id", productHandler.Destroy)
	}
}
