package models

import "github.com/aikidoaikido115/New-Acis-BE/modules/entities"

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
	ResidentID   string   `json:"resident_id"`
	FirstName    string   `json:"first_name"`
	LastName     string   `json:"last_name"`
	Nickname     *string  `json:"nickname"`
	RoomNumber   string   `json:"room_number"`
	IntakeLabels []string `json:"intake_labels"`
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

type VitalSignsOverviewResponse struct {
	Items      []*entities.VitalSign `json:"items"`
	Pagination OverviewPagination    `json:"pagination"`
}

type LaboratoryValuesOverviewResponse struct {
	Items      []*entities.LaboratoryValue `json:"items"`
	Pagination OverviewPagination          `json:"pagination"`
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
