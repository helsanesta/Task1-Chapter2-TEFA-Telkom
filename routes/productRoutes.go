package routes

import (
	controller "task1-ch2-api-product/controllers"
	middleware "task1-ch2-api-product/middleware"

	"github.com/gin-gonic/gin"
)

func ProductRoutes(incomingRoutes *gin.Engine) {
	incomingRoutes.GET("/products", controller.GetProducts())
	incomingRoutes.GET("/products/:id", controller.GetProduct())
	incomingRoutes.POST("/products", middleware.Authentication(), controller.CreateProduct())
	incomingRoutes.PUT("/products/:id", middleware.Authentication(), controller.UpdateProduct())
	incomingRoutes.DELETE("/products/:id", middleware.Authentication(), controller.DeleteProduct())
	incomingRoutes.GET("products/search", controller.SearchProduct())
}
