package repositories

import (
	"github.com/aikidoaikido115/New-Acis-BE/modules/entities"

	"gorm.io/gorm"
)

type GormEmrRepository struct {
	db *gorm.DB
}

func NewGormEmrRepository(db *gorm.DB) *GormEmrRepository {
	return &GormEmrRepository{
		db: db,
	}
}

type EmrRepository interface {
	CreateResident(user *entities.Resident) (*entities.Resident, error)
}