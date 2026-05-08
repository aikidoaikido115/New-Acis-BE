package seed

import (
	"log"
	"strings"

	"github.com/aikidoaikido115/New-Acis-BE/modules/entities"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// SeedAllergies creates default allergy master data in the database.
func SeedAllergies(db *gorm.DB) {
	log.Println("Seeding allergies...")

	allergies := []entities.Allergy{
		{
			ID:          uuid.New().String(),
			AllergyName: "นม",
		},
		{
			ID:          uuid.New().String(),
			AllergyName: "กุ้ง",
		},
		{
			ID:          uuid.New().String(),
			AllergyName: "อาหารทะเล",
		},
		{
			ID:          uuid.New().String(),
			AllergyName: "แพ้อาหารทะเลบางชนิด (ไม่ระบุชนิด)",
		},
		{
			ID:          uuid.New().String(),
			AllergyName: "สงสัยแพ้ผงชูรส (ต้องยืนยันซ้ำ)",
		},
	}

	for _, allergy := range allergies {
		allergy.AllergyName = strings.TrimSpace(allergy.AllergyName)
		if allergy.AllergyName == "" {
			continue
		}

		var existingAllergy entities.Allergy
		result := db.Where("allergy_name = ?", allergy.AllergyName).First(&existingAllergy)

		if result.Error != nil {
			if err := db.Create(&allergy).Error; err != nil {
				log.Printf("❌ Failed to seed allergy '%s': %v", allergy.AllergyName, err)
			} else {
				log.Printf("✅ Seeded allergy: %s (ID: %s)", allergy.AllergyName, allergy.ID)
			}
		} else {
			log.Printf("⏭️  Allergy already exists: %s (ID: %s)", existingAllergy.AllergyName, existingAllergy.ID)
		}
	}

	log.Println("Allergies seeding completed!")
}
