package controllers

import (
	"eventplanner/database"
	"eventplanner/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

// 🟢 Invite a user to an event by email (attendee only)
func InviteUser(c *gin.Context) {
	var body struct {
		EventID uint   `json:"event_id"`
		Email   string `json:"email"`
	}

	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Find user by email
	var user models.User
	if err := database.DB.Where("email = ?", body.Email).First(&user).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User with this email does not exist"})
		return
	}

	// Prevent duplicate invitations
	var existing models.Invitation
	if err := database.DB.Where("event_id = ? AND user_id = ?", body.EventID, user.ID).First(&existing).Error; err == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User already invited"})
		return
	}

	invite := models.Invitation{
		EventID: body.EventID,
		UserID:  user.ID,
		Role:    "attendee",
	}

	if err := database.DB.Create(&invite).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":      "User invited successfully",
		"invited_user": user.Email,
	})
}
