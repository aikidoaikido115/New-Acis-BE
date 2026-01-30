package entities

import "time"

type ResidentLabels struct {
	ResidentID string    `json:"resident_id" gorm:"primaryKey" `
	LabelID    string    `json:"label_id" gorm:"primaryKey" `
	NoteText   string    `json:"note_text" gorm:"not null"`
	NotedAt    time.Time `json:"noted_at" gorm:"not null"`

	Resident    Resident     `json:"-" gorm:"foreignKey:ResidentID;references:ID"`
	IntakeLabel IntakeLabels `json:"-" gorm:"foreignKey:LabelID;references:ID"`
}
