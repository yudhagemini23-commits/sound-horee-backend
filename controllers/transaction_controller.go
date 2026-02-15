package controllers

import (
	"net/http"
	"sound-horee-backend/config"
	"sound-horee-backend/models"
	"strconv"

	"github.com/gin-gonic/gin"
)

// SyncTransactions handles batch upload of offline transactions.
func SyncTransactions(c *gin.Context) {
	var transactions []models.Transaction
	if err := c.ShouldBindJSON(&transactions); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Perform batch insert for performance optimization
	if len(transactions) > 0 {
		result := config.DB.CreateInBatches(&transactions, 100)
		if result.Error != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"status":       "success",
		"synced_count": len(transactions),
	})
}

// GetTransactions retrieves transaction history with flexible filters.
// Supports filtering by: UserID (Required), Date Range (Start/End timestamps), Source App.
func GetTransactions(c *gin.Context) {
	userID := c.Query("user_id")
	if userID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "user_id parameter is required"})
		return
	}

	query := config.DB.Where("user_id = ?", userID)

	startStr := c.Query("start")
	endStr := c.Query("end")

	if startStr != "" && endStr != "" {
		start, _ := strconv.ParseInt(startStr, 10, 64)
		end, _ := strconv.ParseInt(endStr, 10, 64)
		query = query.Where("timestamp BETWEEN ? AND ?", start, end)
	}

	if sourceApp := c.Query("source_app"); sourceApp != "" {
		query = query.Where("source_app LIKE ?", "%"+sourceApp+"%")
	}

	var transactions []models.Transaction
	if err := query.Order("timestamp desc").Find(&transactions).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// --- PERBAIKAN DISINI ---
	// Langsung kirim []transactions agar Android bisa parsing 'as List<TransactionDto>'
	// Jika Mas tetap ingin kirim status & total_amount, Android-nya yang harus dirombak total.
	// Karena Mas minta jangan nyenggol sana-sini, kita samakan Go dengan ekspektasi Android.
	c.JSON(http.StatusOK, transactions)
}
