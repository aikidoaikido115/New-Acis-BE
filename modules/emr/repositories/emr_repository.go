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
	RoomExists(id string) (bool, error)
	GetResidentByID(id string) (*entities.Resident, error)
	GetResidentByRoomID(roomID string) ([]*entities.Resident, error)
	GetAllResidents() ([]*entities.Resident, error)
	GetRoomByID(id string) (*entities.Room, error)
	GetAllRooms() ([]*entities.Room, error)
	//todo UpdateResident
	//todo SoftDeleteResident
	//todo IntakeLabels
	//todo Allergy
	//todo get vital sign เพื่อเอาไป dashboard อันผิดปกติ
}

func (r *GormEmrRepository) CreateResident(resident *entities.Resident) (*entities.Resident, error) {
	err := r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&resident).Error; err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	return r.GetResidentByID(resident.ID)
}

func (r *GormEmrRepository) RoomExists(id string) (bool, error) {
	var count int64
	err := r.db.Model(&entities.Room{}).Where("id = ?", id).Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r *GormEmrRepository) GetResidentByID(id string) (*entities.Resident, error) {
	var resident entities.Resident
	if err := r.db.Preload("Room").Where("id = ?", id).First(&resident).Error; err != nil {
		return nil, err
	}
	return &resident, nil
}

func (r *GormEmrRepository) GetResidentByRoomID(roomID string) ([]*entities.Resident, error) {
	var residents []*entities.Resident
	if err := r.db.Preload("Room").Where("room_id = ?", roomID).Find(&residents).Error; err != nil {
		return nil, err
	}
	return residents, nil
}

func (r *GormEmrRepository) GetAllResidents() ([]*entities.Resident, error) {
	var residents []*entities.Resident
	if err := r.db.Preload("Room").Find(&residents).Error; err != nil {
		return nil, err
	}
	return residents, nil
}

func (r *GormEmrRepository) GetRoomByID(id string) (*entities.Room, error) {
	var room entities.Room
	if err := r.db.Where("id = ?", id).First(&room).Error; err != nil {
		return nil, err
	}
	return &room, nil
}
func (r *GormEmrRepository) GetAllRooms() ([]*entities.Room, error) {
	var rooms []*entities.Room
	if err := r.db.Find(&rooms).Error; err != nil {
		return nil, err
	}
	return rooms, nil
}
