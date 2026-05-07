package entities

import "time"

type ResidentLabels struct {
	ResidentID string    `json:"resident_id" gorm:"primaryKey" `
	LabelID    string    `json:"label_id" gorm:"primaryKey" `
	NoteText   *string   `json:"note_text" gorm:"type:text"` // ใช้ pointer เพื่อให้สามารถเช็ค null ได้
	NotedAt    time.Time `json:"noted_at" gorm:"not null"`

	Resident    Resident     `json:"-" gorm:"foreignKey:ResidentID;references:ID"`
	IntakeLabel IntakeLabels `json:"intake_label,omitempty" gorm:"foreignKey:LabelID;references:ID"`
}
