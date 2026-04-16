package models

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
