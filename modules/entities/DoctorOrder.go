package entities

import "time"

type DoctorOrder struct {
	ID                 string    `json:"doctor_order_id" gorm:"primaryKey"`
	ResidentID         string    `json:"resident_id" gorm:"not null"`
	OrderDate          *string   `json:"order_date" gorm:"type:text"`
	OrderType          *string   `json:"order_type" gorm:"type:text"`
	Title              string    `json:"title" gorm:"not null"`
	Details            *string   `json:"details" gorm:"type:text"`
	StartDate          *string   `json:"start_date" gorm:"type:text"`
	EndDate            *string   `json:"end_date" gorm:"type:text"`
	Frequency          *string   `json:"frequency" gorm:"type:text"`
	OrderedBy          *string   `json:"ordered_by" gorm:"type:text"`
	CreatedByStaffID   string    `json:"created_by_staff_id" gorm:"type:text;not null"`
	CreatedByStaffName string    `json:"created_by_staff_name,omitempty" gorm:"-"`
	CreatedAt          time.Time `json:"created_at" gorm:"not null"`
	UpdatedAt          time.Time `json:"updated_at" gorm:"not null"`

	Resident Resident `json:"-" gorm:"foreignKey:ResidentID;references:ID"`
}
