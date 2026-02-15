package models

// Payment mencatat riwayat pembelian untuk monitoring/audit
type Payment struct {
	ID               uint    `gorm:"primaryKey" json:"id"`
	UserID           string  `gorm:"type:varchar(255);index" json:"user_id"`
	PlanType         string  `gorm:"type:varchar(50)" json:"plan_type"` // weekly / monthly
	Amount           float64 `json:"amount"`
	Status           string  `gorm:"type:varchar(50)" json:"status"` // success / pending
	IapPurchaseToken string  `gorm:"type:text" json:"iap_purchase_token"`
	CreatedAt        int64   `json:"created_at"`
}
