package entities

import "time"

type ResidentAllergies struct {
	ResidentID string    `json:"resident_id" gorm:"primaryKey" `
	AllergyID    string    `json:"allergy_id" gorm:"primaryKey" `
	NoteText   *string   `json:"note_text" gorm:"type:text"` // ใช้ pointer เพื่อให้สามารถเช็ค null ได้
	NotedAt    time.Time `json:"noted_at" gorm:"not null"`

	Resident    Resident     `json:"-" gorm:"foreignKey:ResidentID;references:ID"`
	Allergy Allergy `json:"allergy,omitempty" gorm:"foreignKey:AllergyID;references:ID"`
}
