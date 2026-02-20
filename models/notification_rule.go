package models

import "time"

// NotificationRule merepresentasikan tabel notification_rules di database
type NotificationRule struct {
	ID           uint      `gorm:"primaryKey" json:"id"`
	PackageName  string    `gorm:"unique;not null" json:"package_name"`
	AppName      string    `gorm:"not null" json:"app_name"`
	RegexPattern string    `gorm:"not null" json:"regex_pattern"`
	TtsFormat    string    `gorm:"not null" json:"tts_format"`
	IsActive     bool      `gorm:"default:true" json:"is_active"`
	UpdatedAt    time.Time `json:"updated_at"`
}
