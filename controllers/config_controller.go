package controllers

import (
	"net/http"
	"sound-horee-backend/config"
	"sound-horee-backend/models"

	"github.com/gin-gonic/gin"
)

// GetNotificationRules mengambil daftar regex dan tts format aktif
func GetNotificationRules(c *gin.Context) {
	var rules []models.NotificationRule

	// Tarik data dari DB yang statusnya aktif
	if err := config.DB.Where("is_active = ?", true).Find(&rules).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "Gagal mengambil konfigurasi notifikasi",
		})
		return
	}

	// Kembalikan langsung sebagai Array JSON agar Android gampang parsingnya
	c.JSON(http.StatusOK, rules)
}
