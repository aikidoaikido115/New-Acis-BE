package seed

import (
	"log"

	"github.com/aikidoaikido115/New-Acis-BE/modules/entities"
	user_constants "github.com/aikidoaikido115/New-Acis-BE/modules/user/constants"
	"gorm.io/gorm"
)

// SeedRoles creates default roles in the database
func SeedRoles(db *gorm.DB) {
	log.Println("Seeding roles...")

	roles := []entities.Role{
		{
			ID:   "1",
			Name: user_constants.RoleMedicalStaff,
		},
		{
			ID:   "2",
			Name: user_constants.RoleKitchenStaff,
		},
		{
			ID:   "3",
			Name: user_constants.RoleRelative,
		},
		{
			ID:   "4",
			Name: user_constants.RoleSuperAdmin,
		},
	}

	for _, role := range roles {
		// ตรวจสอบว่า role มีอยู่แล้วหรือยัง (ใช้ name เป็นเกณฑ์)
		var existingRole entities.Role
		result := db.Where("name = ?", role.Name).First(&existingRole)

		if result.Error != nil {
			// ไม่เจอ role → สร้างใหม่
			if err := db.Create(&role).Error; err != nil {
				log.Printf("❌ Failed to seed role '%s': %v", role.Name, err)
			} else {
				log.Printf("✅ Seeded role: %s (ID: %s)", role.Name, role.ID)
			}
		} else {
			log.Printf("⏭️  Role already exists: %s (ID: %s)", existingRole.Name, existingRole.ID)
		}
	}

	log.Println("Roles seeding completed!")
}
