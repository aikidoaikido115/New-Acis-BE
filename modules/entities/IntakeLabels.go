package entities

type IntakeLabels struct {
	ID        string `json:"label_id" gorm:"primaryKey" `
	LabelName string `json:"label_name" gorm:"not null;uniqueIndex"`

	ResidentLabels []ResidentLabels `json:"resident_labels" gorm:"foreignKey:LabelID;references:ID"`
}
