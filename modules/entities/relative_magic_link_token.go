package entities

import "time"

// RelativeMagicLinkToken stores one-time/login tokens for relative portal links.
type RelativeMagicLinkToken struct {
	ID              string     `json:"magic_link_token_id" gorm:"primaryKey"`
	RelativeID      string     `json:"relative_id" gorm:"not null;index"`
	ResidentID      string     `json:"resident_id" gorm:"not null;index"`
	Token           string     `json:"token" gorm:"not null;uniqueIndex"`
	ExpiresAt       time.Time  `json:"expires_at" gorm:"not null;index"`
	LastAccessedAt  *time.Time `json:"last_accessed_at"`
	CreatedByUserID string     `json:"created_by_user_id" gorm:"not null;index"`
	CreatedAt       time.Time  `json:"created_at" gorm:"not null"`
	UpdatedAt       time.Time  `json:"updated_at" gorm:"not null"`

	Relative *Relative `json:"-" gorm:"foreignKey:RelativeID;references:ID"`
	Resident *Resident `json:"-" gorm:"foreignKey:ResidentID;references:ID"`
}
