package entities

import "time"

type NurseNote struct {
	ID                 string    `json:"nurse_note_id" gorm:"primaryKey"`
	ResidentID         string    `json:"resident_id" gorm:"not null"`
	Category           string    `json:"category" gorm:"not null"`
	Content            string    `json:"content" gorm:"not null"`
	Priority           string    `json:"priority" gorm:"not null"`
	SendNote           bool      `json:"send_note" gorm:"not null;default:false"`
	CreatedByStaffID   string    `json:"created_by_staff_id" gorm:"type:text;not null"`
	CreatedByStaffName string    `json:"created_by_staff_name,omitempty" gorm:"-"`
	CreatedAt          time.Time `json:"created_at" gorm:"not null"`
	UpdatedAt          time.Time `json:"updated_at" gorm:"not null"`

	Resident Resident `json:"-" gorm:"foreignKey:ResidentID;references:ID"`
}
