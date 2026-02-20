package config

import (
	"fmt"
	"log"
	"os"
	"sound-horee-backend/models"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

// ConnectDatabase initializes the database connection pool using GORM.
func ConnectDatabase() {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASS"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_NAME"),
	)

	database, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Fatal: Could not connect to the database. ", err)
	}

	// Auto-migrate schema to match Go structs
	// [PERBAIKAN]: Menambahkan Payment dan NotificationRule agar tabelnya dibuat otomatis
	err = database.AutoMigrate(
		&models.Profile{},
		&models.Transaction{},
		&models.Payment{},
		&models.NotificationRule{},
	)
	if err != nil {
		log.Fatal("Fatal: Database migration failed. ", err)
	}

	DB = database
	log.Println("Info: Database connection established successfully.")
}
