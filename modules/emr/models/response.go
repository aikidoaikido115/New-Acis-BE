package models

import (
	"time"

	"github.com/aikidoaikido115/New-Acis-BE/modules/entities"
)

type NumberOfResidentsDashboardResponse struct {
	TotalResidents         int16 `json:"total_residents"`
	IndependentResidents   int16 `json:"independent_residents"`
	PartialAssistResidents int16 `json:"partial_assist_residents"`
	BedriddenResidents     int16 `json:"bedridden_residents"`
}

type ResidentGenderStatsDashboardResponse struct {
	SumOfMale        int16   `json:"sum_of_male"`
	SumOfFemale      int16   `json:"sum_of_female"`
	TotalResidents   int16   `json:"total_residents"`
	MalePercentage   float32 `json:"male_percentage"`
	FemalePercentage float32 `json:"female_percentage"`
}

type VitalSignDashboardSummary struct {
	CurrentNormalResidents   int64 `json:"current_normal_residents"`
	CurrentAbnormalResidents int64 `json:"current_abnormal_residents"`
	CurrentTotalResidents    int64 `json:"current_total_residents"`

	HadAbnormalTodayResidents   int64 `json:"had_abnormal_today_residents"`
	HadNormalOnlyTodayResidents int64 `json:"had_normal_only_today_residents"`
	HadTotalResidents           int64 `json:"had_total_residents"`
}

type DrugPlanTimeOfDayDashboardSummary struct {
	TimeOfDay        string `json:"time_of_day"`
	TotalResidents   int64  `json:"total_residents"`
	GivenResidents   int64  `json:"given_residents"`
	WaitingResidents int64  `json:"waiting_residents"`
	Status           string `json:"status"`
}

type UrineOutputSumResponse struct {
	ResidentID  string  `json:"resident_id"`
	TotalAmount float64 `json:"total_amount"`
}

type UrineOutputSummaryByResidentResponse struct {
	ResidentID string  `json:"resident_id"`
	TotalML    float64 `json:"total_ml"`
	TotalTimes float64 `json:"total_times"`
}

type ResidentOverviewResponse struct {
	ResidentID           string     `json:"resident_id"`
	FirstName            string     `json:"first_name"`
	LastName             string     `json:"last_name"`
	Nickname             *string    `json:"nickname"`
	RoomNumber           string     `json:"room_number"`
	IntakeLabels         []string   `json:"intake_labels"`
	Gender               string     `json:"gender"`
	Status               string     `json:"status"`
	CheckInDate          *time.Time `json:"check_in_date"`
	ExpectedCheckOutDate *time.Time `json:"expected_check_out_date"`
	Floor                *int16     `json:"floor"`
}

type OverviewPagination struct {
	Page       int `json:"page"`
	PageSize   int `json:"page_size"`
	TotalItems int `json:"total_items"`
	TotalPages int `json:"total_pages"`
}

type ResidentOverviewListResponse struct {
	Items      []*ResidentOverviewResponse `json:"items"`
	Pagination OverviewPagination          `json:"pagination"`
}

type VitalSignFieldStatus struct {
	Key        string `json:"key"`
	Label      string `json:"label"`
	IsAbnormal bool   `json:"is_abnormal"`
}

type VitalSignsOverviewItemResponse struct {
	*entities.VitalSign
	NormalList    []string               `json:"normal_list"`
	AbnormalList  []string               `json:"abnormal_list"`
	FieldStatuses []VitalSignFieldStatus `json:"field_statuses"`
}

type VitalSignsOverviewResponse struct {
	Items      []*VitalSignsOverviewItemResponse `json:"items"`
	Pagination OverviewPagination                `json:"pagination"`
}

type LaboratoryValueFieldStatus struct {
	Key        string `json:"key"`
	Label      string `json:"label"`
	IsAbnormal bool   `json:"is_abnormal"`
}

type LaboratoryValuesOverviewItemResponse struct {
	*entities.LaboratoryValue
	NormalList    []string                     `json:"normal_list"`
	AbnormalList  []string                     `json:"abnormal_list"`
	FieldStatuses []LaboratoryValueFieldStatus `json:"field_statuses"`
}

type LaboratoryValuesOverviewResponse struct {
	Items      []*LaboratoryValuesOverviewItemResponse `json:"items"`
	Pagination OverviewPagination                      `json:"pagination"`
}

type AllergyStatisticDashboardResponse struct {
	AllergyID     string `json:"allergy_id"`
	AllergyName   string `json:"allergy_name"`
	ResidentCount int64  `json:"resident_count"`
}

type ResidentAllergyStatsDashboardResponse struct {
	TotalNotAllergic int64                               `json:"total_not_allergic"`
	TotalAllergic    int64                               `json:"total_allergic"`
	AllergyDetails   []AllergyStatisticDashboardResponse `json:"allergy_details"`
}

type ResidentAllergyItemResponse struct {
	AllergyID   string  `json:"allergy_id"`
	AllergyName string  `json:"allergy_name"`
	NoteText    *string `json:"note_text"`
}

type ResidentAllergyListResponse struct {
	ResidentID string                        `json:"resident_id"`
	FirstName  string                        `json:"first_name"`
	LastName   string                        `json:"last_name"`
	Allergies  []ResidentAllergyItemResponse `json:"allergies"`
}

type DrugAllergyStatisticDashboardResponse struct {
	DrugAllergyID string `json:"drug_allergy_id"`
	AllergyName   string `json:"allergy_name"`
	Count         int64  `json:"count"`
}

type ResidentDrugAllergyStatsDashboardResponse struct {
	TotalNotDrugAllergic int64                                   `json:"total_not_drug_allergic"`
	TotalDrugAllergic    int64                                   `json:"total_drug_allergic"`
	DrugAllergyDetails   []DrugAllergyStatisticDashboardResponse `json:"drug_allergy_details"`
}

type ResidentDrugAllergyItemResponse struct {
	DrugAllergyID string  `json:"drug_allergy_id"`
	AllergyName   string  `json:"allergy_name"`
	NoteText      *string `json:"note_text"`
}

type ResidentDrugAllergyListResponse struct {
	ResidentID    string                            `json:"resident_id"`
	FirstName     string                            `json:"first_name"`
	LastName      string                            `json:"last_name"`
	DrugAllergies []ResidentDrugAllergyItemResponse `json:"drug_allergies"`
}
