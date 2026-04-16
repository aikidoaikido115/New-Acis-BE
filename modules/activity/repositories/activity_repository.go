package repositories

import (
	"github.com/aikidoaikido115/New-Acis-BE/modules/entities"
	"gorm.io/gorm"
)

type GormActivityRepository struct {
	db *gorm.DB
}

func NewGormActivityRepository(db *gorm.DB) *GormActivityRepository {
	return &GormActivityRepository{
		db: db,
	}
}

type ActivityRepository interface {
	CreateActivity(activity *entities.Activity) (*entities.Activity, error)
	GetActivityByID(id string) (*entities.Activity, error)
	GetAllActivities() ([]*entities.Activity, error)
	UpdateActivity(activity *entities.Activity) (*entities.Activity, error)
	DeleteActivity(id string) error
}

func (r *GormActivityRepository) CreateActivity(activity *entities.Activity) (*entities.Activity, error) {
	if err := r.db.Create(&activity).Error; err != nil {
		return nil, err
	}

	return r.GetActivityByID(activity.ID)
}

func (r *GormActivityRepository) GetActivityByID(id string) (*entities.Activity, error) {
	var activity entities.Activity
	if err := r.db.Preload("Staff").Where("id = ?", id).First(&activity).Error; err != nil {
		return nil, err
	}

	return &activity, nil
}

func (r *GormActivityRepository) GetAllActivities() ([]*entities.Activity, error) {
	var activities []*entities.Activity
	if err := r.db.Preload("Staff").Find(&activities).Error; err != nil {
		return nil, err
	}

	return activities, nil
}

func (r *GormActivityRepository) UpdateActivity(activity *entities.Activity) (*entities.Activity, error) {
	if err := r.db.Save(&activity).Error; err != nil {
		return nil, err
	}

	return r.GetActivityByID(activity.ID)
}

func (r *GormActivityRepository) DeleteActivity(id string) error {
	if err := r.db.Delete(&entities.Activity{}, "id = ?", id).Error; err != nil {
		return err
	}

	return nil
}
