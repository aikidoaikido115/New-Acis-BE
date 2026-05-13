package entities

import (
	"time"

	"gorm.io/datatypes"
)

type Resident struct {
	ID        string  `json:"resident_id" gorm:"primaryKey" `
	RoomID    *string `json:"room_id"`
	FirstName string  `json:"first_name" gorm:"not null"`
	LastName  string  `json:"last_name" gorm:"not null"`
	Gender    string  `json:"gender" gorm:"not null"`

	Nickname                   *string        `json:"nickname"`
	IdCardNumber               *string        `json:"id_card_number" gorm:"unique"`
	DateOfBirth                time.Time      `json:"date_of_birth" gorm:"not null;type:date"`
	PurposeOfStay              *string        `json:"purpose_of_stay"`
	CheckInDate                *time.Time     `json:"check_in_date" gorm:"type:date"`
	ExpectedCheckOutDate       *time.Time     `json:"expected_check_out_date" gorm:"type:date"`
	Status                     string         `json:"status" gorm:"not null"`
	PreExistingConditions      *string        `json:"pre_existing_conditions"`
	PreExistingConditionsNotes *string        `json:"pre_existing_conditions_notes"`
	ResucitationStatus         *string        `json:"resuscitation_status"`
	SugicalHistory             *string        `json:"surgical_history"`
	PreferredEmergencyHospital *string        `json:"preferred_emergency_hospital"`
	EmergencyHospitalPhone     *string        `json:"emergency_hospital_phone"`
	EmergencyHospitals         datatypes.JSON `json:"emergency_hospitals" gorm:"type:jsonb"`
	EmergencyContacts          datatypes.JSON `json:"emergency_contacts" gorm:"type:jsonb" swaggertype:"object"`
	Room                       Room           `json:"-" gorm:"foreignKey:RoomID;references:ID"`
	ProfileImage               *string        `json:"profile_image"`

	ResidentLabels    []ResidentLabels    `json:"resident_labels" gorm:"foreignKey:ResidentID;references:ID"`
	ResidentAllergies []ResidentAllergies `json:"resident_allergies" gorm:"foreignKey:ResidentID;references:ID"`
	ResidentDA        []ResidentDA        `json:"resident_das" gorm:"foreignKey:ResidentID;references:ID"`
}
