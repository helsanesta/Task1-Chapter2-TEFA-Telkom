package main

import (
	"os"

	routes "task1-ch2-api-product/routes"

	"github.com/gin-gonic/gin"
	_ "github.com/heroku/x/hmetrics/onload"
)

func main() {
	port := os.Getenv("PORT")

	if port == "" {
		port = "8000"
	}

	router := gin.New()
	router.Use(gin.Logger())
	routes.UserRoutes(router)
	routes.ProductRoutes(router)

	router.Run(":" + port)
}
