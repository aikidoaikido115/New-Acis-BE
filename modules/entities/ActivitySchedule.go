package entities

import "time"

type ActivitySchedule struct {
	ID         string    `json:"as_id" gorm:"primaryKey"`
	ActivityID string    `json:"activity_id" gorm:"not null"`
	Date       time.Time `json:"date" gorm:"not null"`
	StartTime  time.Time `json:"start_time" gorm:"not null"`
	EndTime    time.Time `json:"end_time" gorm:"not null"`
	CreatedAt  time.Time `json:"created_at" gorm:"not null"`
	UpdatedAt  time.Time `json:"updated_at" gorm:"not null"`

	Activity Activity `json:"activity" gorm:"foreignKey:ActivityID;references:ID"`
}
