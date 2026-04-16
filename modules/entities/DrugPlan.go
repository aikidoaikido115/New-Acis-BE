package entities

import "time"

type DrugPlan struct {
	ID             string     `json:"dpln_id" gorm:"primaryKey"`
	PdID           string     `json:"pd_id" gorm:"not null"`
	IsTaken        bool       `json:"is_taken" gorm:"not null"`
	TakenAt        *time.Time `json:"taken_at,omitempty" gorm:"type:timestamptz"`
	GivenByStaffID string     `json:"given_by_staff_id" gorm:"not null"`
	IsOmitted      *bool      `json:"is_omitted" gorm:"not null"`
	OmittedReason  *string    `json:"omitted_reason"`
	Notes          *string    `json:"notes"`
	CreatedAt      time.Time  `json:"created_at" gorm:"not null"`
	UpdatedAt      time.Time  `json:"updated_at" gorm:"not null"`

	PersonalDrug PersonalDrug `gorm:"foreignKey:PdID;references:ID"`
}
