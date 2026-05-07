package models

import (
	"time"

	"github.com/aikidoaikido115/New-Acis-BE/modules/entities"
)

type CreateResidentRequest struct {
	RoomID    *string `json:"room_id"`
	FirstName string  `json:"first_name" binding:"required"`
	LastName  string  `json:"last_name" binding:"required"`
	Gender    string  `json:"gender" binding:"required"`

	Nickname                   *string            `json:"nickname"`
	IdCardNumber               *string            `json:"id_card_number"`
	DateOfBirth                time.Time          `json:"date_of_birth" binding:"required"`
	PurposeOfStay              *string            `json:"purpose_of_stay"`
	CheckInDate                *time.Time         `json:"check_in_date"`
	ExpectedCheckOutDate       *time.Time         `json:"expected_check_out_date"`
	Status                     string             `json:"status" binding:"required"`
	PreExistingConditions      *string            `json:"pre_existing_conditions"`
	PreExistingConditionsNotes *string            `json:"pre_existing_conditions_notes"`
	ResucitationStatus         *string            `json:"resuscitation_status" binding:"required"`
	SugicalHistory             *string            `json:"surgical_history"`
	PreferredEmergencyHospital *string            `json:"preferred_emergency_hospital" binding:"required"`
	EmergencyHospitalPhone     *string            `json:"emergency_hospital_phone" binding:"required"`
	ProfileImage               *string            `json:"profile_image"`
	EmergencyContacts          []EmergencyContact `json:"emergency_contacts" binding:"required"`
}

type EmergencyContact struct {
	Name     string `json:"name"`
	Relation string `json:"relation"`
	Phone    string `json:"phone"`
}

type IntakeLabelRequest struct {
	LabelName string  `json:"label_name" binding:"required"`
	NoteText  *string `json:"note_text"`
}

type AllergyRequest struct {
	AllergyName string  `json:"allergy_name" binding:"required"`
	NoteText    *string `json:"note_text"`
}

type DrugAllergyRequest struct {
	AllergyName string  `json:"allergy_name" binding:"required"`
	NoteText    *string `json:"note_text"`
}

type CreateIntakeLabelByResidentRequest struct {
	ResidentID string               `json:"resident_id" binding:"required"`
	Labels     []IntakeLabelRequest `json:"labels" binding:"required"`
}

type CreateAllergyByResidentRequest struct {
	ResidentID string           `json:"resident_id" binding:"required"`
	Allergies  []AllergyRequest `json:"allergies" binding:"required"`
}

type CreateDrugAllergyByResidentRequest struct {
	ResidentID    string               `json:"resident_id" binding:"required"`
	DrugAllergies []DrugAllergyRequest `json:"drug_allergies" binding:"required"`
}

type UpdateResidentRequest struct {
	RoomID      *string    `json:"room_id"`
	FirstName   *string    `json:"first_name"`
	LastName    *string    `json:"last_name"`
	DateOfBirth *time.Time `json:"date_of_birth"`
	Gender      *string    `json:"gender"`

	Nickname                   *string             `json:"nickname"`
	IdCardNumber               *string             `json:"id_card_number"`
	PurposeOfStay              *string             `json:"purpose_of_stay"`
	CheckInDate                *time.Time          `json:"check_in_date"`
	ExpectedCheckOutDate       *time.Time          `json:"expected_check_out_date"`
	Status                     *string             `json:"status"`
	PreExistingConditions      *string             `json:"pre_existing_conditions"`
	PreExistingConditionsNotes *string             `json:"pre_existing_conditions_notes"`
	ResucitationStatus         *string             `json:"resuscitation_status"`
	SugicalHistory             *string             `json:"surgical_history"`
	PreferredEmergencyHospital *string             `json:"preferred_emergency_hospital"`
	EmergencyHospitalPhone     *string             `json:"emergency_hospital_phone"`
	ProfileImage               *string             `json:"profile_image"`
	EmergencyContacts          *[]EmergencyContact `json:"emergency_contacts"`

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
	Date                   string   `json:"date" binding:"required"`
	TimeOfDay              string   `json:"time_of_day" binding:"required"`
	Temperature            *float64 `json:"temperature"`
	HeartRate              *int16   `json:"heart_rate"`
	BreathingRate          *int16   `json:"breathing_rate"`
	BloodPressureSystolic  *int16   `json:"blood_pressure_systolic"`
	BloodPressureDiastolic *int16   `json:"blood_pressure_diastolic"`
	OxygenSaturation       *int16   `json:"oxygen_saturation"`
}

type VitalSignQueryParams struct {
	Date            *string    `json:"date" form:"date" query:"date"`
	TimeOfDay       *string    `json:"time_of_day" form:"time_of_day" query:"time_of_day"`
	ResidentID      *string    `json:"resident_id" form:"resident_id" query:"resident_id"`
	RoomID          *string    `json:"room_id" form:"room_id" query:"room_id"`
	Floor           *int16     `json:"floor" form:"floor" query:"floor"`             // nil = ทุกชั้น, มีค่า = ชั้นนั้น
	LabelIDs        []string   `json:"label_ids" form:"label_ids" query:"label_ids"` // empty = ทุกกลุ่ม, มีค่า = กรองตาม label
	StartDate       *time.Time `json:"start_date" form:"start_date" query:"start_date"`
	EndDate         *time.Time `json:"end_date" form:"end_date" query:"end_date"`
	IsLatest        bool       `json:"is_latest" form:"is_latest" query:"is_latest"` // true = Latest (DISTINCT ON resident_id), false = All
	Page            *int       `json:"page" form:"page" query:"page"`
	PageSize        *int       `json:"page_size" form:"page_size" query:"page_size"`
	Limit           int        `json:"limit" form:"limit" query:"limit"`                                  // default 100
	Offset          int        `json:"offset" form:"offset" query:"offset"`                               // default 0
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

type CreateLaboratoryValueRequest struct {
	ResidentID   string   `json:"resident_id" binding:"required"`
	Date         string   `json:"date" binding:"required"`
	TimeOfDay    string   `json:"time_of_day" binding:"required"`
	BloodGlucose *float64 `json:"blood_glucose"`
	FluidIn      *float64 `json:"fluid_in"`
	FluidOut     *float64 `json:"fluid_out"`
	UrineOutput  *float64 `json:"urine_output"`
	UrineType    *string  `json:"urine_type"` // "ml" หรือ "times"
	Stool        *int16   `json:"stool"`
	DiaperChange *int16   `json:"diaper_change"`
}

type LaboratoryValueQueryParams struct {
	ResidentID            *string    `json:"resident_id" form:"resident_id" query:"resident_id"`
	RoomID                *string    `json:"room_id" form:"room_id" query:"room_id"`
	Date                  *string    `json:"date" form:"date" query:"date"`
	TimeOfDay             *string    `json:"time_of_day" form:"time_of_day" query:"time_of_day"`
	Floor                 *int16     `json:"floor" form:"floor" query:"floor"`
	LabelIDs              []string   `json:"label_ids" form:"label_ids" query:"label_ids"`
	StartDate             *time.Time `json:"start_date" form:"start_date" query:"start_date"`
	EndDate               *time.Time `json:"end_date" form:"end_date" query:"end_date"`
	IsLatest              bool       `json:"is_latest" form:"is_latest" query:"is_latest"`
	Page                  *int       `json:"page" form:"page" query:"page"`
	PageSize              *int       `json:"page_size" form:"page_size" query:"page_size"`
	Limit                 int        `json:"limit" form:"limit" query:"limit"`
	Offset                int        `json:"offset" form:"offset" query:"offset"`
	LaboratoryValueStatus string     `json:"laboratory_value_status" form:"laboratory_value_status" query:"laboratory_value_status"`
}

type UpdateLaboratoryValueRequest struct {
	BloodGlucose *float64 `json:"blood_glucose"`
	FluidIn      *float64 `json:"fluid_in"`
	FluidOut     *float64 `json:"fluid_out"`
	UrineOutput  *float64 `json:"urine_output"`
	UrineType    *string  `json:"urine_type"`
	Stool        *int16   `json:"stool"`
	DiaperChange *int16   `json:"diaper_change"`
}

type ResidentQueryParams struct {
	Floor    *int16   `json:"floor" form:"floor" query:"floor"`             // nil = ทุกชั้น
	LabelIDs []string `json:"label_ids" form:"label_ids" query:"label_ids"` // empty = ทุกกลุ่ม
	Status   *string  `json:"status" form:"status" query:"status"`          // nil = ทั้ง active และ inactive
	Search   *string  `json:"search" form:"search" query:"search"`          // nil = ไม่กรอง, มีค่า = LIKE ชื่อ/นามสกุล/ชื่อเล่น
	Page     *int     `json:"page" form:"page" query:"page"`
	PageSize *int     `json:"page_size" form:"page_size" query:"page_size"`
	Limit    int      `json:"limit" form:"limit" query:"limit"`
	Offset   int      `json:"offset" form:"offset" query:"offset"`
}

type CreateNurseNoteRequest struct {
	ResidentID string `json:"resident_id" binding:"required"`
	Category   string `json:"category" binding:"required"`
	Content    string `json:"content" binding:"required"`
	Priority   string `json:"priority" binding:"required"`
	SendNote   bool   `json:"send_note"`
}

type UpdateNurseNoteRequest struct {
	Category *string `json:"category"`
	Content  *string `json:"content"`
	Priority *string `json:"priority"`
	SendNote *bool   `json:"send_note"`
}

type CreateWoundCareNoteRequest struct {
	ResidentID string  `json:"resident_id" binding:"required"`
	Location   string  `json:"location" binding:"required"`
	WoundType  string  `json:"wound_type" binding:"required"`
	Size       *string `json:"size"`
	Treatment  *string `json:"treatment"`
	Supplies   *string `json:"supplies"`
	Status     *string `json:"status"`
	ImageURL   *string `json:"image_url"`
	Note       *string `json:"note"`
}

type UpdateWoundCareNoteRequest struct {
	Location  *string `json:"location"`
	WoundType *string `json:"wound_type"`
	Size      *string `json:"size"`
	Treatment *string `json:"treatment"`
	Supplies  *string `json:"supplies"`
	Status    *string `json:"status"`
	ImageURL  *string `json:"image_url"`
	Note      *string `json:"note"`
}

type CreateRelativeNoteRequest struct {
	ResidentID string `json:"resident_id" binding:"required"`
	Relation   string `json:"relation"`
	Content    string `json:"content" binding:"required"`
	SendNote   bool   `json:"send_note"`
}

type UpdateRelativeNoteRequest struct {
	Relation *string `json:"relation"`
	Content  *string `json:"content"`
	SendNote *bool   `json:"send_note"`
}

type CreateDoctorOrderRequest struct {
	ResidentID string  `json:"resident_id" binding:"required"`
	OrderDate  *string `json:"order_date"`
	OrderType  *string `json:"order_type"`
	Title      string  `json:"title" binding:"required"`
	Details    *string `json:"details"`
	StartDate  *string `json:"start_date"`
	EndDate    *string `json:"end_date"`
	Frequency  *string `json:"frequency"`
	OrderedBy  *string `json:"ordered_by"`
}

type UpdateDoctorOrderRequest struct {
	OrderDate *string `json:"order_date"`
	OrderType *string `json:"order_type"`
	Title     *string `json:"title"`
	Details   *string `json:"details"`
	StartDate *string `json:"start_date"`
	EndDate   *string `json:"end_date"`
	Frequency *string `json:"frequency"`
	OrderedBy *string `json:"ordered_by"`
}

type IssueRelativeMagicLinkRequest struct {
	ResidentID string `json:"resident_id" binding:"required"`
}

type RelativePortalLoginRequest struct {
	ResidentID string `json:"resident_id"`
	Token      string `json:"token"`
	Password   string `json:"password" binding:"required"`
	Remember   bool   `json:"remember"`
}

type RelativePortalLoginResponse struct {
	Token      string `json:"token"`
	UserID     string `json:"user_id"`
	Username   string `json:"username"`
	Email      string `json:"email"`
	RoleName   string `json:"role_name"`
	ResidentID string `json:"resident_id"`
}

type RelativeMagicLinkResponse struct {
	ResidentID string `json:"resident_id"`
	RelativeID string `json:"relative_id"`
	Token      string `json:"token"`
	MagicLink  string `json:"magic_link"`
	ExpiresAt  string `json:"expires_at"`
}

type RelativeDashboardResponse struct {
	ResidentID     string                           `json:"resident_id"`
	ResidentName   string                           `json:"resident_name"`
	Date           string                           `json:"date"`
	LastUpdatedAt  *string                          `json:"last_updated_at"`
	Notes          []RelativeDashboardNote          `json:"notes"`
	Participations []RelativeDashboardParticipation `json:"participations"`
}

type RelativeDashboardNote struct {
	ID        string `json:"id"`
	Content   string `json:"content"`
	CreatedAt string `json:"created_at"`
}

type RelativeDashboardParticipation struct {
	ResidentID       string                    `json:"resident_id"`
	ASID             string                    `json:"as_id"`
	IsParticipating  bool                      `json:"is_participating"`
	ImgURLs          []entities.ImageURL       `json:"img_urls"`
	ActivitySchedule entities.ActivitySchedule `json:"activity_schedule"`
}

type RelativePatientInfoResponse struct {
	ResidentID                string                      `json:"resident_id"`
	FirstName                 string                      `json:"first_name"`
	LastName                  string                      `json:"last_name"`
	Nickname                  *string                     `json:"nickname"`
	Gender                    string                      `json:"gender"`
	DateOfBirth               string                      `json:"date_of_birth"`
	Age                       int                         `json:"age"`
	IdCardNumber              string                      `json:"id_card_number"`
	PurposeOfStay             *string                     `json:"purpose_of_stay"`
	CheckInDate               string                      `json:"check_in_date"`
	Status                    string                      `json:"status"`
	PreExistingConditions     []string                    `json:"pre_existing_conditions"`
	PreExistingConditionsNote *string                     `json:"pre_existing_conditions_note"`
	SurgicalHistory           []string                    `json:"surgical_history"`
	Medications               []RelativePatientMedication `json:"medications"`
	ResuscitationStatus       *string                     `json:"resuscitation_status"`
	FoodAllergies             []string                    `json:"food_allergies"`
	DrugAllergies             []string                    `json:"drug_allergies"`
	EmergencyHospital         *string                     `json:"emergency_hospital"`
	EmergencyHospitalPhone    *string                     `json:"emergency_hospital_phone"`
	EmergencyContacts         []EmergencyContact          `json:"emergency_contacts"`
}

type RelativePatientMedication struct {
	Name      string `json:"name"`
	Dose      string `json:"dose"`
	Frequency string `json:"frequency"`
	Notes     string `json:"notes"`
}
