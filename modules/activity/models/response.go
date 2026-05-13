package models

import "time"

type ActivityScheduleWithActivitySyncResponse struct {
	ASID         string     `json:"as_id"`        
    ActivityID   string     `json:"activity_id"`
	ActivityName string    `json:"activity_name"`
	ActivityType string    `json:"activity_type"`
	Date         time.Time `json:"date"`
	StartTime    time.Time `json:"start_time"`
	EndTime      time.Time `json:"end_time"`
	Location     *string   `json:"location"`
	Description  *string   `json:"description"`
	SeriesID     *string   `json:"series_id"`
	Status       string    `json:"status"`
}

type ResidentByScheduleResponse struct {
	ResidentID      string   `json:"resident_id"`
	FirstName       string   `json:"first_name"`
	LastName        string   `json:"last_name"`
	Nickname        *string  `json:"nickname"`
	RoomNumber      string   `json:"room_number"`
	Floor           int16    `json:"floor"`
	IntakeLabels    []string `json:"intake_labels"`
	IsParticipating bool     `json:"is_participating"`
}

type ActivityPagination struct {
	Page       int `json:"page"`
	PageSize   int `json:"page_size"`
	TotalItems int `json:"total_items"`
	TotalPages int `json:"total_pages"`
}

type ResidentsByScheduleListResponse struct {
	Items      []*ResidentByScheduleResponse `json:"items"`
	Pagination ActivityPagination            `json:"pagination"`
}
