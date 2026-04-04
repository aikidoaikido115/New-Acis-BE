package repositories

import (
	"github.com/aikidoaikido115/New-Acis-BE/modules/entities"

	"gorm.io/gorm"
)

type GormDrugMasterRepository struct {
	db *gorm.DB
}

func NewGormDrugMasterRepository(db *gorm.DB) *GormDrugMasterRepository {
	return &GormDrugMasterRepository{
		db: db,
	}
}

type DrugMasterRepository interface {
	CreateDrugMaster(drug *entities.DrugMaster) (*entities.DrugMaster, error)
	GetDrugMasterByID(id string) (*entities.DrugMaster, error)
	GetDrugMasterByName(name string) (*entities.DrugMaster, error)
	GetAllDrugMasters() ([]*entities.DrugMaster, error)
	UpdateDrugMaster(drug *entities.DrugMaster) (*entities.DrugMaster, error)
	DeleteDrugMaster(id string) error
	DrugMasterExistsByName(name string) (bool, error)
}

func (r *GormDrugMasterRepository) CreateDrugMaster(drug *entities.DrugMaster) (*entities.DrugMaster, error) {
	if err := r.db.Create(&drug).Error; err != nil {
		return nil, err
	}

	return r.GetDrugMasterByID(drug.ID)
}

func (r *GormDrugMasterRepository) GetDrugMasterByID(id string) (*entities.DrugMaster, error) {
	var drug entities.DrugMaster
	if err := r.db.Where("id = ?", id).First(&drug).Error; err != nil {
		return nil, err
	}

	return &drug, nil
}

func (r *GormDrugMasterRepository) GetDrugMasterByName(name string) (*entities.DrugMaster, error) {
	var drug entities.DrugMaster
	if err := r.db.Where("name = ?", name).First(&drug).Error; err != nil {
		return nil, err
	}

	return &drug, nil
}

func (r *GormDrugMasterRepository) GetAllDrugMasters() ([]*entities.DrugMaster, error) {
	var drugs []*entities.DrugMaster
	if err := r.db.Order("name ASC").Find(&drugs).Error; err != nil {
		return nil, err
	}

	return drugs, nil
}

func (r *GormDrugMasterRepository) UpdateDrugMaster(drug *entities.DrugMaster) (*entities.DrugMaster, error) {
	if err := r.db.Save(&drug).Error; err != nil {
		return nil, err
	}

	return r.GetDrugMasterByID(drug.ID)
}

func (r *GormDrugMasterRepository) DeleteDrugMaster(id string) error {
	if err := r.db.Delete(&entities.DrugMaster{}, "id = ?", id).Error; err != nil {
		return err
	}

	return nil
}

func (r *GormDrugMasterRepository) DrugMasterExistsByName(name string) (bool, error) {
	var count int64
	err := r.db.Model(&entities.DrugMaster{}).Where("name = ?", name).Count(&count).Error
	if err != nil {
		return false, err
	}

	return count > 0, nil
}
