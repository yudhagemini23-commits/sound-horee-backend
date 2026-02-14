package controllers

import (
	"net/http"
	"sound-horee-backend/config"
	"sound-horee-backend/models"
	"sound-horee-backend/utils"

	"github.com/gin-gonic/gin"
)

type LoginInput struct {
	UID         string `json:"uid"`
	Email       string `json:"email"`
	StoreName   string `json:"store_name"`
	PhoneNumber string `json:"phone_number"`
	Category    string `json:"category"`
}

func LoginOrRegister(c *gin.Context) {
	var input LoginInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var user models.Profile
	// 1. Cari berdasarkan UID
	result := config.DB.Where("uid = ?", input.UID).First(&user)

	if result.RowsAffected == 0 {
		// --- REGISTER (User Benar-benar Baru) ---
		user = models.Profile{
			UID:         input.UID,
			Email:       input.Email,
			StoreName:   input.StoreName,
			PhoneNumber: input.PhoneNumber,
			Category:    input.Category,
			JoinedAt:    utils.NowMillis(),
		}
		if err := config.DB.Create(&user).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal membuat user"})
			return
		}
	} else {
		// --- LOGIN (User Lama) ---
		// Hanya update jika input tidak kosong (Silent Check logic)
		updates := make(map[string]interface{})
		if input.Email != "" {
			updates["email"] = input.Email
		}
		if input.StoreName != "" {
			updates["store_name"] = input.StoreName
		}
		if input.PhoneNumber != "" {
			updates["phone_number"] = input.PhoneNumber
		}
		if input.Category != "" {
			updates["category"] = input.Category
		}

		if len(updates) > 0 {
			// Update ke database
			config.DB.Model(&user).Updates(updates)
			// Refresh data 'user' di memori agar sinkron dengan yang ada di DB
			config.DB.First(&user, "uid = ?", input.UID)
		}
	}

	// 2. Generate JWT Token
	token, err := utils.GenerateToken(user.UID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal generate token"})
		return
	}

	// 3. Response JSON (Sangat Krusial untuk Android)
	// Pastikan field "user" berisi object Profile lengkap
	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"token":  token,
		"user":   user,
	})
}
