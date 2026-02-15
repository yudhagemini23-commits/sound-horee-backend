package controllers

import (
	"net/http"
	"sound-horee-backend/config"
	"sound-horee-backend/models"
	"time"

	"github.com/gin-gonic/gin"
)

type UpgradeRequest struct {
	// UserID dihapus dari sini karena kita ambil dari Token (Middleware)
	PlanType string `json:"plan_type" binding:"required"`
}

func UpgradeToPremium(c *gin.Context) {
	userID := c.GetString("user_id")
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User identity not found"})
		return
	}

	var req UpgradeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	var duration time.Duration
	var price float64
	var dbPlanType string // Variabel bantuan untuk MySQL

	// STANDARISASI INPUT
	switch req.PlanType {
	case "weekly", "mingguan":
		duration = 7 * 24 * time.Hour
		price = 5000
		dbPlanType = "weekly" // Sesuai ENUM MySQL
	case "monthly", "bulanan":
		duration = 30 * 24 * time.Hour
		price = 10000
		dbPlanType = "monthly" // Sesuai ENUM MySQL
	default:
		c.JSON(http.StatusBadRequest, gin.H{"error": "Package not recognized"})
		return
	}

	now := time.Now().UnixMilli()
	expiryDate := now + duration.Milliseconds()

	tx := config.DB.Begin()

	// 1. Update Profile
	if err := tx.Model(&models.Profile{}).Where("uid = ?", userID).Updates(map[string]interface{}{
		"is_premium":         1,
		"premium_expires_at": expiryDate,
	}).Error; err != nil {
		tx.Rollback()
		c.JSON(500, gin.H{"error": "Failed to update profile"})
		return
	}

	// 2. Log Payment (Gunakan dbPlanType agar tidak error ENUM)
	paymentLog := models.Payment{
		UserID:    userID,
		PlanType:  dbPlanType, // Pakai "weekly" atau "monthly"
		Amount:    price,
		Status:    "success",
		CreatedAt: now,
	}

	if err := tx.Create(&paymentLog).Error; err != nil {
		tx.Rollback()
		// Error 500 Mas tadi berasal dari baris ini
		c.JSON(500, gin.H{"error": "Failed to log payment audit: " + err.Error()})
		return
	}

	tx.Commit()

	c.JSON(200, gin.H{"status": "success", "expiry_date": expiryDate})
}
