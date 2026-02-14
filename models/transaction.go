package models

type Transaction struct {
	ID uint `gorm:"primaryKey" json:"id"`
	// Tambahkan tag type:varchar(255) agar GORM tidak mengubahnya jadi longtext
	UserID         string  `gorm:"type:varchar(255);index" json:"user_id"`
	SourceApp      string  `json:"source_app"`
	Amount         float64 `json:"amount"`
	RawMessage     string  `gorm:"type:text" json:"raw_message"`
	Timestamp      int64   `json:"timestamp"`
	IsTrialLimited bool    `json:"is_trial_limited"`
	CreatedAt      int64   `json:"created_at"`
}
