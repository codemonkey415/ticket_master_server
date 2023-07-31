package routes

import (
	controller "github.com/octocat0415/controllers"

	"github.com/gin-gonic/gin"
)

// UserRoutes function

func UserRoutes(incomingRoutes *gin.Engine) {
	incomingRoutes.GET("/users/getUser/", controller.GetUser())
	incomingRoutes.GET("/api/users", controller.GetUsers())
	incomingRoutes.POST("/api/users/changeStatus/:userid", controller.ChangeStatus())
	incomingRoutes.POST("/api/users/updateDueDate/:userid/:month", controller.UpdateDueDate())
	incomingRoutes.POST("/api/users/changeRole/:userid", controller.ChangeRole())
	incomingRoutes.POST("/api/users/saveReservations/:userid", controller.SaveReservations())
	incomingRoutes.POST("/api/users/removeReservations/:userid", controller.RemoveReservations())
}

func AuthRoutes(incomingRoutes *gin.Engine) {
	incomingRoutes.POST("/users/register/", controller.SignUp())
	incomingRoutes.POST("/users/login/", controller.Login())
}

func TicketRoutes(incomingRoutes *gin.Engine) {
	incomingRoutes.POST("/api/ticket/", controller.GetAllTickets())
}

func EventRoutes(incomingRoutes *gin.Engine) {
	incomingRoutes.GET("/api/events/", controller.GetAllEvents())
}

func SeatRoutes(incomingRoutes *gin.Engine) {
	incomingRoutes.GET("/api/seat/events", controller.GetEvents())
	incomingRoutes.GET("/api/seat/sections/:eventid", controller.GetSectionNames())
	incomingRoutes.GET("/api/seat/rows/:eventid", controller.GetRows())
	incomingRoutes.POST("/api/seat/", controller.GetAllTickets())
	incomingRoutes.POST("/api/seat/notify", controller.NotifySeat())
}
