package entities

type TempToken struct {
	UserID string `json:"-" gorm:"primaryKey"`
	Token  string `json:"token" gorm:"not null;unique"`
	User   User   `json:"-" gorm:"foreignKey:UserID;references:ID"`
}
