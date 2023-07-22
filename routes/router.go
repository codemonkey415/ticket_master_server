package routes

import (
	controller "github.com/octocat0415/controllers"

	"github.com/gin-gonic/gin"
)

// UserRoutes function
func UserRoutes(incomingRoutes *gin.Engine) {
	incomingRoutes.POST("/users/signup", controller.SignUp())
	incomingRoutes.POST("/users/login", controller.Login())
}

func TicketRoutes(incomingRoutes *gin.Engine) {
	incomingRoutes.POST("/api/ticket", controller.GetAllTickets())
}

func EventRoutes(incomingRoutes *gin.Engine) {
	incomingRoutes.GET("/api/events", controller.GetAllEvents())
}

func SeatRoutes(incomingRoutes *gin.Engine) {
	incomingRoutes.GET("/api/seat/sections", controller.GetSectionNames())
	incomingRoutes.GET("/api/seat/rows", controller.GetRows())
	incomingRoutes.POST("/api/seat", controller.GetAllTickets())
	incomingRoutes.POST("/api/seat/notify", controller.NotifySeat())
}
