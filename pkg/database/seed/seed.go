package seed

import (
	"log"

	"github.com/aikidoaikido115/New-Acis-BE/configs"
	"gorm.io/gorm"
)

// RunAll executes all seeders in order
func RunAll(db *gorm.DB, seedAdmin configs.SeedAdmin) {
	log.Println("Starting database seeding...")

	// เรียก seeder แต่ละตัวตามลำดับ
	SeedRoles(db)
	SeedRooms(db)
	SeedIntakeLabels(db)
	SeedAllergies(db)
	SeedDrugAllergies(db)
	SeedMenus(db)
	SeedDrugMasters(db)
	SeedActivities(db)
	SeedAdminUser(db, seedAdmin)

	log.Println("Database seeding completed!")
}
