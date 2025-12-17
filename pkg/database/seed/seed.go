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

	// อนาคตเพิ่ม seeder อื่นๆ ตรงนี้
	// SeedUsers(db)
	// SeedCategories(db)

	log.Println("Database seeding completed!")
}
