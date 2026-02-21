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

	// --- TAMBAHAN LOGIC BYPASS TESTER ---
	// Jika email adalah tester, kita paksa UID-nya agar selalu masuk ke akun yang sama
	if input.Email == "tester@algoritmakitadigital.id" {
		input.UID = "REVIEWER-GOOGLE-PLAY-001"
		input.StoreName = "Toko Tester Google"
		input.Category = "Digital"
	}
	// ------------------------------------

	var user models.Profile
	result := config.DB.Where("uid = ?", input.UID).First(&user)

	if result.RowsAffected == 0 {
		// Jika tester belum ada di DB, dia akan otomatis terbuat di sini
		user = models.Profile{
			UID:         input.UID,
			Email:       input.Email,
			StoreName:   input.StoreName,
			PhoneNumber: input.PhoneNumber,
			Category:    input.Category,
			JoinedAt:    utils.NowMillis(),
			IsPremium:   true, // Kasih akses premium langsung buat reviewer!
		}
		config.DB.Create(&user)
	} else {
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
			config.DB.First(&user, "uid = ?", input.UID)
		}
	}

	// --- LOGIC TRIAL (PENAMBAHAN) ---
	const trialLimit = 10
	var trialUsage int64
	config.DB.Model(&models.Transaction{}).Where("user_id = ?", user.UID).Count(&trialUsage)

	// Config jatah trial dari backend
	remainingTrial := int(trialLimit) - int(trialUsage)
	if remainingTrial < 0 {
		remainingTrial = 0
	}

	token, err := utils.GenerateToken(user.UID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal generate token"})
		return
	}

	// Response JSON dengan field Subscription baru
	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"token":  token,
		"user":   user,
		"subscription": gin.H{
			// PERBAIKAN: Langsung masukkan variabelnya karena sudah bool
			"is_premium":      user.IsPremium,
			"trial_limit":     trialLimit,
			"trial_usage":     trialUsage,
			"remaining_trial": remainingTrial,
		},
	})
}
