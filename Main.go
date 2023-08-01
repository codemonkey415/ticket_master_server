package main

import (
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/octocat0415/middleware"
	"github.com/octocat0415/routes"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()
	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	router.Use(middleware.CORSMiddleware())

	router.GET("/test", func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization")
		c.JSON(http.StatusOK, gin.H{"success": "Success"})
	})
	routes.AuthRoutes(router)

	router.Use(middleware.Authentication())

	routes.UserRoutes(router)
	routes.TicketRoutes(router)
	routes.EventRoutes(router)
	routes.SeatRoutes(router)

	server := &http.Server{
		Addr:         ":" + port,
		Handler:      router,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  15 * time.Second,
	}
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		panic(err)
	}
}
