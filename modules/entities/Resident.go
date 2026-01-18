package entities

type Resident struct {
	ID       	string `json:"resident_id" gorm:"primaryKey" `
}
