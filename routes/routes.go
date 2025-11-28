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

	// 🔵 RSVP (Update status)
	router.PUT("/rsvp/:id", controllers.UpdateRSVP)

	// 🔵 Get attendee list w/ statuses
	router.GET("/events/:id/attendees", controllers.GetEventAttendees)

	// 🔵 Search events
	router.GET("/search/events", controllers.SearchEvents)

	// Invitation route
	router.POST("/invite", controllers.InviteUser)
}
