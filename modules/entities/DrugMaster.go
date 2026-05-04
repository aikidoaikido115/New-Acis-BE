package entities

type DrugMaster struct {
	ID   string `json:"dm_id" gorm:"primaryKey"`
	Name string `json:"name" gorm:"not null;uniqueIndex:idx_drug_masters_name_dose"`
	Dose string `json:"dose" gorm:"not null;uniqueIndex:idx_drug_masters_name_dose"`
}
