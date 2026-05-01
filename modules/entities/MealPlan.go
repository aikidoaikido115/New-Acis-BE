package entities

import "time"

type MealPlan struct {
	ID               string    `json:"meal_plan_id" gorm:"primaryKey" `
	MenuID           string    `json:"menu_id" gorm:"not null"`
	BackUpMenuID     *string   `json:"backup_menu_id"`
	MainAmount       int16     `json:"main_amount" gorm:"not null"`
	BackUpAmount     *int16    `json:"backup_amount"`
	IsAllergy        bool      `json:"is_allergy"`
	MealType         string    `json:"meal_type" gorm:"not null"` //breakfast, lunch, dinner
	CreatedByStaffID string    `json:"created_by_staff_id" gorm:"type:text"`
	StaffName        string    `json:"staff_name" gorm:"type:text"`
	CreatedAt        time.Time `json:"created_at"`

	Menu Menu `json:"-" gorm:"foreignKey:MenuID;references:ID"`
}
