package models

type MealHistoryItem struct {
	Date           string  `json:"date" gorm:"column:date"`
	MealType       string  `json:"meal_type" gorm:"column:meal_type"`
	MenuName       string  `json:"menu_name" gorm:"column:menu_name"`
	MainAmount     int16   `json:"main_amount" gorm:"column:main_amount"`
	BackupMenuName *string `json:"backup_menu_name" gorm:"column:backup_menu_name"`
	BackupAmount   *int16  `json:"backup_amount" gorm:"column:backup_amount"`
	StaffName      string  `json:"staff_name" gorm:"column:staff_name"`
}

type MealHistoryPagination struct {
	Page       int `json:"page"`
	PageSize   int `json:"page_size"`
	TotalItems int `json:"total_items"`
	TotalPages int `json:"total_pages"`
}

type MealHistoryResponse struct {
	Items      []MealHistoryItem     `json:"items"`
	Pagination MealHistoryPagination `json:"pagination"`
}
