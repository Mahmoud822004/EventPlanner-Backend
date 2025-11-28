package routes

import (
	"eventplanner/controllers"
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(router *gin.Engine) {

	// Event routes
	router.POST("/events", controllers.CreateEvent)
	router.GET("/events/organized/:userId", controllers.GetOrganizedEvents)
	router.GET("/events/invited/:userId", controllers.GetInvitedEvents)
	router.DELETE("/events/:id", controllers.DeleteEvent)

	// Invitation route
	router.POST("/invite", controllers.InviteUser)
}
