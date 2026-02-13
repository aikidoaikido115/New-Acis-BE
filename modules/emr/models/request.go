package models

type CreateResidentRequest struct {
	RoomID    string `json:"room_id" binding:"required"`
	FirstName string `json:"first_name" binding:"required"`
	LastName  string `json:"last_name" binding:"required"`
	Age       *int16 `json:"age" binding:"required"`
	Gender    string `json:"gender" binding:"required"`
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
	RoomID    *string              `json:"room_id"`
	FirstName *string              `json:"first_name"`
	LastName  *string              `json:"last_name"`
	Age       *int16               `json:"age"`
	Gender    *string              `json:"gender"`
	Labels    []IntakeLabelRequest `json:"labels,omitempty"` // ใช้ logic เดิม: label_name มีอยู่ → ใช้เดิม, ไม่มี → สร้างใหม่
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
	HeartRate              *int16    `json:"heart_rate" gorm:"type:text"`
	BreathingRate          *int16    `json:"breathing_rate" gorm:"type:text"`
	BloodPressureSystolic  *int16    `json:"blood_pressure_systolic" gorm:"type:text"`
	BloodPressureDiastolic *int16    `json:"blood_pressure_diastolic" gorm:"type:text"`
	OxygenSaturation       *int16    `json:"oxygen_saturation" gorm:"type:text"`
	CreatedByStaffID       string    `json:"created_by_staff_id" gorm:"type:text"`
}