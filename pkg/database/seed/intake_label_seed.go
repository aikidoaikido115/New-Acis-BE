package seed

import (
	"log"

	emr_constants "github.com/aikidoaikido115/New-Acis-BE/modules/emr/constants"
	"github.com/aikidoaikido115/New-Acis-BE/modules/entities"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// SeedRoles creates default roles in the database
func SeedIntakeLabels(db *gorm.DB) {
	log.Println("Seeding intake labels...")

	intakeLabels := []entities.IntakeLabels{
		{
			ID:        uuid.New().String(),
			LabelName: emr_constants.CareLevelFull,
		},
		{
			ID:        uuid.New().String(),
			LabelName: "ทานอาหารทางสายยาง",
		},
		{
			ID:        uuid.New().String(),
			LabelName: emr_constants.CareLevelPartial,
		},
		{
			ID:        uuid.New().String(),
			LabelName: emr_constants.CareLevelIndependent,
		},
		{
			ID:        uuid.New().String(),
			LabelName: "ใช้รถเข็น",
		},
		{
			ID:        uuid.New().String(),
			LabelName: "ใช้วอคเกอร์",
		},
	}

	for _, intakeLabel := range intakeLabels {
		var existingIntakeLabel entities.IntakeLabels
		result := db.Where("label_name = ?", intakeLabel.LabelName).First(&existingIntakeLabel)

		if result.Error != nil {
			if err := db.Create(&intakeLabel).Error; err != nil {
				log.Printf("❌ Failed to seed intake label '%s': %v", intakeLabel.LabelName, err)
			} else {
				log.Printf("✅ Seeded intake label: %s (ID: %s)", intakeLabel.LabelName, intakeLabel.ID)
			}
		} else {
			log.Printf("⏭️  Intake label already exists: %s (ID: %s)", existingIntakeLabel.LabelName, existingIntakeLabel.ID)
		}
	}

	log.Println("Intake labels seeding completed!")
}
