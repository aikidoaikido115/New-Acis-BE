package entities

import "time"

type PersonalDrug struct {
	ID          string     `json:"pd_id" gorm:"primaryKey"`
	ResidentID  string     `json:"resident_id" gorm:"not null"`
	DmID        string     `json:"dm_id" gorm:"not null"`
	Amount      string     `json:"amount" gorm:"not null"`
	AmountUnit  string     `json:"amount_unit" gorm:"not null"` //เม็ด ช้อนชา
	Frequency   int        `json:"frequency" gorm:"not null"`   //จำนวนครั้งต่อวัน
	TimeOfDay   string     `json:"time_of_day" gorm:"not null"` //เช้า่ กลางวัน เย็น ก่อนนอน
	Timing      string     `json:"timing" gorm:"not null"`      //ก่อนอาหาร หลังอาหาร
	Description string     `json:"description"`
	TakeType    string     `json:"take_type" gorm:"not null"` //regular, as_needed
	StartDate   *time.Time `json:"start_date,omitempty" gorm:"type:date"`
	EndDate     *time.Time `json:"end_date,omitempty" gorm:"type:date"`
	CreatedAt   time.Time  `json:"created_at" gorm:"not null"`
	UpdatedAt   time.Time  `json:"updated_at" gorm:"not null"`

	Resident   Resident   `gorm:"foreignKey:ResidentID;references:ID"`
	DrugMaster DrugMaster `gorm:"foreignKey:DmID;references:ID"`
}
