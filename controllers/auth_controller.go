package controllers

import (
	"fmt"
	"net/http"
	"sound-horee-backend/config"
	"sound-horee-backend/models"
	"sound-horee-backend/utils"
	"time"

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

	// --- 1. LOGIC BYPASS TESTER ---
	if input.Email == "tester@algoritmakitadigital.id" {
		input.UID = "REVIEWER-GOOGLE-PLAY-001"
		input.StoreName = "Toko Tester Google"
		input.Category = "Digital"
	}

	var user models.Profile
	result := config.DB.Where("uid = ?", input.UID).First(&user)

	now := time.Now().UnixMilli()

	if result.RowsAffected > 0 {
		// --- 2. LOGIKA CEK EXPIRED (UNTUK USER LAMA) ---
		// Cek apakah waktu server sekarang sudah melewati batas premium_expires_at
		if user.IsPremium && user.PremiumExpiresAt > 0 && now > user.PremiumExpiresAt {
			user.IsPremium = false // Update di objek memory

			// Update status di Database agar permanen jadi 0 (Free)
			config.DB.Model(&user).Update("is_premium", false)

			fmt.Printf("DEBUG: User %s masa premium habis. Status diupdate ke Free.\n", user.Email)
		}

		// Update data profil jika ada perubahan dari client
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
			config.DB.Model(&user).Updates(updates)
			// Refresh data user setelah update
			config.DB.First(&user, "uid = ?", input.UID)
		}

	} else {
		// --- 3. LOGIKA REGISTER (UNTUK USER BARU) ---
		isUserTester := (input.Email == "tester@algoritmakitadigital.id")

		user = models.Profile{
			UID:         input.UID,
			Email:       input.Email,
			StoreName:   input.StoreName,
			PhoneNumber: input.PhoneNumber,
			Category:    input.Category,
			JoinedAt:    utils.NowMillis(),
			IsPremium:   isUserTester, // Reviewer langsung Premium, user biasa False
		}
		config.DB.Create(&user)
	}

	// --- 4. LOGIC TRIAL USAGE ---
	const trialLimit = 10
	var trialUsage int64
	config.DB.Model(&models.Transaction{}).Where("user_id = ?", user.UID).Count(&trialUsage)

	remainingTrial := int(trialLimit) - int(trialUsage)
	if remainingTrial < 0 {
		remainingTrial = 0
	}

	// --- 5. GENERATE JWT TOKEN ---
	token, err := utils.GenerateToken(user.UID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal generate token"})
		return
	}

	// --- 6. RESPONSE FINAL ---
	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"token":  token,
		"user":   user,
		"subscription": gin.H{
			"is_premium":      user.IsPremium,
			"trial_limit":     trialLimit,
			"trial_usage":     trialUsage,
			"remaining_trial": remainingTrial,
		},
	})
}
