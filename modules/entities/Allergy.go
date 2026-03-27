package entities

type Allergy struct {
	ID        string `json:"allergy_id" gorm:"primaryKey" `
	AllergyName string `json:"allergy_name" gorm:"not null;uniqueIndex"`

	ResidentAllergies []ResidentAllergies `json:"resident_allergies" gorm:"foreignKey:AllergyID;references:ID"`
}
