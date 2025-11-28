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

	// Save event to DB
	if err := database.DB.Create(&event).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Automatically mark creator as "organizer" in Invitations
	organizerInvite := models.Invitation{
		EventID: event.ID,
		UserID:  event.OrganizerID,
		Role:    "organizer",
	}
	database.DB.Create(&organizerInvite)

	c.JSON(http.StatusOK, gin.H{
		"message": "Event created successfully",
		"event":   event,
	})
}

// 🟢 Get all events organized by a specific user
func GetOrganizedEvents(c *gin.Context) {
	userID := c.Param("userId")

	var events []models.Event
	if err := database.DB.Where("organizer_id = ?", userID).Find(&events).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"organized_events": events})
}

// 🟢 Get all events where user is invited (attendee)
func GetInvitedEvents(c *gin.Context) {
	userID := c.Param("userId")

	var invitations []models.Invitation
	if err := database.DB.Where("user_id = ?", userID).Find(&invitations).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var eventIDs []uint
	for _, inv := range invitations {
		eventIDs = append(eventIDs, inv.EventID)
	}

	var events []models.Event
	if len(eventIDs) > 0 {
		if err := database.DB.Where("id IN ?", eventIDs).Find(&events).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{"invited_events": events})
}

// 🟢 Delete event (organizer only)
func DeleteEvent(c *gin.Context) {
	id := c.Param("id")

	// Delete related invitations first
	database.DB.Where("event_id = ?", id).Delete(&models.Invitation{})

	// Delete the event itself
	if err := database.DB.Delete(&models.Event{}, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Event deleted successfully"})
}
