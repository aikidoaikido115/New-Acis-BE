package repositories

import (
	"time"

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
	GetPersonalDrugsTodayCustom(timeOfDay *string, search *string, takeType *string, page int, pageSize int) ([]*entities.PersonalDrug, int64, error)
	GetActivePersonalDrugsForDate(date time.Time, residentID *string) ([]*entities.PersonalDrug, error)
	GetExpiredAsNeededPersonalDrugs(date time.Time, residentID *string) ([]*entities.PersonalDrug, error)
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
	GetDrugPlansTodayCustom(timeOfDay *string, search *string, takeType *string, page int, pageSize int) ([]*entities.DrugPlan, int64, error)
	GetDrugAdministrationHistory(req models.DrugAdministrationHistoryQueryParams, page int, pageSize int) ([]models.DrugAdministrationHistoryItem, int64, error)
	GetDrugPlansTodayResidentSummary() (*models.DrugPlanResidentSummaryResponse, error)
	HasDrugPlanForPersonalDrugOnDate(pdID string, date time.Time) (bool, error)
	DeleteDrugPlansByPdID(pdID string) error
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

func (r *GormDrugRepository) GetPersonalDrugsTodayCustom(timeOfDay *string, search *string, takeType *string, page int, pageSize int) ([]*entities.PersonalDrug, int64, error) {
	var personalDrugs []*entities.PersonalDrug

	applyFilters := func(query *gorm.DB) *gorm.DB {
		query = query.
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

		return query
	}

	var total int64
	if err := applyFilters(r.db.Model(&entities.PersonalDrug{})).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	query := applyFilters(
		r.db.
			Preload("Resident").
			Preload("DrugMaster").
			Model(&entities.PersonalDrug{}),
	)

	offset := (page - 1) * pageSize

	if err := query.Order("personal_drugs.created_at DESC").Offset(offset).Limit(pageSize).Find(&personalDrugs).Error; err != nil {
		return nil, 0, err
	}

	return personalDrugs, total, nil
}

func (r *GormDrugRepository) GetActivePersonalDrugsForDate(date time.Time, residentID *string) ([]*entities.PersonalDrug, error) {
	var personalDrugs []*entities.PersonalDrug

	loc, err := time.LoadLocation("Asia/Bangkok")
	if err != nil {
		return nil, err
	}

	dayStart := time.Date(date.In(loc).Year(), date.In(loc).Month(), date.In(loc).Day(), 0, 0, 0, 0, loc)
	dayEnd := dayStart.AddDate(0, 0, 1)
	dayDate := dayStart.Format("2006-01-02")

	query := r.db.
		Preload("Resident").
		Preload("DrugMaster").
		Where("personal_drugs.created_at < ?", dayEnd).
		Where("(LOWER(personal_drugs.take_type) = 'regular' OR (LOWER(personal_drugs.take_type) = 'as_needed' AND personal_drugs.start_date IS NOT NULL AND personal_drugs.end_date IS NOT NULL AND personal_drugs.start_date <= ?::date AND personal_drugs.end_date >= ?::date))", dayDate, dayDate)

	if residentID != nil && *residentID != "" {
		query = query.Where("personal_drugs.resident_id = ?", *residentID)
	}

	if err := query.Order("personal_drugs.created_at ASC").Find(&personalDrugs).Error; err != nil {
		return nil, err
	}

	return personalDrugs, nil
}

func (r *GormDrugRepository) GetExpiredAsNeededPersonalDrugs(date time.Time, residentID *string) ([]*entities.PersonalDrug, error) {
	var personalDrugs []*entities.PersonalDrug

	loc, err := time.LoadLocation("Asia/Bangkok")
	if err != nil {
		return nil, err
	}

	dayDate := time.Date(date.In(loc).Year(), date.In(loc).Month(), date.In(loc).Day(), 0, 0, 0, 0, loc).Format("2006-01-02")

	query := r.db.
		Preload("Resident").
		Preload("DrugMaster").
		Where("LOWER(personal_drugs.take_type) = 'as_needed'").
		Where("personal_drugs.end_date IS NOT NULL").
		Where("personal_drugs.end_date < ?::date", dayDate)

	if residentID != nil && *residentID != "" {
		query = query.Where("personal_drugs.resident_id = ?", *residentID)
	}

	if err := query.Order("personal_drugs.updated_at ASC").Find(&personalDrugs).Error; err != nil {
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
		Preload("PersonalDrug.Resident").
		Preload("PersonalDrug.DrugMaster").
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
		Preload("PersonalDrug.Resident").
		Preload("PersonalDrug.DrugMaster").
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
		Preload("PersonalDrug.Resident").
		Preload("PersonalDrug.DrugMaster").
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
		Preload("PersonalDrug.Resident").
		Preload("PersonalDrug.DrugMaster").
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
		Preload("PersonalDrug.Resident").
		Preload("PersonalDrug.DrugMaster").
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

func (r *GormDrugRepository) GetDrugPlansTodayCustom(timeOfDay *string, search *string, takeType *string, page int, pageSize int) ([]*entities.DrugPlan, int64, error) {
	var drugPlans []*entities.DrugPlan

	applyFilters := func(query *gorm.DB) *gorm.DB {
		query = query.
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

		return query
	}

	var total int64
	if err := applyFilters(r.db.Model(&entities.DrugPlan{})).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	query := applyFilters(
		r.db.
			Preload("PersonalDrug").
			Preload("PersonalDrug.Resident").
			Preload("PersonalDrug.DrugMaster").
			Model(&entities.DrugPlan{}),
	)

	offset := (page - 1) * pageSize

	if err := query.Order("drug_plans.created_at DESC").Offset(offset).Limit(pageSize).Find(&drugPlans).Error; err != nil {
		return nil, 0, err
	}

	return drugPlans, total, nil
}

func (r *GormDrugRepository) GetDrugAdministrationHistory(req models.DrugAdministrationHistoryQueryParams, page int, pageSize int) ([]models.DrugAdministrationHistoryItem, int64, error) {
	items := make([]models.DrugAdministrationHistoryItem, 0)

	query := r.db.Table("drug_plans dp").
		Joins("JOIN personal_drugs pd ON dp.pd_id = pd.id").
		Joins("JOIN drug_masters dm ON pd.dm_id = dm.id").
		Joins("JOIN residents r ON pd.resident_id = r.id").
		Joins("LEFT JOIN staffs s ON dp.given_by_staff_id = s.id").
		Joins("LEFT JOIN users su ON s.user_id = su.id")

	if req.Date != nil && *req.Date != "" {
		query = query.Where("DATE(dp.created_at AT TIME ZONE 'Asia/Bangkok') = ?::date", *req.Date)
	} else {
		query = query.
			Where("dp.created_at >= DATE_TRUNC('day', NOW() AT TIME ZONE 'Asia/Bangkok') AT TIME ZONE 'Asia/Bangkok'").
			Where("dp.created_at < DATE_TRUNC('day', NOW() AT TIME ZONE 'Asia/Bangkok') AT TIME ZONE 'Asia/Bangkok' + INTERVAL '1 day'")
	}

	if req.TimeOfDay != nil && *req.TimeOfDay != "" {
		query = query.Where("pd.time_of_day ILIKE ?", "%"+*req.TimeOfDay+"%")
	}

	if req.Status != nil && *req.Status != "" {
		switch *req.Status {
		case "taken":
			query = query.Where("dp.is_taken = TRUE")
		case "omitted":
			query = query.Where("dp.is_omitted = TRUE")
		case "pending":
			query = query.Where("dp.is_taken = FALSE AND dp.is_omitted = FALSE")
		}
	}

	if req.Search != nil && *req.Search != "" {
		like := "%" + *req.Search + "%"
		query = query.Where("(r.first_name ILIKE ? OR r.last_name ILIKE ? OR r.nickname ILIKE ?)", like, like, like)
	}

	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	if err := query.
		Select(`
			dp.id AS drug_plan_id,
			COALESCE(dp.taken_at, dp.updated_at, dp.created_at) AS action_at,
			TRIM(CONCAT(COALESCE(r.first_name, ''), ' ', COALESCE(r.last_name, ''))) AS resident_name,
			dm.name AS drug_name,
			dm.dose AS drug_dose,
			CASE
				WHEN dp.is_taken = TRUE THEN 'taken'
				WHEN dp.is_omitted = TRUE THEN 'omitted'
				ELSE 'pending'
			END AS status,
			CASE
				WHEN dp.notes IS NOT NULL AND TRIM(dp.notes) <> '' THEN dp.notes
				WHEN dp.is_omitted = TRUE THEN dp.omitted_reason
				ELSE NULL
			END AS note,
			NULLIF(TRIM(CONCAT(COALESCE(su.first_name, ''), ' ', COALESCE(su.last_name, ''))), '') AS given_by_staff_name,
			pd.time_of_day AS time_of_day
		`).
		Order("COALESCE(dp.taken_at, dp.updated_at, dp.created_at) DESC").
		Offset(offset).
		Limit(pageSize).
		Scan(&items).Error; err != nil {
		return nil, 0, err
	}

	return items, total, nil
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

func (r *GormDrugRepository) HasDrugPlanForPersonalDrugOnDate(pdID string, date time.Time) (bool, error) {
	loc, err := time.LoadLocation("Asia/Bangkok")
	if err != nil {
		return false, err
	}

	dayStart := time.Date(date.In(loc).Year(), date.In(loc).Month(), date.In(loc).Day(), 0, 0, 0, 0, loc)
	dayEnd := dayStart.AddDate(0, 0, 1)

	var count int64
	err = r.db.Model(&entities.DrugPlan{}).
		Where("pd_id = ?", pdID).
		Where("created_at >= ?", dayStart).
		Where("created_at < ?", dayEnd).
		Count(&count).Error
	if err != nil {
		return false, err
	}

	return count > 0, nil
}

func (r *GormDrugRepository) DeleteDrugPlansByPdID(pdID string) error {
	if err := r.db.Delete(&entities.DrugPlan{}, "pd_id = ?", pdID).Error; err != nil {
		return err
	}

	return nil
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
