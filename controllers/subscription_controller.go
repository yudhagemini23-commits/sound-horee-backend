package controllers

import (
	"net/http"
	"sound-horee-backend/config"
	"sound-horee-backend/models"
	"time"

	"github.com/gin-gonic/gin"
)

type UpgradeRequest struct {
	UserID           string `json:"user_id"` // Jika masih dikirim dari client
	PlanType         string `json:"plan_type" binding:"required"`
	IAPPurchaseToken string `json:"iap_purchase_token"` // TAMBAHKAN INI
	IAPOrderID       string `json:"iap_order_id"`       // TAMBAHKAN INI
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

	// --- PROTEKSI 1: CEK APAKAH TOKEN SUDAH PERNAH DIPAKAI ---
	var existingPayment models.Payment
	checkToken := config.DB.Where("iap_purchase_token = ?", req.IAPPurchaseToken).First(&existingPayment)
	if checkToken.RowsAffected > 0 {
		c.JSON(http.StatusConflict, gin.H{"error": "Token pembelian ini sudah pernah digunakan"})
		return
	}

	var duration time.Duration
	var price float64
	var dbPlanType string

	switch req.PlanType {
	case "monthly", "bulanan":
		duration = 30 * 24 * time.Hour
		price = 49000
		dbPlanType = "monthly"
	default:
		c.JSON(http.StatusBadRequest, gin.H{"error": "Paket tidak dikenali"})
		return
	}

	now := time.Now().UnixMilli()
	expiryDate := now + duration.Milliseconds()

	tx := config.DB.Begin()

	// 2. Update Profile User
	if err := tx.Model(&models.Profile{}).Where("uid = ?", userID).Updates(map[string]interface{}{
		"is_premium":         1,
		"premium_expires_at": expiryDate,
	}).Error; err != nil {
		tx.Rollback()
		c.JSON(500, gin.H{"error": "Gagal update profile"})
		return
	}

	// 3. Simpan Log Pembayaran
	paymentLog := models.Payment{
		UserID:           userID,
		PlanType:         dbPlanType,
		Amount:           price,
		Status:           "success",
		IapPurchaseToken: req.IAPPurchaseToken,
		IapOrderID:       req.IAPOrderID,
		CreatedAt:        now,
		UpdatedAt:        now, // Tambahkan ini agar tidak nol
	}

	if err := tx.Create(&paymentLog).Error; err != nil {
		tx.Rollback()
		c.JSON(500, gin.H{"error": "Gagal mencatat audit pembayaran"})
		return
	}

	tx.Commit()
	c.JSON(200, gin.H{"status": "success", "expiry_date": expiryDate})
}
