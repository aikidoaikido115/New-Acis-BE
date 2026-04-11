package models

type DrugPlanResidentSummaryResponse struct {
	TotalResidents   int64 `json:"total_residents"`
	GivenResidents   int64 `json:"given_residents"`
	WaitingResidents int64 `json:"waiting_residents"`
}
