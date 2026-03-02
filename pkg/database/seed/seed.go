package seed

import (
	"log"

	"gorm.io/gorm"
)

// RunAll executes all seeders in order
func RunAll(db *gorm.DB) {
	log.Println("Starting database seeding...")

	// เรียก seeder แต่ละตัวตามลำดับ
	SeedRoles(db)
	SeedRooms(db)
	SeedIntakeLabels(db)

	log.Println("Database seeding completed!")
}
