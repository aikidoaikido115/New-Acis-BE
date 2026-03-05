package models

import "time"

type CreateResidentRequest struct {
	RoomID    string `json:"room_id" binding:"required"`
	FirstName string `json:"first_name" binding:"required"`
	LastName  string `json:"last_name" binding:"required"`
	Gender    string `json:"gender" binding:"required"`

	Nickname                   *string    `json:"nickname"`
	IdCardNumber               string     `json:"id_card_number" binding:"required"`
	DateOfBirth                time.Time  `json:"date_of_birth" binding:"required"`
	PurposeOfStay              *string    `json:"purpose_of_stay"`
	CheckInDate                time.Time  `json:"check_in_date" binding:"required"`
	ExpectedCheckOutDate       *time.Time `json:"expected_check_out_date"`
	Status                     string     `json:"status" binding:"required"`
	PreExistingConditions      *string    `json:"pre_existing_conditions"`
	PreExistingConditionsNotes *string    `json:"pre_existing_conditions_notes"`
	ResucitationStatus         *string    `json:"resuscitation_status"`
	SugicalHistory             *string    `json:"surgical_history"`
	PreferredEmergencyHospital *string    `json:"preferred_emergency_hospital"`
	EmergencyHospitalPhone     *string    `json:"emergency_hospital_phone"`
}

type IntakeLabelRequest struct {
	LabelName string  `json:"label_name" binding:"required"`
	NoteText  *string `json:"note_text"`
}

type CreateIntakeLabelByResidentRequest struct {
	ResidentID string               `json:"resident_id" binding:"required"`
	Labels     []IntakeLabelRequest `json:"labels" binding:"required"`
}

type UpdateResidentRequest struct {
	RoomID      *string    `json:"room_id"`
	FirstName   *string    `json:"first_name"`
	LastName    *string    `json:"last_name"`
	DateOfBirth *time.Time `json:"date_of_birth"`
	Gender      *string    `json:"gender"`

	Nickname                   *string    `json:"nickname"`
	IdCardNumber               *string    `json:"id_card_number"`
	PurposeOfStay              *string    `json:"purpose_of_stay"`
	CheckInDate                *time.Time `json:"check_in_date"`
	ExpectedCheckOutDate       *time.Time `json:"expected_check_out_date"`
	Status                     *string    `json:"status"`
	PreExistingConditions      *string    `json:"pre_existing_conditions"`
	PreExistingConditionsNotes *string    `json:"pre_existing_conditions_notes"`
	ResucitationStatus         *string    `json:"resuscitation_status"`
	SugicalHistory             *string    `json:"surgical_history"`
	PreferredEmergencyHospital *string    `json:"preferred_emergency_hospital"`
	EmergencyHospitalPhone     *string    `json:"emergency_hospital_phone"`

	Labels []IntakeLabelRequest `json:"labels,omitempty"` // ใช้ logic เดิม: label_name มีอยู่ → ใช้อันเดิม, ไม่มี → สร้างใหม่
}

type UpdateRoomRequest struct {
	StaffID *string `json:"staff_id"`
}

type CreateRoomRequest struct {
	StaffID    *string `json:"staff_id" binding:"required"`
	Floor      int16   `json:"floor" binding:"required"`
	RoomNumber string  `json:"room_number" binding:"required"`
}

type CreateVitalSignRequest struct {
	ResidentID             string   `json:"resident_id" binding:"required"`
	Temperature            *float64 `json:"temperature"`
	HeartRate              *int16   `json:"heart_rate"`
	BreathingRate          *int16   `json:"breathing_rate"`
	BloodPressureSystolic  *int16   `json:"blood_pressure_systolic"`
	BloodPressureDiastolic *int16   `json:"blood_pressure_diastolic"`
	OxygenSaturation       *int16   `json:"oxygen_saturation"`
}

type VitalSignQueryParams struct {
    ResidentID      *string    `json:"resident_id" form:"resident_id" query:"resident_id"`
    RoomID          *string    `json:"room_id" form:"room_id" query:"room_id"`
    Floor           *int16     `json:"floor" form:"floor" query:"floor"`// nil = ทุกชั้น, มีค่า = ชั้นนั้น
    LabelIDs        []string   `json:"label_ids" form:"label_ids" query:"label_ids"`// empty = ทุกกลุ่ม, มีค่า = กรองตาม label
    StartDate       *time.Time `json:"start_date" form:"start_date" query:"start_date"`
    EndDate         *time.Time `json:"end_date" form:"end_date" query:"end_date"`
    IsLatest        bool       `json:"is_latest" form:"is_latest" query:"is_latest"`// true = Latest (DISTINCT ON resident_id), false = All
    Limit           int        `json:"limit" form:"limit" query:"limit"`// default 100
    Offset          int        `json:"offset" form:"offset" query:"offset"`// default 0
    VitalSignStatus string     `json:"vitalsign_status" form:"vitalsign_status" query:"vitalsign_status"` // เพิ่ม query tag ตรงนี้
}

type UpdateVitalSignRequest struct {
	Temperature            *float64 `json:"temperature"`
	HeartRate              *int16   `json:"heart_rate"`
	BreathingRate          *int16   `json:"breathing_rate"`
	BloodPressureSystolic  *int16   `json:"blood_pressure_systolic"`
	BloodPressureDiastolic *int16   `json:"blood_pressure_diastolic"`
	OxygenSaturation       *int16   `json:"oxygen_saturation"`
}
