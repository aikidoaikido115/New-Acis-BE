package entities

import "time"

// DailyUpdate เก็บบันทึกการอัปเดตประจำวันที่จะแจ้งให้ญาติทราบ
type DailyUpdate struct {
	ID               string    `json:"daily_id" gorm:"primaryKey"`
	RelativeID       string    `json:"relative_id" gorm:"not null;index"`
	Description      string    `json:"description" gorm:"type:text;not null"`
	CreatedByStaffID string    `json:"created_by_staff_id" gorm:"not null;index"`
	CreatedAt        time.Time `json:"created_at" gorm:"not null"`
	UpdatedAt        time.Time `json:"updated_at" gorm:"not null"`

	// Relationships
	Relative   *Relative   `json:"-" gorm:"foreignKey:RelativeID;references:ID"`
	DailyFiles []DailyFile `json:"daily_files" gorm:"foreignKey:DailyID;references:ID"`
}
