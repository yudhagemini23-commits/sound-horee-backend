package routes

import (
	"sound-horee-backend/controllers"
	"sound-horee-backend/middlewares"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(r *gin.Engine) {
	v1 := r.Group("/api/v1")
	{
		// --- PUBLIC ROUTES ---
		auth := v1.Group("/auth")
		{
			auth.POST("/login", controllers.LoginOrRegister)
		}

		// --- PROTECTED ROUTES (Wajib JWT) ---
		protected := v1.Group("/", middlewares.AuthRequired())
		{
			// Profile
			protected.POST("/profile/sync", controllers.SyncProfile)
			protected.GET("/profile/:uid", controllers.GetProfile)

			// Transactions (Data Uang Masuk)
			protected.POST("/transactions/sync", controllers.SyncTransactions)
			protected.GET("/transactions", controllers.GetTransactions)

			// --- TAMBAHAN: Subscription & Payments (Monitoring) ---
			// Ini untuk menangani upgrade premium dan monitoring IAP nantinya
			subscription := protected.Group("/subscription")
			{
				subscription.POST("/upgrade", controllers.UpgradeToPremium)
				// subscription.GET("/history", controllers.GetPaymentHistory) // Contoh kedepannya
			}
		}

		configGroup := v1.Group("/config")
		{
			// Endpoint ini sengaja dibuat public (tanpa middleware Auth)
			// agar Android bisa tarik kapan saja, atau Mas bisa pasang Auth jika mau.
			configGroup.GET("/rules", controllers.GetNotificationRules)
		}
	}
}
