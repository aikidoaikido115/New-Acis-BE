package entities

type Room struct {
	ID         string  `json:"room_id" gorm:"primaryKey" `
	StaffID    *string `json:"staff_id"` //pointer to allow null value
	Floor      int16   `json:"floor" gorm:"not null"`
	RoomNumber string  `json:"room_number" gorm:"not null;uniqueIndex"`

	Staff Staff `json:"-" gorm:"foreignKey:StaffID;references:ID"`
}
