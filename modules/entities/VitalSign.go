package entities

import "time"

type VitalSign struct {
	ID                     string    `json:"vital_sign_id" gorm:"primaryKey"`
	ResidentID             string    `json:"resident_id" gorm:"primaryKey" `
	Temperature            *float64  `json:"temperature" gorm:"type:text"`
	HeartRate              *int16    `json:"heart_rate" gorm:"type:text"`
	BreathingRate          *int16    `json:"breathing_rate" gorm:"type:text"`
	BloodPressureSystolic  *int16    `json:"blood_pressure_systolic" gorm:"type:text"`
	BloodPressureDiastolic *int16    `json:"blood_pressure_diastolic" gorm:"type:text"`
	OxygenSaturation       *int16    `json:"oxygen_saturation" gorm:"type:text"`
	CreatedAt              time.Time `json:"created_at" gorm:"not null"`
	UpdatedAt              time.Time `json:"updated_at" gorm:"not null"`

	Resident Resident `json:"-" gorm:"foreignKey:ResidentID;references:ID"`
}
