package models

// Transaction represents financial data synced from the Android app.
type Transaction struct {
	ID             uint    `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID         string  `gorm:"type:varchar(255);index" json:"user_id"` // Indexed for faster filtering
	SourceApp      string  `gorm:"type:varchar(255)" json:"source_app"`
	Amount         float64 `json:"amount"`
	RawMessage     string  `gorm:"type:text" json:"raw_message"`
	Timestamp      int64   `json:"timestamp"` // Transaction time (Epoch Millis)
	IsTrialLimited bool    `json:"is_trial_limited"`
	CreatedAt      int64   `gorm:"autoCreateTime" json:"created_at"`
}
