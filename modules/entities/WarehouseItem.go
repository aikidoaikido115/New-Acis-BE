package entities

import "time"

type WarehouseItem struct {
	ID              string    `json:"id" gorm:"primaryKey"`
	Code            string    `json:"code" gorm:"not null;uniqueIndex"`
	Name            string    `json:"name" gorm:"not null"`
	Description     string    `json:"description"`
	Quantity        int       `json:"quantity" gorm:"not null;default:0"`
	MinimumQuantity int       `json:"minimumQuantity" gorm:"not null;default:0"`
	Unit            string    `json:"unit" gorm:"not null"`
	Category        string    `json:"category" gorm:"not null"`
	CreatedAt       time.Time `json:"createdAt"`
	UpdatedAt       time.Time `json:"updatedAt"`
}
