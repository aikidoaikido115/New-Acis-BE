package models

import (
	"time"

	"github.com/aikidoaikido115/New-Acis-BE/modules/entities"
)

type DrugPlanResidentSummaryResponse struct {
	TotalResidents   int64 `json:"total_residents"`
	GivenResidents   int64 `json:"given_residents"`
	WaitingResidents int64 `json:"waiting_residents"`
}

type DrugPlanTimeOfDaySummary struct {
	TimeOfDay        string `json:"time_of_day" gorm:"column:time_of_day"`
	TotalResidents   int64  `json:"total_residents" gorm:"column:total_residents"`
	GivenResidents   int64  `json:"given_residents" gorm:"column:given_residents"`
	WaitingResidents int64  `json:"waiting_residents" gorm:"column:waiting_residents"`
	Status           string `json:"status,omitempty" gorm:"-"`
}

type DrugPlanGenerationResponse struct {
	GeneratedCount       int    `json:"generated_count"`
	SkippedExistingCount int    `json:"skipped_existing_count"`
	ExpiredDeletedCount  int    `json:"expired_deleted_count"`
	Scope                string `json:"scope"`
	ResidentID           string `json:"resident_id,omitempty"`
}

type DrugAdministrationHistoryItem struct {
	DrugPlanID       string     `json:"drug_plan_id" gorm:"column:drug_plan_id"`
	ActionAt         *time.Time `json:"action_at" gorm:"column:action_at"`
	ResidentName     string     `json:"resident_name" gorm:"column:resident_name"`
	DrugName         string     `json:"drug_name" gorm:"column:drug_name"`
	DrugDose         string     `json:"drug_dose" gorm:"column:drug_dose"`
	Status           string     `json:"status" gorm:"column:status"`
	Note             *string    `json:"note" gorm:"column:note"`
	GivenByStaffName *string    `json:"given_by_staff_name" gorm:"column:given_by_staff_name"`
	TimeOfDay        string     `json:"time_of_day" gorm:"column:time_of_day"`
}

type DrugAdministrationHistoryPagination struct {
	Page       int `json:"page"`
	PageSize   int `json:"page_size"`
	TotalItems int `json:"total_items"`
	TotalPages int `json:"total_pages"`
}

type DrugAdministrationHistoryResponse struct {
	Items      []DrugAdministrationHistoryItem     `json:"items"`
	Pagination DrugAdministrationHistoryPagination `json:"pagination"`
}

type PersonalDrugOverviewResponse struct {
	Items      []*entities.PersonalDrug            `json:"items"`
	Pagination DrugAdministrationHistoryPagination `json:"pagination"`
}

type DrugPlanOverviewResponse struct {
	Items      []*entities.DrugPlan                `json:"items"`
	Pagination DrugAdministrationHistoryPagination `json:"pagination"`
}
