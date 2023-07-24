package main

import (

	// middleware "smt-go-server/middleware"
	// routes "smt-go-server/routes"

	"os"

	"github.com/gin-gonic/gin"
	_ "github.com/heroku/x/hmetrics/onload"
	"github.com/octocat0415/routes"
)

func main() {
	port := os.Getenv("PORT")

	if port == "" {
		port = "8080"
	}

	gin.SetMode(gin.ReleaseMode)
	router := gin.New()
	router.Use(gin.Logger())

	routes.UserRoutes(router)
	routes.TicketRoutes(router)
	routes.EventRoutes(router)
	routes.SeatRoutes(router)

	// Test API
	router.GET("/api/test", func(c *gin.Context) {
		c.JSON(200, gin.H{"success": "Success"})

	})

	// router.Use(middleware.Authentication())

	// API - 1
	router.GET("/api-2", func(c *gin.Context) {
		c.JSON(200, gin.H{"success": "Access granted for api-2"})
	})

	router.Run(":" + port)
}
