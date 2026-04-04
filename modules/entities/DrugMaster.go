package entities

type DrugMaster struct {
	ID   string `json:"dm_id" gorm:"primaryKey"`
	Name string `json:"name" gorm:"not null;unique"`
	Dose string `json:"dose" gorm:"not null"`
}
