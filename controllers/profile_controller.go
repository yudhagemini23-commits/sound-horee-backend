package controllers

import (
	"net/http"
	"sound-horee-backend/config"
	"sound-horee-backend/models"

	"github.com/gin-gonic/gin"
)

// SyncProfile handles user registration and profile updates.
// Strategy: Upsert (Insert if new, Update if exists).
func SyncProfile(c *gin.Context) {
	var input models.Profile
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var existing models.Profile
	result := config.DB.First(&existing, "uid = ?", input.UID)

	if result.RowsAffected == 0 {
		// Create new profile
		config.DB.Create(&input)
	} else {
		// Update existing profile details (excluding subscription status)
		config.DB.Model(&existing).Updates(models.Profile{
			Email:       input.Email,
			StoreName:   input.StoreName,
			PhoneNumber: input.PhoneNumber,
			Category:    input.Category,
		})
	}

	c.JSON(http.StatusOK, gin.H{"status": "success", "data": input})
}

// GetProfile fetches profile details, including subscription status.
func GetProfile(c *gin.Context) {
	uid := c.Param("uid")
	var profile models.Profile

	if err := config.DB.First(&profile, "uid = ?", uid).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "success", "data": profile})
}
