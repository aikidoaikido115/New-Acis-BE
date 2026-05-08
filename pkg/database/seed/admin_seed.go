package seed

import (
	"errors"
	"log"

	"github.com/aikidoaikido115/New-Acis-BE/configs"
	"github.com/aikidoaikido115/New-Acis-BE/modules/entities"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// SeedAdminUser creates default admin user from configuration
func SeedAdminUser(db *gorm.DB, seedAdmin configs.SeedAdmin) {
	log.Println("Seeding admin user...")

	// Check if admin already exists (by email or username)
	var existingUser entities.User
	result := db.Where("email = ? OR username = ?", seedAdmin.Email, seedAdmin.Username).First(&existingUser)

	if result.Error == nil {
		var existingStaff entities.Staff
		staffResult := db.Where("user_id = ?", existingUser.ID).First(&existingStaff)
		if staffResult.Error != nil {
			if errors.Is(staffResult.Error, gorm.ErrRecordNotFound) {
				staff := &entities.Staff{
					ID:     uuid.New().String(),
					UserID: existingUser.ID,
				}

				if err := db.Create(staff).Error; err != nil {
					log.Printf("❌ Failed to create staff record for existing admin user '%s': %v", existingUser.Email, err)
					return
				}

				log.Printf("✅ Created staff record for existing admin user: %s (staff ID: %s)", existingUser.Email, staff.ID)
			} else {
				log.Printf("❌ Failed to check staff record for '%s': %v", existingUser.Email, staffResult.Error)
				return
			}
		}

		log.Printf("⏭️  Admin user already exists: %s (ID: %s)", existingUser.Email, existingUser.ID)
		return
	}

	// Get Admin role
	var adminRole entities.Role
	roleResult := db.Where("name = ?", seedAdmin.RoleName).First(&adminRole)
	if roleResult.Error != nil {
		log.Printf("❌ Failed to find role '%s': %v", seedAdmin.RoleName, roleResult.Error)
		return
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(seedAdmin.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Printf("❌ Failed to hash password: %v", err)
		return
	}

	// Create admin user
	adminUser := &entities.User{
		ID:        uuid.New().String(),
		RoleID:    adminRole.ID,
		Username:  seedAdmin.Username,
		Email:     seedAdmin.Email,
		Password:  string(hashedPassword),
		IsApprove: true, // Admin is auto-approved
		FirstName: seedAdmin.FirstName,
		LastName:  seedAdmin.LastName,
		Nickname:  seedAdmin.Nickname,
	}

	if err := db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(adminUser).Error; err != nil {
			return err
		}

		adminStaff := &entities.Staff{
			ID:     uuid.New().String(),
			UserID: adminUser.ID,
		}

		if err := tx.Create(adminStaff).Error; err != nil {
			return err
		}

		return nil
	}); err != nil {
		log.Printf("❌ Failed to seed admin user '%s': %v", seedAdmin.Email, err)
		return
	}

	log.Printf("✅ Seeded admin user: %s (ID: %s)", adminUser.Email, adminUser.ID)
	log.Printf("✅ Created staff record for admin user: %s", adminUser.Email)
	log.Println("Admin user seeding completed!")
}
