package entities

import "time"

type AuditLogs struct {
	ID        string    `json:"-" gorm:"primaryKey"`
	TableName string    `json:"table_name" gorm:"not null"`
	RecordID  string    `json:"record_id" gorm:"not null"`
	UserID    string    `json:"user_id" gorm:"not null"`
	Action    string    `json:"action" gorm:"not null"`
	OldValue  string    `json:"old_value" gorm:"type:text"`
	NewValue  string    `json:"new_value" gorm:"type:text"`
	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime:milli"`
}
