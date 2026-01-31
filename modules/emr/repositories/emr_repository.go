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
	// Resident operations
	CreateResident(resident *entities.Resident) (*entities.Resident, error)
	GetResidentByID(id string) (*entities.Resident, error)
	GetResidentByRoomID(roomID string) ([]*entities.Resident, error)
	GetAllResidents() ([]*entities.Resident, error)
	GetNumberOfResidents() (int16, error)

	// Room operations
	RoomExists(id string) (bool, error)
	GetRoomByID(id string) (*entities.Room, error)
	GetAllRooms() ([]*entities.Room, error)

	// IntakeLabel operations
	CreateIntakeLabel(label *entities.IntakeLabels) (*entities.IntakeLabels, error)
	GetIntakeLabelByID(id string) (*entities.IntakeLabels, error)
	GetIntakeLabelByName(labelName string) (*entities.IntakeLabels, error)
	GetAllIntakeLabels() ([]*entities.IntakeLabels, error)

	// ResidentLabel operations (many-to-many)
	CreateIntakeLabelByResidentID(residentLabel *entities.ResidentLabels) (*entities.ResidentLabels, error)
	GetResidentLabelsByResidentID(residentID string) ([]*entities.ResidentLabels, error)

	//todo UpdateResident
	//todo SoftDeleteResident
	//todo Allergy
	//todo get vital sign เพื่อเอาไป dashboard อันผิดปกติ
}

func (r *GormEmrRepository) CreateResident(resident *entities.Resident) (*entities.Resident, error) {
	if err := r.db.Create(&resident).Error; err != nil {
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

func (r *GormEmrRepository) GetNumberOfResidents() (int16, error) {
	var count int64
	if err := r.db.Model(&entities.Resident{}).Count(&count).Error; err != nil {
		return 0, err
	}
	return int16(count), nil
}

func (r *GormEmrRepository) CreateIntakeLabelByResidentID(residentLabel *entities.ResidentLabels) (*entities.ResidentLabels, error) {
	if err := r.db.Create(&residentLabel).Error; err != nil {
		return nil, err
	}

	var result entities.ResidentLabels
	if err := r.db.Preload("Resident").Preload("IntakeLabel").
		Where("resident_id = ? AND label_id = ?", residentLabel.ResidentID, residentLabel.LabelID).
		First(&result).Error; err != nil {
		return nil, err
	}
	return &result, nil
}

func (r *GormEmrRepository) GetIntakeLabelByName(labelName string) (*entities.IntakeLabels, error) {
	var label entities.IntakeLabels
	if err := r.db.Where("label_name = ?", labelName).First(&label).Error; err != nil {
		return nil, err
	}
	return &label, nil
}

func (r *GormEmrRepository) CreateIntakeLabel(label *entities.IntakeLabels) (*entities.IntakeLabels, error) {
	if err := r.db.Create(&label).Error; err != nil {
		return nil, err
	}
	return label, nil
}

func (r *GormEmrRepository) GetIntakeLabelByID(id string) (*entities.IntakeLabels, error) {
	var label entities.IntakeLabels
	if err := r.db.Where("id = ?", id).First(&label).Error; err != nil {
		return nil, err
	}
	return &label, nil
}

func (r *GormEmrRepository) GetAllIntakeLabels() ([]*entities.IntakeLabels, error) {
	var labels []*entities.IntakeLabels
	if err := r.db.Find(&labels).Error; err != nil {
		return nil, err
	}
	return labels, nil
}

func (r *GormEmrRepository) GetResidentLabelsByResidentID(residentID string) ([]*entities.ResidentLabels, error) {
	var residentLabels []*entities.ResidentLabels
	if err := r.db.Preload("IntakeLabel").Where("resident_id = ?", residentID).Find(&residentLabels).Error; err != nil {
		return nil, err
	}
	return residentLabels, nil
}
