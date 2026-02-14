package controllers

import (
	"net/http"
	"sound-horee-backend/config"
	"sound-horee-backend/models"
	"sound-horee-backend/utils"

	"github.com/gin-gonic/gin"
)

type LoginRequest struct {
	UID         string `json:"uid" binding:"required"`
	Email       string `json:"email" binding:"required"`
	StoreName   string `json:"store_name"`
	PhoneNumber string `json:"phone_number"`
}

func LoginOrRegister(c *gin.Context) {
	var input LoginRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 1. Cek apakah user sudah ada di database?
	var user models.Profile
	result := config.DB.First(&user, "uid = ?", input.UID)

	if result.RowsAffected == 0 {
		// --- REGISTER (User Baru) ---
		user = models.Profile{
			UID:         input.UID,
			Email:       input.Email,
			StoreName:   input.StoreName, // Default dari nama Google
			PhoneNumber: input.PhoneNumber,
			JoinedAt:    utils.NowMillis(), // Helper timestamp (opsional) atau time.Now().UnixMilli()
			IsPremium:   false,
		}
		config.DB.Create(&user)
	} else {
		// --- LOGIN (User Lama) ---
		// Opsional: Update email kalau berubah di Google
		config.DB.Model(&user).Update("email", input.Email)
	}

	// 2. Generate JWT Token
	token, err := utils.GenerateToken(user.UID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	// 3. Kirim Token + Data User ke Android
	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"token":  token,
		"user":   user,
	})
}
