package entities

import "time"

// DailyFile เก็บข้อมูลไฟล์หรือรูปภาพที่แนบไปกับอัปเดตประจำวัน
type DailyFile struct {
	ID        string    `json:"file_id" gorm:"primaryKey"`
	DailyID   string    `json:"daily_id" gorm:"not null;index"`
	FileName  string    `json:"file_name" gorm:"not null"`
	FileType  string    `json:"file_type"`
	FileSize  int64     `json:"file_size"`
	URL       string    `json:"url" gorm:"not null"`
	CreatedAt time.Time `json:"created_at" gorm:"not null"`
	UpdatedAt time.Time `json:"updated_at" gorm:"not null"`

	// Relationships
	DailyUpdate *DailyUpdate `json:"-" gorm:"foreignKey:DailyID;references:ID"`
}
