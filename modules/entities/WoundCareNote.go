package entities

import "time"

type WoundCareNote struct {
	ID               string    `json:"wound_care_note_id" gorm:"primaryKey"`
	ResidentID       string    `json:"resident_id" gorm:"not null"`
	Location         string    `json:"location" gorm:"not null"`
	WoundType        string    `json:"wound_type" gorm:"not null"`
	Size             *string   `json:"size"`
	Treatment        *string   `json:"treatment"`
	Supplies         *string   `json:"supplies"`
	Status           *string   `json:"status"`
	ImageURL         *string   `json:"image_url"`
	Note             *string   `json:"note"`
	CreatedByStaffID string    `json:"created_by_staff_id" gorm:"type:text;not null"`
	CreatedAt        time.Time `json:"created_at" gorm:"not null"`
	UpdatedAt        time.Time `json:"updated_at" gorm:"not null"`

	Resident Resident `json:"-" gorm:"foreignKey:ResidentID;references:ID"`
}
