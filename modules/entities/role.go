package entities

type Role struct {
	ID   string `json:"-" gorm:"primaryKey"`
	Name string `json:"name" gorm:"not null;unique"`
}
