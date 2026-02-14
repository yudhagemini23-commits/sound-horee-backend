package routes

import (
	"sound-horee-backend/controllers"
	"sound-horee-backend/middlewares" // Import middleware

	"github.com/gin-gonic/gin"
)

func SetupRoutes(r *gin.Engine) {
	v1 := r.Group("/api/v1")
	{
		// --- PUBLIC ROUTES (Bisa diakses siapa saja) ---
		auth := v1.Group("/auth")
		{
			auth.POST("/login", controllers.LoginOrRegister) // Ini endpoint LOGIN
		}

		// --- PROTECTED ROUTES (Harus pakai Token) ---
		// Kita pasang middleware di sini
		protected := v1.Group("/", middlewares.AuthRequired())
		{
			// Profile
			protected.POST("/profile/sync", controllers.SyncProfile)
			protected.GET("/profile/:uid", controllers.GetProfile)

			// Transactions
			protected.POST("/transactions/sync", controllers.SyncTransactions)
			protected.GET("/transactions", controllers.GetTransactions)
		}
	}
}
