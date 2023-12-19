package main

import (
	"github.com/KhetwalDevesh/restaurant-management/middleware"
	"github.com/KhetwalDevesh/restaurant-management/routes"
	"github.com/gin-gonic/gin"
	"os"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8000"
	}
	router := gin.New()
	router.Use(gin.Logger())
	routes.UserRoutes(router)
	routes.Use(middleware.Authentication())
	routes.FoodRoutes(router)
	routes.InvoiceRoutes(router)
	routes.MenuRoutes(router)
	routes.OrderRoutes(router)
	routes.OrderItemRoutes(router)
	routes.TableRoutes(router)
	router.Run(":" + port)
}
