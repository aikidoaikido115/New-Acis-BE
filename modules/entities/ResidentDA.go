package entities

import "time"

type ResidentDA struct {
	ResidentID string    `json:"resident_id" gorm:"primaryKey" `
	DrugAllergyID    string    `json:"drug_allergy_id" gorm:"primaryKey" `
	NoteText   *string   `json:"note_text" gorm:"type:text"` // ใช้ pointer เพื่อให้สามารถเช็ค null ได้
	NotedAt    time.Time `json:"noted_at" gorm:"not null"`

	Resident    Resident     `json:"-" gorm:"foreignKey:ResidentID;references:ID"`
	DrugAllergy DrugAllergy `json:"drug_allergy,omitempty" gorm:"foreignKey:DrugAllergyID;references:ID"`
}
