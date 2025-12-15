package entities

import "time"

type User struct {
	ID                string    `json:"user_id" gorm:"primaryKey" `
	Username          string    `json:"username" gorm:"unique;not null"`
	NumberOfUsernames int       `json:"number_of_usernames" gorm:"default:0"`
	Email             string    `json:"email" gorm:"unique;not null"`
	Password          string    `json:"-"`
	ProfileImage      string    `json:"profile_image" gorm:"default:https://www.isranews.org/article/images/2025/Harry/6/Hun_Sen_July_2019.jpg"`
	CreatedAt         time.Time `json:"created_at"`
	UpdatedAt         time.Time `json:"updated_at"`
}
