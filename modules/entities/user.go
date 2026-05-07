package entities

import "time"

type User struct {
	ID           string    `json:"user_id" gorm:"primaryKey" `
	RoleID       string    `json:"role_id" gorm:"not null"`
	Username     string    `json:"username" gorm:"unique;not null"`
	Email        string    `json:"email" gorm:"unique;not null"`
	Password     string    `json:"-"`
	IsApprove    bool      `json:"is_approve" gorm:"not null;default:false"`
	FirstName    string    `json:"first_name" gorm:"default:null"`
	LastName     string    `json:"last_name" gorm:"default:null"`
	Nickname     string    `json:"nickname" gorm:"default:null"`
	Gender       string    `json:"gender" gorm:"default:null"`
	Phone        string    `json:"phone" gorm:"default:null"`
	ProfileImage string    `json:"profile_image" gorm:"default:https://www.isranews.org/article/images/2025/Harry/6/Hun_Sen_July_2019.jpg"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`

	Role Role `json:"role" gorm:"foreignKey:RoleID;references:ID"`
}
