package models

// Profile represents the user/store entity.
type Profile struct {
	UID         string `gorm:"primaryKey;type:varchar(255)" json:"uid"`
	Email       string `json:"email"`
	StoreName   string `json:"store_name"`
	PhoneNumber string `json:"phone_number"`
	Category    string `json:"category"`
	JoinedAt    int64  `json:"joined_at"`

	// Subscription Status
	IsPremium        bool  `gorm:"default:false" json:"is_premium"`
	PremiumExpiresAt int64 `gorm:"default:0" json:"premium_expires_at"`

	// Relations
	Transactions []Transaction `gorm:"foreignKey:UserID;references:UID" json:"transactions,omitempty"`
}
