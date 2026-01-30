package entities

type Resident struct {
	ID       	string `json:"resident_id" gorm:"primaryKey" `
	RoomID  	string `json:"room_id" gorm:"not null"`
	FirstName  	string `json:"first_name" gorm:"not null"`
	LastName   	string `json:"last_name" gorm:"not null"`
	Age        	*int16    `json:"age" gorm:"not null"` // ใช้ pointer เพื่อให้สามารถเช็ค null ได้
	Gender	 	string `json:"gender" gorm:"not null"`

	Room 		Room   `json:"-" gorm:"foreignKey:RoomID;references:ID"`
	ResidentLabels []ResidentLabels `json:"resident_labels" gorm:"foreignKey:ResidentID;references:ID"`
}
