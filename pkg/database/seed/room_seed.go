package seed

import (
	"log"

	"github.com/aikidoaikido115/New-Acis-BE/modules/entities"
	"gorm.io/gorm"
	"github.com/google/uuid"

)

// SeedRoles creates default roles in the database
func SeedRooms(db *gorm.DB) {
	log.Println("Seeding rooms...")

	rooms := []entities.Room{
		{
			ID:   uuid.New().String(),
			Floor: 1,
			RoomNumber: "101",
		},
		{
			ID:   uuid.New().String(),
			Floor: 1,
			RoomNumber: "102",
		},
		{
			ID:   uuid.New().String(),
			Floor: 4,
			RoomNumber: "401",
		},
		{
			ID:   uuid.New().String(),
			Floor: 4,
			RoomNumber: "402",
		},
	}

	for _, room := range rooms {
		var existingRoom entities.Room
		result := db.Where("room_number = ?", room.RoomNumber).First(&existingRoom)

		if result.Error != nil {
			if err := db.Create(&room).Error; err != nil {
				log.Printf("❌ Failed to seed room '%s': %v", room.RoomNumber, err)
			} else {
				log.Printf("✅ Seeded room: %s (ID: %s)", room.RoomNumber, room.ID)
			}
		} else {
			log.Printf("⏭️  Room already exists: %s (ID: %s)", existingRoom.RoomNumber, existingRoom.ID)
		}
	}

	log.Println("Rooms seeding completed!")
}
