package database

import (
	"fmt"
	"log"

	"github.com/aikidoaikido115/New-Acis-BE/configs"
	"github.com/aikidoaikido115/New-Acis-BE/modules/entities"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var db *gorm.DB

func InitDB(config configs.PostgreSQL) {
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=%s",
		config.Host,
		config.Username,
		config.Password,
		config.Database,
		config.Port,
		config.SSLMode,
	)

	var err error
	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Database connection failed: %v", err)
	}

	if err := db.AutoMigrate(
		&entities.User{},
		&entities.OTP{},
		&entities.TempToken{},
	); err != nil {
		log.Fatalf("Database migration failed: %v", err)
	}

	log.Printf("Database connected: %s@%s:%s/%s", config.Username, config.Host, config.Port, config.Database)
}

func GetDB() *gorm.DB {
	if db == nil {
		log.Fatal("Database not initialized. Call InitDB() first")
	}
	return db
}
