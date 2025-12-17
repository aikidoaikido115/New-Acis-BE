package entities

type StaffsFiles struct {
	ID      string `json:"file_id" gorm:"primaryKey" `
	StaffID string `json:"staff_id" gorm:"not null"`

	Staff Staff `json:"staff" gorm:"foreignKey:StaffID;references:ID"`
}
