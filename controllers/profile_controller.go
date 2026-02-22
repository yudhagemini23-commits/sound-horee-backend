package controllers

import (
	"fmt"
	"net/http"
	"sound-horee-backend/config"
	"sound-horee-backend/models"
	"time" // Hanya menambah ini

	"github.com/gin-gonic/gin"
)

// SyncProfile handles user registration and profile updates.
func SyncProfile(c *gin.Context) {
	var input models.Profile
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var existing models.Profile
	result := config.DB.First(&existing, "uid = ?", input.UID)

	if result.RowsAffected == 0 {
		config.DB.Create(&input)
	} else {
		// --- TAMBAHAN LOGIK CHECK EXPIRED ---
		now := time.Now().UnixMilli()
		if existing.IsPremium && existing.PremiumExpiresAt > 0 && now > existing.PremiumExpiresAt {
			existing.IsPremium = false
			config.DB.Model(&existing).Update("is_premium", false)
			fmt.Printf("Silent Check (Sync): %s expired\n", existing.Email)
		}
		// ------------------------------------

		config.DB.Model(&existing).Updates(models.Profile{
			Email:       input.Email,
			StoreName:   input.StoreName,
			PhoneNumber: input.PhoneNumber,
			Category:    input.Category,
		})
	}

	c.JSON(http.StatusOK, gin.H{"status": "success", "data": existing})
}

// GetProfile fetches profile details, including subscription status.
func GetProfile(c *gin.Context) {
	uid := c.Param("uid")
	var profile models.Profile

	if err := config.DB.First(&profile, "uid = ?", uid).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	// --- TAMBAHAN LOGIK CHECK EXPIRED ---
	now := time.Now().UnixMilli()
	if profile.IsPremium && profile.PremiumExpiresAt > 0 && now > profile.PremiumExpiresAt {
		profile.IsPremium = false
		config.DB.Model(&profile).Update("is_premium", false)
		fmt.Printf("Silent Check (GetProfile): %s expired\n", profile.Email)
	}
	// ------------------------------------

	c.JSON(http.StatusOK, gin.H{"status": "success", "data": profile})
}
