package controllers

import (
	"eventplanner/database"
	"eventplanner/models"
	"net/http"

	"github.com/gin-gonic/gin"
)


// 🟢 Create new event (organizer only)
func CreateEvent(c *gin.Context) {
	var event models.Event

	if err := c.ShouldBindJSON(&event); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if event.OrganizerID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Organizer ID required"})
		return
	}

	// Save event
	if err := database.DB.Create(&event).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Add organizer to invitations
	organizerInvite := models.Invitation{
		EventID: event.ID,
		UserID:  event.OrganizerID,
		Role:    "organizer",
		Status:  "going", // organizer is automatically going
	}
	database.DB.Create(&organizerInvite)

	c.JSON(http.StatusOK, gin.H{
		"message": "Event created successfully",
		"event":   event,
	})
}



// 🟢 Get all events organized by the user
func GetOrganizedEvents(c *gin.Context) {
	userID := c.Param("userId")

	var events []models.Event
	if err := database.DB.Where("organizer_id = ?", userID).Find(&events).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"organized_events": events})
}



// 🟢 Get all events where user is invited
func GetInvitedEvents(c *gin.Context) {
	userID := c.Param("userId")

	var invites []models.Invitation
	if err := database.DB.Where("user_id = ?", userID).Find(&invites).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var eventIDs []uint
	for _, inv := range invites {
		eventIDs = append(eventIDs, inv.EventID)
	}

	var events []models.Event
	if len(eventIDs) > 0 {
		database.DB.Where("id IN ?", eventIDs).Find(&events)
	}

	c.JSON(http.StatusOK, gin.H{"invited_events": events})
}



	// 🟢 Delete event
	func DeleteEvent(c *gin.Context) {
	id := c.Param("id")

	// Delete invitations first
	database.DB.Where("event_id = ?", id).Delete(&models.Invitation{})

	// Delete event
	if err := database.DB.Delete(&models.Event{}, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Event deleted successfully"})
}




// 🟢 NEW — Update RSVP Status (Going, Maybe, Not Going)
func UpdateRSVP(c *gin.Context) {
	invID := c.Param("id")

	var body struct {
		Status string `json:"status"`
	}

	if err := c.BindJSON(&body); err != nil {
		c.JSON(400, gin.H{"error": "Invalid request"})
		return
	}

	allowed := map[string]bool{
		"going":     true,
		"maybe":     true,
		"not_going": true,
	}

	if !allowed[body.Status] {
		c.JSON(400, gin.H{"error": "Invalid status"})
		return
	}

	var invitation models.Invitation
	if err := database.DB.First(&invitation, invID).Error; err != nil {
		c.JSON(404, gin.H{"error": "Invitation not found"})
		return
	}

	invitation.Status = body.Status
	database.DB.Save(&invitation)

	c.JSON(200, gin.H{"message": "RSVP updated", "invitation": invitation})
}



// 🟢 NEW — Get attendee list + statuses for an event
func GetEventAttendees(c *gin.Context) {
	eventID := c.Param("id")

	var attendees []models.Invitation
	if err := database.DB.Preload("User").Where("event_id = ?", eventID).Find(&attendees).Error; err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, attendees)
}




// 🟢 NEW — Search + filter events
func SearchEvents(c *gin.Context) {
	keyword := c.Query("keyword")
	date := c.Query("date")
	role := c.Query("role") // organizer or attendee
	userID := c.Query("userId")

	var events []models.Event

	query := database.DB.Model(&models.Event{})

	// Keyword search (name or description)
	if keyword != "" {
		query = query.Where("name LIKE ? OR description LIKE ?", "%"+keyword+"%", "%"+keyword+"%")
	}

	// Filter by date YYYY-MM-DD
	if date != "" {
		query = query.Where("DATE(date) = ?", date)
	}

	// Filter by role
	if role != "" && userID != "" {
		if role == "organizer" {
			query = query.Where("organizer_id = ?", userID)
		} else {
			var invites []models.Invitation
			database.DB.Where("user_id = ?", userID).Find(&invites)

			var ids []uint
			for _, i := range invites {
				ids = append(ids, i.EventID)
			}
			if len(ids) > 0 {
				query = query.Where("id IN ?", ids)
			}
		}
	}

	query.Find(&events)
	c.JSON(200, events)
}
