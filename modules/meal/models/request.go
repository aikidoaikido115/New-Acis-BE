package models

type CreateMenuRequest struct {
	MenuName    string `json:"menu_name" binding:"required"`
	Description string `json:"description" binding:"required"`
}

type UpdateMenuRequest struct {
	MenuName    *string `json:"menu_name"`
	Description *string `json:"description"`
}

type CreateMealPlanRequest struct {
	MenuID         string  `json:"menu_id" binding:"required"`
	BackUpMenuID   *string `json:"backup_menu_id"`
	MainAmount     int16   `json:"main_amount" binding:"required"`
	BackUpAmount   *int16  `json:"backup_amount"`
	MealType       string  `json:"meal_type" binding:"required"`
	HumanInTheLoop *bool   `json:"human_in_the_loop"`
}

type UpdateMealPlanRequest struct {
	MenuID       *string `json:"menu_id"`
	BackUpMenuID *string `json:"backup_menu_id"`
	MainAmount   *int16  `json:"main_amount"`
	BackUpAmount *int16  `json:"backup_amount"`
	MealType     *string `json:"meal_type"`
}

type MealHistoryQueryParams struct {
	Date     *string `json:"date" form:"date" query:"date"`
	MealType *string `json:"meal_type" form:"meal_type" query:"meal_type"`
	Search   *string `json:"search" form:"search" query:"search"`
	Page     *int    `json:"page" form:"page" query:"page"`
	PageSize *int    `json:"page_size" form:"page_size" query:"page_size"`
}
