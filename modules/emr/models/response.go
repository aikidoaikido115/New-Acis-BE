package models

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