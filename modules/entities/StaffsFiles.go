package entities

type StaffsFiles struct {
	ID       	string `json:"file_id" gorm:"primaryKey" `
	StaffID  	string `json:"staff_id" gorm:"not null"`
	FileName 	string `json:"file_name" gorm:"not null"`
	FileType 	string `json:"file_type" gorm:"not null"`
	FileSize 	int64  `json:"file_size" gorm:"not null"`
	File     	string `json:"file" gorm:"not null"`

	Staff 		Staff `json:"-" gorm:"foreignKey:StaffID;references:ID"`
}
