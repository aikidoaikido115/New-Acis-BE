package repositories

import (
	"github.com/aikidoaikido115/New-Acis-BE/modules/entities"
	"github.com/aikidoaikido115/New-Acis-BE/modules/medicine/models"

	"gorm.io/gorm"
)

type GormDrugRepository struct {
	db *gorm.DB
}

func NewGormDrugRepository(db *gorm.DB) *GormDrugRepository {
	return &GormDrugRepository{
		db: db,
	}
}

type DrugRepository interface {

	// DrugMaster Operations
	CreateDrugMaster(drug *entities.DrugMaster) (*entities.DrugMaster, error)
	GetDrugMasterByID(id string) (*entities.DrugMaster, error)
	GetDrugMasterByName(name string) (*entities.DrugMaster, error)
	GetAllDrugMasters() ([]*entities.DrugMaster, error)
	UpdateDrugMaster(drug *entities.DrugMaster) (*entities.DrugMaster, error)
	DeleteDrugMaster(id string) error
	DrugMasterExistsByName(name string) (bool, error)
	DrugMasterExistsByNameAndDose(name string, dose string) (bool, error)

	// PersonalDrug Operations
	CreatePersonalDrug(personalDrug *entities.PersonalDrug) (*entities.PersonalDrug, error)
	GetPersonalDrugByID(id string) (*entities.PersonalDrug, error)
	GetAllPersonalDrugs() ([]*entities.PersonalDrug, error)
	GetPersonalDrugsToday() ([]*entities.PersonalDrug, error)
	GetPersonalDrugsByResidentID(residentID string) ([]*entities.PersonalDrug, error)
	GetPersonalDrugsByResidentIDToday(residentID string) ([]*entities.PersonalDrug, error)
	GetPersonalDrugsByTimeOfDayToday(timeOfDay string) ([]*entities.PersonalDrug, error)
	GetPersonalDrugsByTakeTypeToday(takeType string) ([]*entities.PersonalDrug, error)
	SearchPersonalDrugsTodayByResidentName(search string) ([]*entities.PersonalDrug, error)
	GetPersonalDrugsTodayCustom(timeOfDay *string, search *string, takeType *string) ([]*entities.PersonalDrug, error)
	ResidentExistsByID(id string) (bool, error)
	UpdatePersonalDrug(personalDrug *entities.PersonalDrug) (*entities.PersonalDrug, error)
	DeletePersonalDrug(id string) error

	// DrugPlan Operations
	CreateDrugPlan(drugPlan *entities.DrugPlan) (*entities.DrugPlan, error)
	GetDrugPlanByID(id string) (*entities.DrugPlan, error)
	GetAllDrugPlans() ([]*entities.DrugPlan, error)
	GetDrugPlansToday() ([]*entities.DrugPlan, error)
	GetDrugPlansByResidentID(residentID string) ([]*entities.DrugPlan, error)
	GetDrugPlansByResidentIDToday(residentID string) ([]*entities.DrugPlan, error)
	GetDrugPlansTodayCustom(timeOfDay *string, search *string, takeType *string) ([]*entities.DrugPlan, error)
	GetDrugPlansTodayResidentSummary() (*models.DrugPlanResidentSummaryResponse, error)
	UpdateDrugPlan(drugPlan *entities.DrugPlan) (*entities.DrugPlan, error)
	DeleteDrugPlan(id string) error
}

func (r *GormDrugRepository) CreateDrugMaster(drug *entities.DrugMaster) (*entities.DrugMaster, error) {
	if err := r.db.Create(&drug).Error; err != nil {
		return nil, err
	}

	return r.GetDrugMasterByID(drug.ID)
}

func (r *GormDrugRepository) GetDrugMasterByID(id string) (*entities.DrugMaster, error) {
	var drug entities.DrugMaster
	if err := r.db.Where("id = ?", id).First(&drug).Error; err != nil {
		return nil, err
	}

	return &drug, nil
}

func (r *GormDrugRepository) GetDrugMasterByName(name string) (*entities.DrugMaster, error) {
	var drug entities.DrugMaster
	if err := r.db.Where("name = ?", name).First(&drug).Error; err != nil {
		return nil, err
	}

	return &drug, nil
}

func (r *GormDrugRepository) GetAllDrugMasters() ([]*entities.DrugMaster, error) {
	var drugs []*entities.DrugMaster
	if err := r.db.Order("name ASC").Find(&drugs).Error; err != nil {
		return nil, err
	}

	return drugs, nil
}

func (r *GormDrugRepository) UpdateDrugMaster(drug *entities.DrugMaster) (*entities.DrugMaster, error) {
	if err := r.db.Save(&drug).Error; err != nil {
		return nil, err
	}

	return r.GetDrugMasterByID(drug.ID)
}

func (r *GormDrugRepository) DeleteDrugMaster(id string) error {
	if err := r.db.Delete(&entities.DrugMaster{}, "id = ?", id).Error; err != nil {
		return err
	}

	return nil
}

func (r *GormDrugRepository) DrugMasterExistsByName(name string) (bool, error) {
	var count int64
	err := r.db.Model(&entities.DrugMaster{}).Where("name = ?", name).Count(&count).Error
	if err != nil {
		return false, err
	}

	return count > 0, nil
}

func (r *GormDrugRepository) DrugMasterExistsByNameAndDose(name string, dose string) (bool, error) {
	var count int64
	err := r.db.Model(&entities.DrugMaster{}).Where("name = ? AND dose = ?", name, dose).Count(&count).Error
	if err != nil {
		return false, err
	}

	return count > 0, nil
}

func (r *GormDrugRepository) CreatePersonalDrug(personalDrug *entities.PersonalDrug) (*entities.PersonalDrug, error) {
	if err := r.db.Create(&personalDrug).Error; err != nil {
		return nil, err
	}

	return r.GetPersonalDrugByID(personalDrug.ID)
}

func (r *GormDrugRepository) GetPersonalDrugByID(id string) (*entities.PersonalDrug, error) {
	var personalDrug entities.PersonalDrug
	if err := r.db.Preload("Resident").Preload("DrugMaster").Where("id = ?", id).First(&personalDrug).Error; err != nil {
		return nil, err
	}

	return &personalDrug, nil
}

func (r *GormDrugRepository) GetAllPersonalDrugs() ([]*entities.PersonalDrug, error) {
	var personalDrugs []*entities.PersonalDrug
	if err := r.db.Preload("Resident").Preload("DrugMaster").Find(&personalDrugs).Error; err != nil {
		return nil, err
	}

	return personalDrugs, nil
}

func (r *GormDrugRepository) GetPersonalDrugsToday() ([]*entities.PersonalDrug, error) {
	var personalDrugs []*entities.PersonalDrug
	if err := r.db.
		Preload("Resident").
		Preload("DrugMaster").
		Where("personal_drugs.created_at >= DATE_TRUNC('day', NOW() AT TIME ZONE 'Asia/Bangkok') AT TIME ZONE 'Asia/Bangkok'").
		Where("personal_drugs.created_at < DATE_TRUNC('day', NOW() AT TIME ZONE 'Asia/Bangkok') AT TIME ZONE 'Asia/Bangkok' + INTERVAL '1 day'").
		Order("personal_drugs.created_at DESC").
		Find(&personalDrugs).Error; err != nil {
		return nil, err
	}

	return personalDrugs, nil
}

func (r *GormDrugRepository) GetPersonalDrugsByResidentID(residentID string) ([]*entities.PersonalDrug, error) {
	var personalDrugs []*entities.PersonalDrug
	if err := r.db.
		Preload("Resident").
		Preload("DrugMaster").
		Where("personal_drugs.resident_id = ?", residentID).
		Order("personal_drugs.created_at DESC").
		Find(&personalDrugs).Error; err != nil {
		return nil, err
	}

	return personalDrugs, nil
}

func (r *GormDrugRepository) GetPersonalDrugsByResidentIDToday(residentID string) ([]*entities.PersonalDrug, error) {
	var personalDrugs []*entities.PersonalDrug
	if err := r.db.
		Preload("Resident").
		Preload("DrugMaster").
		Where("personal_drugs.resident_id = ?", residentID).
		Where("personal_drugs.created_at >= DATE_TRUNC('day', NOW() AT TIME ZONE 'Asia/Bangkok') AT TIME ZONE 'Asia/Bangkok'").
		Where("personal_drugs.created_at < DATE_TRUNC('day', NOW() AT TIME ZONE 'Asia/Bangkok') AT TIME ZONE 'Asia/Bangkok' + INTERVAL '1 day'").
		Order("personal_drugs.created_at DESC").
		Find(&personalDrugs).Error; err != nil {
		return nil, err
	}

	return personalDrugs, nil
}

func (r *GormDrugRepository) GetPersonalDrugsByTimeOfDayToday(timeOfDay string) ([]*entities.PersonalDrug, error) {
	var personalDrugs []*entities.PersonalDrug
	if err := r.db.
		Preload("Resident").
		Preload("DrugMaster").
		Where("personal_drugs.time_of_day ILIKE ?", "%"+timeOfDay+"%").
		Where("personal_drugs.created_at >= DATE_TRUNC('day', NOW() AT TIME ZONE 'Asia/Bangkok') AT TIME ZONE 'Asia/Bangkok'").
		Where("personal_drugs.created_at < DATE_TRUNC('day', NOW() AT TIME ZONE 'Asia/Bangkok') AT TIME ZONE 'Asia/Bangkok' + INTERVAL '1 day'").
		Order("personal_drugs.created_at DESC").
		Find(&personalDrugs).Error; err != nil {
		return nil, err
	}

	return personalDrugs, nil
}

func (r *GormDrugRepository) GetPersonalDrugsByTakeTypeToday(takeType string) ([]*entities.PersonalDrug, error) {
	var personalDrugs []*entities.PersonalDrug
	if err := r.db.
		Preload("Resident").
		Preload("DrugMaster").
		Where("LOWER(personal_drugs.take_type) = LOWER(?)", takeType).
		Where("personal_drugs.created_at >= DATE_TRUNC('day', NOW() AT TIME ZONE 'Asia/Bangkok') AT TIME ZONE 'Asia/Bangkok'").
		Where("personal_drugs.created_at < DATE_TRUNC('day', NOW() AT TIME ZONE 'Asia/Bangkok') AT TIME ZONE 'Asia/Bangkok' + INTERVAL '1 day'").
		Order("personal_drugs.created_at DESC").
		Find(&personalDrugs).Error; err != nil {
		return nil, err
	}

	return personalDrugs, nil
}

func (r *GormDrugRepository) SearchPersonalDrugsTodayByResidentName(search string) ([]*entities.PersonalDrug, error) {
	like := "%" + search + "%"
	var personalDrugs []*entities.PersonalDrug
	if err := r.db.
		Preload("Resident").
		Preload("DrugMaster").
		Joins("JOIN residents ON personal_drugs.resident_id = residents.id").
		Where("(residents.first_name ILIKE ? OR residents.last_name ILIKE ? OR residents.nickname ILIKE ?)", like, like, like).
		Where("personal_drugs.created_at >= DATE_TRUNC('day', NOW() AT TIME ZONE 'Asia/Bangkok') AT TIME ZONE 'Asia/Bangkok'").
		Where("personal_drugs.created_at < DATE_TRUNC('day', NOW() AT TIME ZONE 'Asia/Bangkok') AT TIME ZONE 'Asia/Bangkok' + INTERVAL '1 day'").
		Order("personal_drugs.created_at DESC").
		Find(&personalDrugs).Error; err != nil {
		return nil, err
	}

	return personalDrugs, nil
}

func (r *GormDrugRepository) GetPersonalDrugsTodayCustom(timeOfDay *string, search *string, takeType *string) ([]*entities.PersonalDrug, error) {
	var personalDrugs []*entities.PersonalDrug

	query := r.db.
		Preload("Resident").
		Preload("DrugMaster").
		Where("personal_drugs.created_at >= DATE_TRUNC('day', NOW() AT TIME ZONE 'Asia/Bangkok') AT TIME ZONE 'Asia/Bangkok'").
		Where("personal_drugs.created_at < DATE_TRUNC('day', NOW() AT TIME ZONE 'Asia/Bangkok') AT TIME ZONE 'Asia/Bangkok' + INTERVAL '1 day'")

	if timeOfDay != nil && *timeOfDay != "" {
		query = query.Where("personal_drugs.time_of_day ILIKE ?", "%"+*timeOfDay+"%")
	}

	if takeType != nil && *takeType != "" {
		query = query.Where("LOWER(personal_drugs.take_type) = LOWER(?)", *takeType)
	}

	if search != nil && *search != "" {
		like := "%" + *search + "%"
		query = query.
			Joins("JOIN residents ON personal_drugs.resident_id = residents.id").
			Where("(residents.first_name ILIKE ? OR residents.last_name ILIKE ? OR residents.nickname ILIKE ?)", like, like, like)
	}

	if err := query.Order("personal_drugs.created_at DESC").Find(&personalDrugs).Error; err != nil {
		return nil, err
	}

	return personalDrugs, nil
}

func (r *GormDrugRepository) ResidentExistsByID(id string) (bool, error) {
	var count int64
	err := r.db.Model(&entities.Resident{}).Where("id = ?", id).Count(&count).Error
	if err != nil {
		return false, err
	}

	return count > 0, nil
}

func (r *GormDrugRepository) UpdatePersonalDrug(personalDrug *entities.PersonalDrug) (*entities.PersonalDrug, error) {
	if err := r.db.Save(&personalDrug).Error; err != nil {
		return nil, err
	}

	return r.GetPersonalDrugByID(personalDrug.ID)
}

func (r *GormDrugRepository) DeletePersonalDrug(id string) error {
	if err := r.db.Delete(&entities.PersonalDrug{}, "id = ?", id).Error; err != nil {
		return err
	}

	return nil
}

func (r *GormDrugRepository) CreateDrugPlan(drugPlan *entities.DrugPlan) (*entities.DrugPlan, error) {
	if err := r.db.Create(&drugPlan).Error; err != nil {
		return nil, err
	}

	return r.GetDrugPlanByID(drugPlan.ID)
}

func (r *GormDrugRepository) GetDrugPlanByID(id string) (*entities.DrugPlan, error) {
	var drugPlan entities.DrugPlan
	if err := r.db.
		Preload("PersonalDrug").
		Where("id = ?", id).
		First(&drugPlan).Error; err != nil {
		return nil, err
	}

	return &drugPlan, nil
}

func (r *GormDrugRepository) GetAllDrugPlans() ([]*entities.DrugPlan, error) {
	var drugPlans []*entities.DrugPlan
	if err := r.db.
		Preload("PersonalDrug").
		Order("created_at DESC").
		Find(&drugPlans).Error; err != nil {
		return nil, err
	}

	return drugPlans, nil
}

func (r *GormDrugRepository) GetDrugPlansToday() ([]*entities.DrugPlan, error) {
	var drugPlans []*entities.DrugPlan
	if err := r.db.
		Preload("PersonalDrug").
		Where("drug_plans.created_at >= DATE_TRUNC('day', NOW() AT TIME ZONE 'Asia/Bangkok') AT TIME ZONE 'Asia/Bangkok'").
		Where("drug_plans.created_at < DATE_TRUNC('day', NOW() AT TIME ZONE 'Asia/Bangkok') AT TIME ZONE 'Asia/Bangkok' + INTERVAL '1 day'").
		Order("drug_plans.created_at DESC").
		Find(&drugPlans).Error; err != nil {
		return nil, err
	}

	return drugPlans, nil
}

func (r *GormDrugRepository) GetDrugPlansByResidentID(residentID string) ([]*entities.DrugPlan, error) {
	var drugPlans []*entities.DrugPlan
	if err := r.db.
		Preload("PersonalDrug").
		Joins("JOIN personal_drugs ON drug_plans.pd_id = personal_drugs.id").
		Where("personal_drugs.resident_id = ?", residentID).
		Order("drug_plans.created_at DESC").
		Find(&drugPlans).Error; err != nil {
		return nil, err
	}

	return drugPlans, nil
}

func (r *GormDrugRepository) GetDrugPlansByResidentIDToday(residentID string) ([]*entities.DrugPlan, error) {
	var drugPlans []*entities.DrugPlan
	if err := r.db.
		Preload("PersonalDrug").
		Joins("JOIN personal_drugs ON drug_plans.pd_id = personal_drugs.id").
		Where("personal_drugs.resident_id = ?", residentID).
		Where("drug_plans.created_at >= DATE_TRUNC('day', NOW() AT TIME ZONE 'Asia/Bangkok') AT TIME ZONE 'Asia/Bangkok'").
		Where("drug_plans.created_at < DATE_TRUNC('day', NOW() AT TIME ZONE 'Asia/Bangkok') AT TIME ZONE 'Asia/Bangkok' + INTERVAL '1 day'").
		Order("drug_plans.created_at DESC").
		Find(&drugPlans).Error; err != nil {
		return nil, err
	}

	return drugPlans, nil
}

func (r *GormDrugRepository) GetDrugPlansTodayCustom(timeOfDay *string, search *string, takeType *string) ([]*entities.DrugPlan, error) {
	var drugPlans []*entities.DrugPlan

	query := r.db.
		Preload("PersonalDrug").
		Joins("JOIN personal_drugs ON drug_plans.pd_id = personal_drugs.id").
		Where("drug_plans.created_at >= DATE_TRUNC('day', NOW() AT TIME ZONE 'Asia/Bangkok') AT TIME ZONE 'Asia/Bangkok'").
		Where("drug_plans.created_at < DATE_TRUNC('day', NOW() AT TIME ZONE 'Asia/Bangkok') AT TIME ZONE 'Asia/Bangkok' + INTERVAL '1 day'")

	if timeOfDay != nil && *timeOfDay != "" {
		query = query.Where("personal_drugs.time_of_day ILIKE ?", "%"+*timeOfDay+"%")
	}

	if takeType != nil && *takeType != "" {
		query = query.Where("LOWER(personal_drugs.take_type) = LOWER(?)", *takeType)
	}

	if search != nil && *search != "" {
		like := "%" + *search + "%"
		query = query.
			Joins("JOIN residents ON personal_drugs.resident_id = residents.id").
			Where("(residents.first_name ILIKE ? OR residents.last_name ILIKE ? OR residents.nickname ILIKE ?)", like, like, like)
	}

	if err := query.Order("drug_plans.created_at DESC").Find(&drugPlans).Error; err != nil {
		return nil, err
	}

	return drugPlans, nil
}

func (r *GormDrugRepository) GetDrugPlansTodayResidentSummary() (*models.DrugPlanResidentSummaryResponse, error) {
	var summary models.DrugPlanResidentSummaryResponse

	err := r.db.Raw(`
		SELECT
			COUNT(*) AS total_residents,
			COUNT(*) FILTER (WHERE s.waiting_count > 0) AS waiting_residents,
			COUNT(*) FILTER (WHERE s.waiting_count = 0 AND s.taken_count > 0) AS given_residents
		FROM (
			SELECT
				pd.resident_id,
				SUM(CASE WHEN dp.is_taken = FALSE THEN 1 ELSE 0 END) AS waiting_count,
				SUM(CASE WHEN dp.is_taken = TRUE THEN 1 ELSE 0 END) AS taken_count
			FROM drug_plans dp
			JOIN personal_drugs pd ON dp.pd_id = pd.id
			WHERE dp.created_at >= DATE_TRUNC('day', NOW() AT TIME ZONE 'Asia/Bangkok') AT TIME ZONE 'Asia/Bangkok'
				AND dp.created_at < DATE_TRUNC('day', NOW() AT TIME ZONE 'Asia/Bangkok') AT TIME ZONE 'Asia/Bangkok' + INTERVAL '1 day'
			GROUP BY pd.resident_id
		) s
	`).Scan(&summary).Error
	if err != nil {
		return nil, err
	}

	return &summary, nil
}

func (r *GormDrugRepository) UpdateDrugPlan(drugPlan *entities.DrugPlan) (*entities.DrugPlan, error) {
	if err := r.db.Save(&drugPlan).Error; err != nil {
		return nil, err
	}

	return r.GetDrugPlanByID(drugPlan.ID)
}

func (r *GormDrugRepository) DeleteDrugPlan(id string) error {
	if err := r.db.Delete(&entities.DrugPlan{}, "id = ?", id).Error; err != nil {
		return err
	}

	return nil
}
