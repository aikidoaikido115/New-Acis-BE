package entities

// Relative เก็บข้อมูลญาติของผู้ป่วยที่เชื่อมโยงกัน
type Relative struct {
	ID               string `json:"relative_id" gorm:"primaryKey"`
	UserID           string `json:"user_id" gorm:"index"` // สามารถ null ได้ถ้ายังไม่ผูก User Account
	ResidentID       string `json:"resident_id" gorm:"not null;index"`
	RelativePassword string `json:"relative_password"` // ควรเก็บเป็น Hash
	Relation         string `json:"relation" gorm:"not null"`
	Phone            string `json:"phone"`

	// Relationships
	Resident     *Resident     `json:"-" gorm:"foreignKey:ResidentID;references:ID"`
	DailyUpdates []DailyUpdate `json:"daily_updates" gorm:"foreignKey:RelativeID;references:ID"`
}