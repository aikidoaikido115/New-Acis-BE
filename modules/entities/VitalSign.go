package entities

import "time"

type VitalSign struct {
	ID                     string    `json:"vital_sign_id" gorm:"primaryKey"`
	ResidentID             string    `json:"resident_id" gorm:"primaryKey" `
	Temperature            *float64  `json:"temperature"`
	HeartRate              *int16    `json:"heart_rate"`
	BreathingRate          *int16    `json:"breathing_rate"`
	BloodPressureSystolic  *int16    `json:"blood_pressure_systolic"`
	BloodPressureDiastolic *int16    `json:"blood_pressure_diastolic"`
	OxygenSaturation       *int16    `json:"oxygen_saturation"`
	CreatedByStaffID       string    `json:"created_by_staff_id" gorm:"type:text"`
	CreatedAt              time.Time `json:"created_at" gorm:"not null"`
	UpdatedAt              time.Time `json:"updated_at" gorm:"not null"`

	Resident Resident `json:"-" gorm:"foreignKey:ResidentID;references:ID"`
}
