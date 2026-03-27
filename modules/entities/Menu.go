package entities

type Menu struct {
	ID         string  `json:"menu_id" gorm:"primaryKey" `
	MenuName   string  `json:"menu_name" gorm:"not null"`
	Description string  `json:"description" gorm:"not null"` //such as ingredients
}
