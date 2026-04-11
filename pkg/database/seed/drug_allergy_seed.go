package seed

import (
	"log"
	"strings"

	"github.com/aikidoaikido115/New-Acis-BE/modules/entities"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// SeedDrugAllergies creates default drug allergy master data in the database.
func SeedDrugAllergies(db *gorm.DB) {
	log.Println("Seeding drug allergies...")

	drugAllergies := []entities.DrugAllergy{
		{
			ID:          uuid.New().String(),
			AllergyName: "Penicillin",
		},
		{
			ID:          uuid.New().String(),
			AllergyName: "Sulfonamides",
		},
		{
			ID:          uuid.New().String(),
			AllergyName: "NSAIDs",
		},
		{
			ID:          uuid.New().String(),
			AllergyName: "Aspirin",
		},
		{
			ID:          uuid.New().String(),
			AllergyName: "Cephalosporins",
		},
	}

	for _, drugAllergy := range drugAllergies {
		drugAllergy.AllergyName = strings.TrimSpace(drugAllergy.AllergyName)
		if drugAllergy.AllergyName == "" {
			continue
		}

		var existingDrugAllergy entities.DrugAllergy
		result := db.Where("allergy_name = ?", drugAllergy.AllergyName).First(&existingDrugAllergy)

		if result.Error != nil {
			if err := db.Create(&drugAllergy).Error; err != nil {
				log.Printf("Failed to seed drug allergy '%s': %v", drugAllergy.AllergyName, err)
			} else {
				log.Printf("Seeded drug allergy: %s (ID: %s)", drugAllergy.AllergyName, drugAllergy.ID)
			}
		} else {
			log.Printf("Drug allergy already exists: %s (ID: %s)", existingDrugAllergy.AllergyName, existingDrugAllergy.ID)
		}
	}

	log.Println("Drug allergies seeding completed!")
}
