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
	config.AllowAllOrigins = true
	config.AllowHeaders = []string{"X-Requested-With", "Content-Type", "Authorization"}
	config.AllowMethods = []string{"GET", "POST", "PUT", "DELETE"}
	router.Use(cors.New(config))

	router.Use(gin.Logger())

	// Test API
	router.GET("/test", func(c *gin.Context) {
		c.JSON(200, gin.H{"success": "Success"})
	})

	routes.AuthRoutes(router)
	router.Use(middleware.Authentication())
	routes.UserRoutes(router)
	routes.TicketRoutes(router)
	routes.EventRoutes(router)
	routes.SeatRoutes(router)

	router.Run(":" + port)
}
