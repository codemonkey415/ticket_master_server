package main

import (

	// middleware "smt-go-server/middleware"
	// routes "smt-go-server/routes"

	"os"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	_ "github.com/heroku/x/hmetrics/onload"
	"github.com/octocat0415/middleware"
	"github.com/octocat0415/routes"
)

func main() {
	port := os.Getenv("PORT")

	if port == "" {
		port = "8080"
	}

	// gin.SetMode(gin.ReleaseMode)
	router := gin.New()

	config := cors.DefaultConfig()
	config.AllowHeaders = []string{"X-Requested-With", "Content-Type", "Authorization"}
	config.AllowMethods = []string{"GET", "POST", "PUT", "DELETE"}
	config.AllowOrigins = []string{"*"}
	router.Use(cors.New(config))

	router.Use(gin.Logger())
	
	// Test API
	router.GET("/api/test/", func(c *gin.Context) {
		c.JSON(200, gin.H{"success": "Success"})
	})

	// Test API
	router.GET("/api/test/", func(c *gin.Context) {
		c.JSON(200, gin.H{"success": "Success"})
	})

	routes.AuthRoutes(router)
	router.Use(middleware.Authentication())
	routes.UserRoutes(router)
	routes.TicketRoutes(router)
	routes.EventRoutes(router)
	routes.SeatRoutes(router)

	// API - 1
	router.GET("/api-2/", func(c *gin.Context) {
		c.JSON(200, gin.H{"success": "Access granted for api-2"})
	})

	router.Run(":" + port)
}
