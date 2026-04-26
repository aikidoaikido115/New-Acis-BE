package entities

import "time"

type RelativeNote struct {
	ID                 string    `json:"relative_note_id" gorm:"primaryKey"`
	ResidentID         string    `json:"resident_id" gorm:"not null"`
	Relation           string    `json:"relation" gorm:"not null"`
	Content            string    `json:"content" gorm:"not null"`
	SendNote           bool      `json:"send_note" gorm:"not null;default:true"`
	CreatedByStaffID   string    `json:"created_by_staff_id" gorm:"type:text;not null"`
	CreatedByStaffName string    `json:"created_by_staff_name,omitempty" gorm:"-"`
	CreatedAt          time.Time `json:"created_at" gorm:"not null"`
	UpdatedAt          time.Time `json:"updated_at" gorm:"not null"`

	Resident Resident `json:"-" gorm:"foreignKey:ResidentID;references:ID"`
}
