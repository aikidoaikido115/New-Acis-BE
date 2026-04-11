package models

type CreateDrugMasterRequest struct {
	Name string `json:"name" binding:"required"`
	Dose string `json:"dose" binding:"required"`
}

type UpdateDrugMasterRequest struct {
	Name *string `json:"name"`
	Dose *string `json:"dose"`
}

type CreatePersonalDrugRequest struct {
	ResidentID  string  `json:"resident_id" binding:"required"`
	DmID        string  `json:"dm_id" binding:"required"`
	Amount      string  `json:"amount" binding:"required"`
	AmountUnit  string  `json:"amount_unit" binding:"required"`
	Frequency   int     `json:"frequency" binding:"required"`
	TimeOfDay   string  `json:"time_of_day" binding:"required"`
	Timing      string  `json:"timing" binding:"required"`
	Description *string `json:"description"`
	TakeType    string  `json:"take_type" binding:"required"`
	StartDate   *string `json:"start_date"`
	EndDate     *string `json:"end_date"`
}

type PersonalDrugOverviewQueryParams struct {
	TimeOfDay *string `json:"time_of_day" form:"time_of_day" query:"time_of_day"`
	TakeType  *string `json:"take_type" form:"take_type" query:"take_type"`
	Search    *string `json:"search" form:"search" query:"search"`
}

type UpdatePersonalDrugRequest struct {
	ResidentID  *string `json:"resident_id"`
	DmID        *string `json:"dm_id"`
	Amount      *string `json:"amount"`
	AmountUnit  *string `json:"amount_unit"`
	Frequency   *int    `json:"frequency"`
	TimeOfDay   *string `json:"time_of_day"`
	Timing      *string `json:"timing"`
	Description *string `json:"description"`
	TakeType    *string `json:"take_type"`
	StartDate   *string `json:"start_date"`
	EndDate     *string `json:"end_date"`
}

type CreateDrugPlanRequest struct {
	PdID           string  `json:"pd_id" binding:"required"`
	IsTaken        bool    `json:"is_taken"`
	TakenAt        *string `json:"taken_at"`
	GivenByStaffID string  `json:"given_by_staff_id" binding:"required"`
	IsOmmitted     *bool   `json:"is_omitted"`
	OmmittedReason *string `json:"omitted_reason"`
	Notes          *string `json:"notes"`
}

type UpdateDrugPlanRequest struct {
	PdID           *string `json:"pd_id"`
	IsTaken        *bool   `json:"is_taken"`
	TakenAt        *string `json:"taken_at"`
	GivenByStaffID *string `json:"given_by_staff_id"`
	IsOmmitted     *bool   `json:"is_omitted"`
	OmmittedReason *string `json:"omitted_reason"`
	Notes          *string `json:"notes"`
}

type DrugPlanOverviewQueryParams struct {
	TimeOfDay *string `json:"time_of_day" form:"time_of_day" query:"time_of_day"`
	TakeType  *string `json:"take_type" form:"take_type" query:"take_type"`
	Search    *string `json:"search" form:"search" query:"search"`
}
