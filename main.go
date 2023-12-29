package main

import (
	config "github.com/KhetwalDevesh/restaurant-management/database"
	"github.com/KhetwalDevesh/restaurant-management/middleware"
	"github.com/KhetwalDevesh/restaurant-management/routes"
	"github.com/gin-gonic/gin"
	"os"
)

func main() {
	port := os.Getenv("SERVER_PORT")
	if port == "" {
		port = "8000"
	}
	config.ConfigDB()
	router := gin.New()
	router.Use(gin.Logger())
	routes.UserRoutes(router)
	router.Use(middleware.Authentication())
	routes.FoodRoutes(router)
	routes.InvoiceRoutes(router)
	routes.MenuRoutes(router)
	routes.OrderRoutes(router)
	routes.OrderItemRoutes(router)
	routes.TableRoutes(router)
	router.Run(":" + port)
}
