package entities

type DrugAllergy struct {
	ID          string `json:"drug_allergy_id" gorm:"primaryKey" `
	AllergyName string `json:"allergy_name" gorm:"not null;uniqueIndex"`

	ResidentDA []ResidentDA `json:"resident_das" gorm:"foreignKey:DrugAllergyID;references:ID"`
}
