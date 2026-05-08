package entities

type Staff struct {
	ID     string `json:"staff_id" gorm:"primaryKey" `
	UserID string `json:"user_id" gorm:"not null"`

	User User `json:"user" gorm:"foreignKey:UserID;references:ID"`
}
