package entities

import "time"

type LaboratoryValue struct {
	ID               string    `json:"laboratory_value_id" gorm:"primaryKey"`
	ResidentID       string    `json:"resident_id" gorm:"primaryKey" `
	BloodGlucose     *float64  `json:"blood_glucose"`
	FluidIn          *float64  `json:"fluid_in"`
	FluidOut         *float64  `json:"fluid_out"`
	UrineOutput      *float64  `json:"urine_output"`
	UrineType        *string   `json:"urine_type"`
	Stool            *int16   `json:"stool"`
	DiaperChange     *int16     `json:"diaper_change"`
	CreatedByStaffID string    `json:"created_by_staff_id" gorm:"type:text"`
	CreatedAt        time.Time `json:"created_at" gorm:"not null"`
	UpdatedAt        time.Time `json:"updated_at" gorm:"not null"`

	Resident Resident `json:"-" gorm:"foreignKey:ResidentID;references:ID"`
}
