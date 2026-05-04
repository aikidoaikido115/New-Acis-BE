package entities

import "time"

type SupportTicket struct {
	ID              string    `json:"support_ticket_id" gorm:"primaryKey"`
	Name            string    `json:"name" gorm:"not null"`
	Email           string    `json:"email" gorm:"not null"`
	Subject         string    `json:"subject" gorm:"not null"`
	Message         string    `json:"message" gorm:"not null"`
	Status          string    `json:"status" gorm:"not null;default:open"`
	ReporterRole    string    `json:"reporter_role" gorm:"not null"`
	CreatedByUserID string    `json:"created_by_user_id" gorm:"type:text;not null"`
	CreatedAt       time.Time `json:"created_at" gorm:"not null"`
	UpdatedAt       time.Time `json:"updated_at" gorm:"not null"`

	CreatedByUser User `json:"-" gorm:"foreignKey:CreatedByUserID;references:ID"`
}
