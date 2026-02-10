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
