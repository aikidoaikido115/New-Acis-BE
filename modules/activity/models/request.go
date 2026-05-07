package models

import "time"

type CreateActivityRequest struct {
	ActivityName string  `json:"activity_name" binding:"required"`
	ActivityType string  `json:"activity_type" binding:"required"`
	Description  *string `json:"description"`
	Location     *string `json:"location"`
}

type UpdateActivityRequest struct {
	StaffID      *string `json:"staff_id"`
	ActivityName *string `json:"activity_name"`
	ActivityType *string `json:"activity_type"`
	Description  *string `json:"description"`
	Location     *string `json:"location"`
}

type CreateActivityScheduleRequest struct {
	ActivityID string    `json:"activity_id" binding:"required"`
	Date       time.Time `json:"date" binding:"required"`
	StartTime  time.Time `json:"start_time" binding:"required"`
	EndTime    time.Time `json:"end_time" binding:"required"`
}

type UpdateActivityScheduleRequest struct {
	ActivityID *string    `json:"activity_id"`
	Date       *time.Time `json:"date"`
	StartTime  *time.Time `json:"start_time"`
	EndTime    *time.Time `json:"end_time"`
}

type CreateActivityScheduleWithActivitySyncRequest struct {
	ActivityName string    `json:"activity_name" binding:"required"`
	ActivityType string    `json:"activity_type" binding:"required"`
	Date         time.Time `json:"date" binding:"required"`
	StartTime    time.Time `json:"start_time" binding:"required"`
	EndTime      time.Time `json:"end_time" binding:"required"`
	Location     *string   `json:"location"`
	Description  *string   `json:"description"`
}

type UpdateActivityScheduleWithActivitySyncRequest struct {
	ActivityName *string    `json:"activity_name"`
	ActivityType *string    `json:"activity_type"`
	Date         *time.Time `json:"date"`
	StartTime    *time.Time `json:"start_time"`
	EndTime      *time.Time `json:"end_time"`
	Location     *string    `json:"location"`
	Description  *string    `json:"description"`
}

type ParticipationImageURLRequest struct {
	URL string `json:"url"`
}

type CreateParticipationRequest struct {
	ResidentID      string `json:"resident_id" binding:"required"`
	ASID            string `json:"as_id" binding:"required"`
	IsParticipating bool   `json:"is_participating"`
}

type UpdateParticipationRequest struct {
	IsParticipating *bool `json:"is_participating"`
	ClearImage      *bool  `json:"clear_image" form:"clear_image"`
}

type BulkUpdateParticipationIsParticipatingByResidentIDsRequest struct {
	ASID            string   `json:"as_id" binding:"required"`
	ResidentIDs     []string `json:"resident_ids" binding:"required"`
	IsParticipating *bool    `json:"is_participating"`
}

type ResidentsByScheduleQueryParams struct {
	Search   *string  `json:"search" form:"search" query:"search"`
	Floor    *int16   `json:"floor" form:"floor" query:"floor"`
	LabelIDs []string `json:"label_ids" form:"label_ids" query:"label_ids"`
	Page     *int     `json:"page" form:"page" query:"page"`
	PageSize *int     `json:"page_size" form:"page_size" query:"page_size"`
	Limit    int      `json:"limit" form:"limit" query:"limit"`
	Offset   int      `json:"offset" form:"offset" query:"offset"`
}
