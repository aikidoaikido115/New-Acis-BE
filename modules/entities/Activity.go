package entities

type Activity struct {
	ID        string `json:"activity_id" gorm:"primaryKey" `
	StaffID   string `json:"staff_id" gorm:"not null"`
	ActivityName string `json:"activity_name" gorm:"not null"`
	ActivityType string `json:"activity_type" gorm:"not null"`
	Description *string `json:"description"`
	Location *string `json:"location"`

	Staff Staff `json:"staff" gorm:"foreignKey:StaffID;references:ID"`
}
