package seed

import (
	"log"

	"github.com/aikidoaikido115/New-Acis-BE/modules/entities"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// SeedDrugMasters creates default drug master data in the database.
func SeedDrugMasters(db *gorm.DB) {
	log.Println("Seeding drug masters...")

	drugs := []entities.DrugMaster{
		{
			ID:   uuid.New().String(),
			Name: "Paracetamol",
			Dose: "500 mg",
		},
		{
			ID:   uuid.New().String(),
			Name: "Aspirin",
			Dose: "81 mg",
		},
		{
			ID:   uuid.New().String(),
			Name: "Amoxicillin",
			Dose: "500 mg",
		},
		{
			ID:   uuid.New().String(),
			Name: "Metformin",
			Dose: "500 mg",
		},
		{
			ID:   uuid.New().String(),
			Name: "Atorvastatin",
			Dose: "20 mg",
		},
		{
			ID:   uuid.New().String(),
			Name: "Amlodipine",
			Dose: "5 mg",
		},
		{
			ID:   uuid.New().String(),
			Name: "Losartan",
			Dose: "50 mg",
		},
		{
			ID:   uuid.New().String(),
			Name: "Omeprazole",
			Dose: "20 mg",
		},
	}

	for _, drug := range drugs {
		var existingDrug entities.DrugMaster
		result := db.Where("name = ?", drug.Name).First(&existingDrug)

		if result.Error != nil {
			if err := db.Create(&drug).Error; err != nil {
				log.Printf("❌ Failed to seed drug '%s': %v", drug.Name, err)
			} else {
				log.Printf("✅ Seeded drug: %s (ID: %s)", drug.Name, drug.ID)
			}
		} else {
			log.Printf("⏭️  Drug already exists: %s (ID: %s)", existingDrug.Name, existingDrug.ID)
		}
	}

	log.Println("Drug masters seeding completed!")
}
