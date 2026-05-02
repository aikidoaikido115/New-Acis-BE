package repositories

import (
	"time"

	emr_constants "github.com/aikidoaikido115/New-Acis-BE/modules/emr/constants"
	"github.com/aikidoaikido115/New-Acis-BE/modules/emr/models"
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
	UpdateResident(resident *entities.Resident) (*entities.Resident, error)
	ResidentExists(id string) (bool, error)
	IdCardNumberExists(idCardNumber string) (bool, error)

	// Dashboard operations
	GetNumberOfResidentsDashboard() (models.NumberOfResidentsDashboardResponse, error)
	GetNumberOfResidentGender() (models.ResidentGenderStatsDashboardResponse, error)
	GetResidentAllergyStatsDashboard() (models.ResidentAllergyStatsDashboardResponse, error)
	GetResidentDrugAllergyStatsDashboard() (models.ResidentDrugAllergyStatsDashboardResponse, error)

	// Room operations
	RoomExists(id string) (bool, error)
	GetRoomByID(id string) (*entities.Room, error)
	GetAllRooms() ([]*entities.Room, error)
	CreateRoom(room *entities.Room) (*entities.Room, error)
	UpdateRoom(room *entities.Room) (*entities.Room, error)
	RoomNumberExists(roomNumber string) (bool, error)

	GetResidentsCustom(params models.ResidentQueryParams) ([]*entities.Resident, int64, error)

	// IntakeLabel operations
	CreateIntakeLabel(label *entities.IntakeLabels) (*entities.IntakeLabels, error)
	GetIntakeLabelByID(id string) (*entities.IntakeLabels, error)
	GetIntakeLabelByName(labelName string) (*entities.IntakeLabels, error)
	GetAllIntakeLabels() ([]*entities.IntakeLabels, error)
	LabelExists(labelName string) (bool, error)
	DeleteIntakeLabel(id string) error

	// Allergy operations
	CreateAllergy(allergy *entities.Allergy) (*entities.Allergy, error)
	GetAllergyByID(id string) (*entities.Allergy, error)
	GetAllergyByName(allergyName string) (*entities.Allergy, error)
	GetAllAllergies() ([]*entities.Allergy, error)
	AllergyExists(allergyName string) (bool, error)
	DeleteAllergy(id string) error

	// DrugAllergy operations
	CreateDrugAllergy(drugAllergy *entities.DrugAllergy) (*entities.DrugAllergy, error)
	GetDrugAllergyByID(id string) (*entities.DrugAllergy, error)
	GetDrugAllergyByName(allergyName string) (*entities.DrugAllergy, error)
	GetAllDrugAllergies() ([]*entities.DrugAllergy, error)
	DrugAllergyExists(allergyName string) (bool, error)
	DeleteDrugAllergy(id string) error

	// ResidentLabel operations (many-to-many)
	CreateIntakeLabelByResidentID(residentLabel *entities.ResidentLabels) (*entities.ResidentLabels, error)
	GetResidentLabelsByResidentID(residentID string) ([]*entities.ResidentLabels, error)
	ResidentLabelExists(residentID, labelID string) (bool, error)
	DeleteResidentLabelsByResidentID(residentID string) error

	// ResidentAllergy operations (many-to-many)
	CreateAllergyByResidentID(residentAllergy *entities.ResidentAllergies) (*entities.ResidentAllergies, error)
	GetResidentAllergiesByResidentID(residentID string) ([]*entities.ResidentAllergies, error)
	GetAllResidentAllergies() ([]*models.ResidentAllergyListResponse, error)
	ResidentAllergyExists(residentID, allergyID string) (bool, error)
	DeleteResidentAllergiesByResidentID(residentID string) error

	// ResidentDA operations (many-to-many)
	CreateDrugAllergyByResidentID(residentDA *entities.ResidentDA) (*entities.ResidentDA, error)
	GetResidentDrugAllergiesByResidentID(residentID string) ([]*entities.ResidentDA, error)
	GetAllResidentDrugAllergies() ([]*models.ResidentDrugAllergyListResponse, error)
	ResidentDrugAllergyExists(residentID, drugAllergyID string) (bool, error)
	DeleteResidentDrugAllergiesByResidentID(residentID string) error

	// VitalSign operations
	CreateVitalSign(vitalSign *entities.VitalSign) (*entities.VitalSign, error)
	VitalSignSlotExists(residentID string, measurementDate time.Time, timeOfDay string) (bool, error)

	GetVitalSignByID(id string) (*entities.VitalSign, error)
	GetVitalSignsByResidentIDOnDate(residentID string, dayDate time.Time, isLatest bool) ([]*entities.VitalSign, error)
	GetVitalSignsOnDate(dayDate time.Time, isLatest bool) ([]*entities.VitalSign, error)
	GetVitalSignsByRoomIDToday(roomID string, isLatest bool) ([]*entities.VitalSign, error)
	GetVitalSignsByFloorToday(floor int16, isLatest bool) ([]*entities.VitalSign, error)
	GetVitalSignsByResidentIDToday(residentID string, isLatest bool) ([]*entities.VitalSign, error)
	GetVitalSignsHistory(residentID string) ([]*entities.VitalSign, error)
	GetVitalSignsToday(isLatest bool) ([]*entities.VitalSign, error)
	GetVitalSignsCustom(params models.VitalSignQueryParams) ([]*entities.VitalSign, int64, error)

	UpdateVitalSignByID(vitalSign *entities.VitalSign) (*entities.VitalSign, error)

	//LaboratoryValue operations
	CreateLaboratoryValue(laboratoryValue *entities.LaboratoryValue) (*entities.LaboratoryValue, error)

	GetLaboratoryValueByID(id string) (*entities.LaboratoryValue, error)
	GetLaboratoryValuesByRoomIDToday(roomID string, isLatest bool) ([]*entities.LaboratoryValue, error)
	GetLaboratoryValuesByFloorToday(floor int16, isLatest bool) ([]*entities.LaboratoryValue, error)
	GetLaboratoryValuesByResidentIDToday(residentID string, isLatest bool) ([]*entities.LaboratoryValue, error)
	GetLaboratoryValuesHistory(residentID string) ([]*entities.LaboratoryValue, error)
	GetLaboratoryValuesToday(isLatest bool) ([]*entities.LaboratoryValue, error)
	GetLaboratoryValuesCustom(params models.LaboratoryValueQueryParams) ([]*entities.LaboratoryValue, int64, error)

	LaboratoryValueSlotExists(residentID string, measurementDate time.Time, timeOfDay string) (bool, error)
	GetLaboratoryValuesByResidentIDOnDate(residentID string, dayDate time.Time, isLatest bool) ([]*entities.LaboratoryValue, error)
	GetLaboratoryValuesOnDate(dayDate time.Time, isLatest bool) ([]*entities.LaboratoryValue, error)

	UpdateLaboratoryValueByID(laboratoryValue *entities.LaboratoryValue) (*entities.LaboratoryValue, error)

	GetUrineOutputSumGroupByResident(params models.LaboratoryValueQueryParams, urineType string) (*models.UrineOutputSumResponse, error)

	// NurseNote operations
	CreateNurseNote(note *entities.NurseNote) (*entities.NurseNote, error)
	GetNurseNoteByID(id string) (*entities.NurseNote, error)
	GetNurseNotesOverviewOnDate(dayDate time.Time) ([]*entities.NurseNote, error)
	GetNurseNotesByResidentIDOnDate(residentID string, dayDate time.Time) ([]*entities.NurseNote, error)
	UpdateNurseNoteByID(note *entities.NurseNote) (*entities.NurseNote, error)
	DeleteNurseNoteByID(id string) error

	// WoundCareNote operations
	CreateWoundCareNote(note *entities.WoundCareNote) (*entities.WoundCareNote, error)
	GetWoundCareNoteByID(id string) (*entities.WoundCareNote, error)
	GetWoundCareNotesOverviewOnDate(dayDate time.Time) ([]*entities.WoundCareNote, error)
	GetWoundCareNotesByResidentIDOnDate(residentID string, dayDate time.Time) ([]*entities.WoundCareNote, error)
	UpdateWoundCareNoteByID(note *entities.WoundCareNote) (*entities.WoundCareNote, error)
	DeleteWoundCareNoteByID(id string) error

	// RelativeNote operations
	CreateRelativeNote(note *entities.RelativeNote) (*entities.RelativeNote, error)
	GetRelativeNoteByID(id string) (*entities.RelativeNote, error)
	GetRelativeNotesOverviewOnDate(dayDate time.Time) ([]*entities.RelativeNote, error)
	GetRelativeNotesByResidentIDOnDate(residentID string, dayDate time.Time) ([]*entities.RelativeNote, error)
	UpdateRelativeNoteByID(note *entities.RelativeNote) (*entities.RelativeNote, error)
	DeleteRelativeNoteByID(id string) error

	// DoctorOrder operations
	CreateDoctorOrder(order *entities.DoctorOrder) (*entities.DoctorOrder, error)
	GetDoctorOrderByID(id string) (*entities.DoctorOrder, error)
	GetDoctorOrdersOverviewOnDate(dayDate time.Time) ([]*entities.DoctorOrder, error)
	GetDoctorOrdersByResidentIDOnDate(residentID string, dayDate time.Time) ([]*entities.DoctorOrder, error)
	UpdateDoctorOrderByID(order *entities.DoctorOrder) (*entities.DoctorOrder, error)
	DeleteDoctorOrderByID(id string) error

	// todo เพราะมันเจาะจงว่า ค่าไหนของ vital sign อีกทีนึง
	// GetLatestVitalSignsGreaterThanCustom(params models.VitalSignQueryParams, greaterThan float64) ([]*entities.VitalSign, error)
	// GetLatestVitalSignsLessThanCustom(params models.VitalSignQueryParams, lessThan float64) ([]*entities.VitalSign, error)

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
	if err := r.db.Preload("Room").Preload("ResidentLabels.IntakeLabel").Preload("ResidentAllergies.Allergy").Preload("ResidentDA.DrugAllergy").Where("id = ?", id).First(&resident).Error; err != nil {
		return nil, err
	}
	return &resident, nil
}

func (r *GormEmrRepository) GetResidentByRoomID(roomID string) ([]*entities.Resident, error) {
	var residents []*entities.Resident
	if err := r.db.Preload("Room").Preload("ResidentLabels.IntakeLabel").Preload("ResidentAllergies.Allergy").Preload("ResidentDA.DrugAllergy").Where("room_id = ?", roomID).Find(&residents).Error; err != nil {
		return nil, err
	}
	return residents, nil
}

func (r *GormEmrRepository) GetAllResidents() ([]*entities.Resident, error) {
	var residents []*entities.Resident
	if err := r.db.Preload("Room").Preload("ResidentLabels.IntakeLabel").Preload("ResidentAllergies.Allergy").Preload("ResidentDA.DrugAllergy").Find(&residents).Error; err != nil {
		return nil, err
	}
	return residents, nil
}

func (r *GormEmrRepository) GetResidentsCustom(params models.ResidentQueryParams) ([]*entities.Resident, int64, error) {
	var residents []*entities.Resident
	applyFilters := func(query *gorm.DB) *gorm.DB {
		needRoomsJoin := false
		needLabelsJoin := false

		if params.Floor != nil {
			needRoomsJoin = true
			query = query.Where("rooms.floor = ?", *params.Floor)
		}

		if len(params.LabelIDs) > 0 {
			needLabelsJoin = true
			query = query.Where("resident_labels.label_id IN ?", params.LabelIDs).
				Group("residents.id").
				Having("COUNT(DISTINCT resident_labels.label_id) = ?", len(params.LabelIDs))
		}

		if params.Status != nil && *params.Status != "" {
			query = query.Where("residents.status = ?", *params.Status)
		}

		if params.Search != nil && *params.Search != "" {
			like := "%" + *params.Search + "%"
			query = query.Where(
				"residents.first_name ILIKE ? OR residents.last_name ILIKE ? OR residents.nickname ILIKE ?",
				like, like, like,
			)
		}

		if needRoomsJoin {
			query = query.Joins("JOIN rooms ON residents.room_id = rooms.id")
		}
		if needLabelsJoin {
			query = query.Joins("JOIN resident_labels ON residents.id = resident_labels.resident_id")
		}

		return query
	}

	countBase := applyFilters(r.db.Model(&entities.Resident{}))
	countSubQuery := countBase.Select("residents.id")
	var total int64
	if err := r.db.Table("(?) AS filtered_residents", countSubQuery).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	query := applyFilters(
		r.db.Preload("Room").Preload("ResidentLabels.IntakeLabel").Model(&entities.Resident{}),
	)

	query = query.Order("residents.check_in_date DESC")

	if params.Limit > 0 {
		query = query.Limit(params.Limit)
	}
	if params.Offset > 0 {
		query = query.Offset(params.Offset)
	}

	if err := query.Find(&residents).Error; err != nil {
		return nil, 0, err
	}
	return residents, total, nil
}

func (r *GormEmrRepository) UpdateResident(resident *entities.Resident) (*entities.Resident, error) {
	if err := r.db.Save(&resident).Error; err != nil {
		return nil, err
	}
	return r.GetResidentByID(resident.ID)
}

func (r *GormEmrRepository) ResidentExists(id string) (bool, error) {
	var count int64
	err := r.db.Model(&entities.Resident{}).Where("id = ?", id).Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r *GormEmrRepository) IdCardNumberExists(idCardNumber string) (bool, error) {
	var count int64
	err := r.db.Model(&entities.Resident{}).Where("id_card_number = ?", idCardNumber).Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
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

func (r *GormEmrRepository) CreateRoom(room *entities.Room) (*entities.Room, error) {
	if err := r.db.Create(&room).Error; err != nil {
		return nil, err
	}
	return room, nil
}

func (r *GormEmrRepository) UpdateRoom(room *entities.Room) (*entities.Room, error) {
	if err := r.db.Save(&room).Error; err != nil {
		return nil, err
	}
	return room, nil
}

func (r *GormEmrRepository) RoomNumberExists(roomNumber string) (bool, error) {
	var count int64
	err := r.db.Model(&entities.Room{}).Where("room_number = ?", roomNumber).Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r *GormEmrRepository) GetNumberOfResidentsDashboard() (models.NumberOfResidentsDashboardResponse, error) {
	var response models.NumberOfResidentsDashboardResponse

	err := r.db.Model(&entities.Resident{}).
		Joins("LEFT JOIN resident_labels ON residents.id = resident_labels.resident_id").
		Joins("LEFT JOIN intake_labels ON resident_labels.label_id = intake_labels.id").
		Where("residents.status = ?", "active").
		Select(`
            COUNT(DISTINCT CASE WHEN intake_labels.label_name = ? THEN residents.id END) as independent_residents,
            COUNT(DISTINCT CASE WHEN intake_labels.label_name = ? THEN residents.id END) as partial_assist_residents,
            COUNT(DISTINCT CASE WHEN intake_labels.label_name = ? THEN residents.id END) as bedridden_residents,
            COUNT(DISTINCT residents.id) as total_residents
        `, emr_constants.CareLevelIndependent, emr_constants.CareLevelPartial, emr_constants.CareLevelFull).
		Scan(&response).Error
	if err != nil {
		return models.NumberOfResidentsDashboardResponse{}, err
	}

	return response, nil
}

func (r *GormEmrRepository) GetNumberOfResidentGender() (models.ResidentGenderStatsDashboardResponse, error) {
	var response models.ResidentGenderStatsDashboardResponse

	err := r.db.Model(&entities.Resident{}).
		Select(`
			COUNT(CASE WHEN gender = 'male' THEN 1 END) as sum_of_male,
			COUNT(CASE WHEN gender = 'female' THEN 1 END) as sum_of_female,
			COUNT(*) as total_residents
		`).Scan(&response).Error

	if err != nil {
		return models.ResidentGenderStatsDashboardResponse{}, err
	}
	return response, nil
}

func (r *GormEmrRepository) GetResidentAllergyStatsDashboard() (models.ResidentAllergyStatsDashboardResponse, error) {
	var response models.ResidentAllergyStatsDashboardResponse

	// Count allergic active residents (those with at least one allergy)
	if err := r.db.Model(&entities.ResidentAllergies{}).
		Joins("JOIN residents ON resident_allergies.resident_id = residents.id").
		Where("residents.status = ?", "active").
		Distinct("resident_allergies.resident_id").
		Count(&response.TotalAllergic).Error; err != nil {
		return models.ResidentAllergyStatsDashboardResponse{}, err
	}

	// Count non-allergic active residents (those with no allergies)
	if err := r.db.Model(&entities.Resident{}).
		Joins("LEFT JOIN resident_allergies ON residents.id = resident_allergies.resident_id").
		Where("residents.status = ? AND resident_allergies.resident_id IS NULL", "active").
		Distinct("residents.id").
		Count(&response.TotalNotAllergic).Error; err != nil {
		return models.ResidentAllergyStatsDashboardResponse{}, err
	}

	// Get allergy details for active residents grouped by allergy combination per resident
	// First, get all allergies per active resident, then group by combination
	var allergyGroupings []struct {
		AllergyNames string
		ResidentCount int64
	}
	
	if err := r.db.Model(&entities.ResidentAllergies{}).
		Joins("JOIN residents ON resident_allergies.resident_id = residents.id").
		Joins("JOIN allergies ON resident_allergies.allergy_id = allergies.id").
		Where("residents.status = ?", "active").
		Select("STRING_AGG(DISTINCT allergies.allergy_name, ' + ' ORDER BY allergies.allergy_name) as allergy_names, COUNT(DISTINCT resident_allergies.resident_id) as resident_count").
		Group("resident_allergies.resident_id").
		Scan(&allergyGroupings).Error; err != nil {
		return models.ResidentAllergyStatsDashboardResponse{}, err
	}

	// Transform results into the response format
	allergyCountMap := make(map[string]int64)
	for _, group := range allergyGroupings {
		allergyCountMap[group.AllergyNames] += group.ResidentCount
	}

	for allergyCombo, count := range allergyCountMap {
		response.AllergyDetails = append(response.AllergyDetails, models.AllergyStatisticDashboardResponse{
			AllergyID:     "",
			AllergyName:   allergyCombo,
			ResidentCount: count,
		})
	}

	return response, nil
}

func (r *GormEmrRepository) GetResidentDrugAllergyStatsDashboard() (models.ResidentDrugAllergyStatsDashboardResponse, error) {
	var response models.ResidentDrugAllergyStatsDashboardResponse

	if err := r.db.Model(&entities.ResidentDA{}).
		Distinct("resident_id").
		Count(&response.TotalDrugAllergic).Error; err != nil {
		return models.ResidentDrugAllergyStatsDashboardResponse{}, err
	}

	if err := r.db.Model(&entities.Resident{}).
		Joins("LEFT JOIN resident_das ON residents.id = resident_das.resident_id").
		Where("resident_das.resident_id IS NULL").
		Distinct("residents.id").
		Count(&response.TotalNotDrugAllergic).Error; err != nil {
		return models.ResidentDrugAllergyStatsDashboardResponse{}, err
	}

	if err := r.db.Model(&entities.ResidentDA{}).
		Joins("JOIN drug_allergies ON resident_das.drug_allergy_id = drug_allergies.id").
		Select("drug_allergies.id as drug_allergy_id, drug_allergies.allergy_name as allergy_name, COUNT(DISTINCT resident_das.resident_id) as count").
		Group("drug_allergies.id, drug_allergies.allergy_name").
		Order("count DESC").
		Scan(&response.DrugAllergyDetails).Error; err != nil {
		return models.ResidentDrugAllergyStatsDashboardResponse{}, err
	}

	return response, nil
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

// func (r *GormEmrRepository) UpdateIntakeLabelByResidentID(residentLabel *entities.ResidentLabels) (*entities.ResidentLabels, error) {
// 	if err := r.db.Save(&residentLabel).Error; err != nil {
// 		return nil, err
// 	}

// 	var result entities.ResidentLabels
// 	if err := r.db.Preload("Resident").Preload("IntakeLabel").
// 		Where("resident_id = ? AND label_id = ?", residentLabel.ResidentID, residentLabel.LabelID).
// 		First(&result).Error; err != nil {
// 		return nil, err
// 	}
// 	return &result, nil
// }

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

func (r *GormEmrRepository) LabelExists(labelName string) (bool, error) {
	var count int64
	err := r.db.Model(&entities.IntakeLabels{}).Where("label_name = ?", labelName).Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r *GormEmrRepository) DeleteIntakeLabel(id string) error {
	if err := r.db.Delete(&entities.IntakeLabels{}, "id = ?", id).Error; err != nil {
		return err
	}
	return nil
}

func (r *GormEmrRepository) GetAllergyByName(allergyName string) (*entities.Allergy, error) {
	var allergy entities.Allergy
	if err := r.db.Where("allergy_name = ?", allergyName).First(&allergy).Error; err != nil {
		return nil, err
	}
	return &allergy, nil
}

func (r *GormEmrRepository) CreateAllergy(allergy *entities.Allergy) (*entities.Allergy, error) {
	if err := r.db.Create(&allergy).Error; err != nil {
		return nil, err
	}
	return allergy, nil
}

func (r *GormEmrRepository) GetAllergyByID(id string) (*entities.Allergy, error) {
	var allergy entities.Allergy
	if err := r.db.Where("id = ?", id).First(&allergy).Error; err != nil {
		return nil, err
	}
	return &allergy, nil
}

func (r *GormEmrRepository) GetAllAllergies() ([]*entities.Allergy, error) {
	var allergies []*entities.Allergy
	if err := r.db.Find(&allergies).Error; err != nil {
		return nil, err
	}
	return allergies, nil
}

func (r *GormEmrRepository) AllergyExists(allergyName string) (bool, error) {
	var count int64
	err := r.db.Model(&entities.Allergy{}).Where("allergy_name = ?", allergyName).Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r *GormEmrRepository) DeleteAllergy(id string) error {
	if err := r.db.Delete(&entities.Allergy{}, "id = ?", id).Error; err != nil {
		return err
	}
	return nil
}

func (r *GormEmrRepository) GetDrugAllergyByName(allergyName string) (*entities.DrugAllergy, error) {
	var drugAllergy entities.DrugAllergy
	if err := r.db.Where("allergy_name = ?", allergyName).First(&drugAllergy).Error; err != nil {
		return nil, err
	}
	return &drugAllergy, nil
}

func (r *GormEmrRepository) CreateDrugAllergy(drugAllergy *entities.DrugAllergy) (*entities.DrugAllergy, error) {
	if err := r.db.Create(&drugAllergy).Error; err != nil {
		return nil, err
	}
	return drugAllergy, nil
}

func (r *GormEmrRepository) GetDrugAllergyByID(id string) (*entities.DrugAllergy, error) {
	var drugAllergy entities.DrugAllergy
	if err := r.db.Where("id = ?", id).First(&drugAllergy).Error; err != nil {
		return nil, err
	}
	return &drugAllergy, nil
}

func (r *GormEmrRepository) GetAllDrugAllergies() ([]*entities.DrugAllergy, error) {
	var drugAllergies []*entities.DrugAllergy
	if err := r.db.Find(&drugAllergies).Error; err != nil {
		return nil, err
	}
	return drugAllergies, nil
}

func (r *GormEmrRepository) DrugAllergyExists(allergyName string) (bool, error) {
	var count int64
	err := r.db.Model(&entities.DrugAllergy{}).Where("allergy_name = ?", allergyName).Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r *GormEmrRepository) DeleteDrugAllergy(id string) error {
	if err := r.db.Delete(&entities.DrugAllergy{}, "id = ?", id).Error; err != nil {
		return err
	}
	return nil
}

func (r *GormEmrRepository) GetResidentLabelsByResidentID(residentID string) ([]*entities.ResidentLabels, error) {
	var residentLabels []*entities.ResidentLabels
	if err := r.db.Preload("IntakeLabel").Where("resident_id = ?", residentID).Find(&residentLabels).Error; err != nil {
		return nil, err
	}
	return residentLabels, nil
}

func (r *GormEmrRepository) ResidentLabelExists(residentID, labelID string) (bool, error) {
	var count int64
	err := r.db.Model(&entities.ResidentLabels{}).Where("resident_id = ? AND label_id = ?", residentID, labelID).Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r *GormEmrRepository) DeleteResidentLabelsByResidentID(residentID string) error {
	if err := r.db.Where("resident_id = ?", residentID).Delete(&entities.ResidentLabels{}).Error; err != nil {
		return err
	}
	return nil
}

func (r *GormEmrRepository) CreateAllergyByResidentID(residentAllergy *entities.ResidentAllergies) (*entities.ResidentAllergies, error) {
	if err := r.db.Create(&residentAllergy).Error; err != nil {
		return nil, err
	}

	var result entities.ResidentAllergies
	if err := r.db.Preload("Resident").Preload("Allergy").
		Where("resident_id = ? AND allergy_id = ?", residentAllergy.ResidentID, residentAllergy.AllergyID).
		First(&result).Error; err != nil {
		return nil, err
	}
	return &result, nil
}

func (r *GormEmrRepository) GetResidentAllergiesByResidentID(residentID string) ([]*entities.ResidentAllergies, error) {
	var residentAllergies []*entities.ResidentAllergies
	if err := r.db.Preload("Allergy").Where("resident_id = ?", residentID).Find(&residentAllergies).Error; err != nil {
		return nil, err
	}
	return residentAllergies, nil
}

func (r *GormEmrRepository) GetAllResidentAllergies() ([]*models.ResidentAllergyListResponse, error) {
	var residents []*entities.Resident
	if err := r.db.
		Select("id", "first_name", "last_name").
		Preload("ResidentAllergies.Allergy").
		Find(&residents).Error; err != nil {
		return nil, err
	}

	result := make([]*models.ResidentAllergyListResponse, 0, len(residents))
	for _, resident := range residents {
		allergyItems := make([]models.ResidentAllergyItemResponse, 0, len(resident.ResidentAllergies))
		for _, residentAllergy := range resident.ResidentAllergies {
			allergyItems = append(allergyItems, models.ResidentAllergyItemResponse{
				AllergyID:   residentAllergy.AllergyID,
				AllergyName: residentAllergy.Allergy.AllergyName,
				NoteText:    residentAllergy.NoteText,
			})
		}

		result = append(result, &models.ResidentAllergyListResponse{
			ResidentID: resident.ID,
			FirstName:  resident.FirstName,
			LastName:   resident.LastName,
			Allergies:  allergyItems,
		})
	}

	return result, nil
}

func (r *GormEmrRepository) ResidentAllergyExists(residentID, allergyID string) (bool, error) {
	var count int64
	err := r.db.Model(&entities.ResidentAllergies{}).Where("resident_id = ? AND allergy_id = ?", residentID, allergyID).Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r *GormEmrRepository) DeleteResidentAllergiesByResidentID(residentID string) error {
	if err := r.db.Where("resident_id = ?", residentID).Delete(&entities.ResidentAllergies{}).Error; err != nil {
		return err
	}
	return nil
}

func (r *GormEmrRepository) CreateDrugAllergyByResidentID(residentDA *entities.ResidentDA) (*entities.ResidentDA, error) {
	if err := r.db.Create(&residentDA).Error; err != nil {
		return nil, err
	}

	var result entities.ResidentDA
	if err := r.db.Preload("Resident").Preload("DrugAllergy").
		Where("resident_id = ? AND drug_allergy_id = ?", residentDA.ResidentID, residentDA.DrugAllergyID).
		First(&result).Error; err != nil {
		return nil, err
	}
	return &result, nil
}

func (r *GormEmrRepository) GetResidentDrugAllergiesByResidentID(residentID string) ([]*entities.ResidentDA, error) {
	var residentDAs []*entities.ResidentDA
	if err := r.db.Preload("DrugAllergy").Where("resident_id = ?", residentID).Find(&residentDAs).Error; err != nil {
		return nil, err
	}
	return residentDAs, nil
}

func (r *GormEmrRepository) GetAllResidentDrugAllergies() ([]*models.ResidentDrugAllergyListResponse, error) {
	var residents []*entities.Resident
	if err := r.db.
		Select("id", "first_name", "last_name").
		Preload("ResidentDA.DrugAllergy").
		Find(&residents).Error; err != nil {
		return nil, err
	}

	result := make([]*models.ResidentDrugAllergyListResponse, 0, len(residents))
	for _, resident := range residents {
		drugAllergyItems := make([]models.ResidentDrugAllergyItemResponse, 0, len(resident.ResidentDA))
		for _, residentDA := range resident.ResidentDA {
			drugAllergyItems = append(drugAllergyItems, models.ResidentDrugAllergyItemResponse{
				DrugAllergyID: residentDA.DrugAllergyID,
				AllergyName:   residentDA.DrugAllergy.AllergyName,
				NoteText:      residentDA.NoteText,
			})
		}

		result = append(result, &models.ResidentDrugAllergyListResponse{
			ResidentID:    resident.ID,
			FirstName:     resident.FirstName,
			LastName:      resident.LastName,
			DrugAllergies: drugAllergyItems,
		})
	}

	return result, nil
}

func (r *GormEmrRepository) ResidentDrugAllergyExists(residentID, drugAllergyID string) (bool, error) {
	var count int64
	err := r.db.Model(&entities.ResidentDA{}).Where("resident_id = ? AND drug_allergy_id = ?", residentID, drugAllergyID).Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r *GormEmrRepository) DeleteResidentDrugAllergiesByResidentID(residentID string) error {
	if err := r.db.Where("resident_id = ?", residentID).Delete(&entities.ResidentDA{}).Error; err != nil {
		return err
	}
	return nil
}

func (r *GormEmrRepository) CreateVitalSign(vitalSign *entities.VitalSign) (*entities.VitalSign, error) {
	if err := r.db.Create(&vitalSign).Error; err != nil {
		return nil, err
	}
	return vitalSign, nil
}

func (r *GormEmrRepository) VitalSignSlotExists(residentID string, measurementDate time.Time, timeOfDay string) (bool, error) {
	var count int64
	err := r.db.Model(&entities.VitalSign{}).
		Where("resident_id = ?", residentID).
		Where("measurement_date = ?::date", measurementDate.Format("2006-01-02")).
		Where("LOWER(time_of_day) = LOWER(?)", timeOfDay).
		Count(&count).Error
	if err != nil {
		return false, err
	}

	return count > 0, nil
}

func (r *GormEmrRepository) GetVitalSignByID(id string) (*entities.VitalSign, error) {
	var vitalSign entities.VitalSign
	if err := r.db.Where("id = ?", id).First(&vitalSign).Error; err != nil {
		return nil, err
	}
	return &vitalSign, nil
}

func (r *GormEmrRepository) GetVitalSignsHistory(residentID string) ([]*entities.VitalSign, error) {
	var vitalSigns []*entities.VitalSign
	if err := r.db.Where("resident_id = ?", residentID).Find(&vitalSigns).Error; err != nil {
		return nil, err
	}
	return vitalSigns, nil
}

func (r *GormEmrRepository) GetVitalSignsByRoomIDToday(roomID string, isLatest bool) ([]*entities.VitalSign, error) {
	var vitalSigns []*entities.VitalSign

	query := r.db

	if isLatest {
		query = query.Table("vital_signs").
			Select("DISTINCT ON (vital_signs.resident_id) vital_signs.*")
	}

	query = query.Joins("JOIN residents ON vital_signs.resident_id = residents.id").
		Where("residents.room_id = ?", roomID).
		Where("vital_signs.created_at >= DATE_TRUNC('day', NOW() AT TIME ZONE 'Asia/Bangkok') AT TIME ZONE 'Asia/Bangkok'").
		Where("vital_signs.created_at < DATE_TRUNC('day', NOW() AT TIME ZONE 'Asia/Bangkok') AT TIME ZONE 'Asia/Bangkok' + INTERVAL '1 day'")

	if isLatest {
		query = query.Order("vital_signs.resident_id, vital_signs.created_at DESC")
	} else {
		query = query.Order("vital_signs.created_at DESC")
	}

	err := query.Find(&vitalSigns).Error
	if err != nil {
		return nil, err
	}

	return vitalSigns, nil
}

func (r *GormEmrRepository) GetVitalSignsByFloorToday(floor int16, isLatest bool) ([]*entities.VitalSign, error) {
	var vitalSigns []*entities.VitalSign

	query := r.db

	if isLatest {
		query = query.Table("vital_signs").
			Select("DISTINCT ON (vital_signs.resident_id) vital_signs.*")
	}

	query = query.Joins("JOIN residents ON vital_signs.resident_id = residents.id").
		Joins("JOIN rooms ON residents.room_id = rooms.id").
		Where("rooms.floor = ?", floor).
		Where("vital_signs.created_at >= DATE_TRUNC('day', NOW() AT TIME ZONE 'Asia/Bangkok') AT TIME ZONE 'Asia/Bangkok'").
		Where("vital_signs.created_at < DATE_TRUNC('day', NOW() AT TIME ZONE 'Asia/Bangkok') AT TIME ZONE 'Asia/Bangkok' + INTERVAL '1 day'")

	if isLatest {
		query = query.Order("vital_signs.resident_id, vital_signs.created_at DESC")
	} else {
		query = query.Order("vital_signs.created_at DESC")
	}

	err := query.Find(&vitalSigns).Error
	if err != nil {
		return nil, err
	}

	return vitalSigns, nil
}

func (r *GormEmrRepository) GetVitalSignsByResidentIDToday(residentID string, isLatest bool) ([]*entities.VitalSign, error) {
	bangkokNow := time.Now().In(time.FixedZone("ICT", 7*60*60))
	return r.GetVitalSignsByResidentIDOnDate(residentID, bangkokNow, isLatest)
}

func (r *GormEmrRepository) GetVitalSignsByResidentIDOnDate(residentID string, dayDate time.Time, isLatest bool) ([]*entities.VitalSign, error) {
	var vitalSigns []*entities.VitalSign

	query := r.db.
		Where("resident_id = ?", residentID).
		Where("measurement_date = ?::date", dayDate.Format("2006-01-02")).
		Order("created_at DESC")

	if isLatest {
		query = query.Limit(1)
	}

	err := query.Find(&vitalSigns).Error
	if err != nil {
		return nil, err
	}

	return vitalSigns, nil
}

func (r *GormEmrRepository) GetVitalSignsToday(isLatest bool) ([]*entities.VitalSign, error) {
	bangkokNow := time.Now().In(time.FixedZone("ICT", 7*60*60))
	return r.GetVitalSignsOnDate(bangkokNow, isLatest)
}

func (r *GormEmrRepository) GetVitalSignsOnDate(dayDate time.Time, isLatest bool) ([]*entities.VitalSign, error) {
	var vitalSigns []*entities.VitalSign

	query := r.db

	if isLatest {
		query = query.Table("vital_signs").
			Select("DISTINCT ON (vital_signs.resident_id) vital_signs.*")
	}

	query = query.Where("vital_signs.measurement_date = ?::date", dayDate.Format("2006-01-02"))

	if isLatest {
		query = query.Order("vital_signs.resident_id, vital_signs.created_at DESC")
	} else {
		query = query.Order("vital_signs.created_at DESC")
	}

	err := query.Find(&vitalSigns).Error
	if err != nil {
		return nil, err
	}

	return vitalSigns, nil
}

func (r *GormEmrRepository) GetVitalSignsCustom(params models.VitalSignQueryParams) ([]*entities.VitalSign, int64, error) {
	var vitalSigns []*entities.VitalSign

	buildQuery := func(withPagination bool) *gorm.DB {
		var query *gorm.DB

		if params.IsLatest {
			query = r.db.Table("vital_signs").Select("DISTINCT ON (vital_signs.resident_id) vital_signs.*")
		} else {
			query = r.db.Model(&entities.VitalSign{})
		}

		needResidentsJoin := false
		needRoomsJoin := false

		if params.ResidentID != nil && *params.ResidentID != "" {
			query = query.Where("vital_signs.resident_id = ?", *params.ResidentID)
		}

		if params.RoomID != nil && *params.RoomID != "" {
			needResidentsJoin = true
			query = query.Where("residents.room_id = ?", *params.RoomID)
		}

		if params.Floor != nil {
			needResidentsJoin = true
			needRoomsJoin = true
			query = query.Where("rooms.floor = ?", *params.Floor)
		}

		if len(params.LabelIDs) > 0 {
			subQuery := r.db.Table("resident_labels").
				Select("resident_id").
				Where("label_id IN ?", params.LabelIDs).
				Group("resident_id").
				Having("COUNT(DISTINCT label_id) = ?", len(params.LabelIDs))
			query = query.Where("vital_signs.resident_id IN (?)", subQuery)
		}

		if params.TimeOfDay != nil && *params.TimeOfDay != "" {
			query = query.Where("LOWER(vital_signs.time_of_day) = LOWER(?)", *params.TimeOfDay)
		}

		if needResidentsJoin {
			query = query.Joins("JOIN residents ON vital_signs.resident_id = residents.id")
		}
		if needRoomsJoin {
			query = query.Joins("JOIN rooms ON residents.room_id = rooms.id")
		}

		if params.Date != nil && *params.Date != "" {
			query = query.Where("vital_signs.measurement_date = ?::date", *params.Date)
		} else if params.StartDate != nil {
			query = query.Where("vital_signs.created_at >= ?", *params.StartDate)
		} else {
			query = query.Where("vital_signs.created_at >= DATE_TRUNC('day', NOW() AT TIME ZONE 'Asia/Bangkok') AT TIME ZONE 'Asia/Bangkok'")
		}
		if params.Date != nil && *params.Date != "" {
			// date filter already handles full-day selection
		} else if params.EndDate != nil {
			endDateInclusive := params.EndDate.AddDate(0, 0, 1)
			query = query.Where("vital_signs.created_at < ?", endDateInclusive)
		} else {
			query = query.Where("vital_signs.created_at < DATE_TRUNC('day', NOW() AT TIME ZONE 'Asia/Bangkok') AT TIME ZONE 'Asia/Bangkok' + INTERVAL '1 day'")
		}

		if params.IsLatest {
			query = query.Order("vital_signs.resident_id, vital_signs.created_at DESC")
		} else {
			query = query.Order("vital_signs.created_at DESC")
		}

		if withPagination {
			if params.Limit > 0 {
				query = query.Limit(params.Limit)
			}
			if params.Offset > 0 {
				query = query.Offset(params.Offset)
			}
		}

		return query
	}

	countSubQuery := buildQuery(false).Select("vital_signs.id")
	var total int64
	if err := r.db.Table("(?) AS filtered_vital_signs", countSubQuery).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := buildQuery(true).Find(&vitalSigns).Error
	if err != nil {
		return nil, 0, err
	}

	return vitalSigns, total, nil
}

func (r *GormEmrRepository) UpdateVitalSignByID(vitalSign *entities.VitalSign) (*entities.VitalSign, error) {
	if err := r.db.Save(&vitalSign).Error; err != nil {
		return nil, err
	}
	return vitalSign, nil
}

func (r *GormEmrRepository) CreateLaboratoryValue(laboratoryValue *entities.LaboratoryValue) (*entities.LaboratoryValue, error) {
	if err := r.db.Create(&laboratoryValue).Error; err != nil {
		return nil, err
	}
	return laboratoryValue, nil
}

func (r *GormEmrRepository) GetLaboratoryValueByID(id string) (*entities.LaboratoryValue, error) {
	var lab entities.LaboratoryValue
	if err := r.db.Where("id = ?", id).First(&lab).Error; err != nil {
		return nil, err
	}
	return &lab, nil
}

func (r *GormEmrRepository) GetLaboratoryValuesHistory(residentID string) ([]*entities.LaboratoryValue, error) {
	var labs []*entities.LaboratoryValue
	if err := r.db.Where("resident_id = ?", residentID).Order("created_at DESC").Find(&labs).Error; err != nil {
		return nil, err
	}
	return labs, nil
}

func (r *GormEmrRepository) GetLaboratoryValuesByRoomIDToday(roomID string, isLatest bool) ([]*entities.LaboratoryValue, error) {
	var labs []*entities.LaboratoryValue

	query := r.db

	if isLatest {
		query = query.Table("laboratory_values").
			Select("DISTINCT ON (laboratory_values.resident_id) laboratory_values.*")
	}

	query = query.Joins("JOIN residents ON laboratory_values.resident_id = residents.id").
		Where("residents.room_id = ?", roomID).
		Where("laboratory_values.created_at >= DATE_TRUNC('day', NOW() AT TIME ZONE 'Asia/Bangkok') AT TIME ZONE 'Asia/Bangkok'").
		Where("laboratory_values.created_at < DATE_TRUNC('day', NOW() AT TIME ZONE 'Asia/Bangkok') AT TIME ZONE 'Asia/Bangkok' + INTERVAL '1 day'")

	if isLatest {
		query = query.Order("laboratory_values.resident_id, laboratory_values.created_at DESC")
	} else {
		query = query.Order("laboratory_values.created_at DESC")
	}

	if err := query.Find(&labs).Error; err != nil {
		return nil, err
	}
	return labs, nil
}

func (r *GormEmrRepository) GetLaboratoryValuesByFloorToday(floor int16, isLatest bool) ([]*entities.LaboratoryValue, error) {
	var labs []*entities.LaboratoryValue

	query := r.db

	if isLatest {
		query = query.Table("laboratory_values").
			Select("DISTINCT ON (laboratory_values.resident_id) laboratory_values.*")
	}

	query = query.Joins("JOIN residents ON laboratory_values.resident_id = residents.id").
		Joins("JOIN rooms ON residents.room_id = rooms.id").
		Where("rooms.floor = ?", floor).
		Where("laboratory_values.created_at >= DATE_TRUNC('day', NOW() AT TIME ZONE 'Asia/Bangkok') AT TIME ZONE 'Asia/Bangkok'").
		Where("laboratory_values.created_at < DATE_TRUNC('day', NOW() AT TIME ZONE 'Asia/Bangkok') AT TIME ZONE 'Asia/Bangkok' + INTERVAL '1 day'")

	if isLatest {
		query = query.Order("laboratory_values.resident_id, laboratory_values.created_at DESC")
	} else {
		query = query.Order("laboratory_values.created_at DESC")
	}

	if err := query.Find(&labs).Error; err != nil {
		return nil, err
	}
	return labs, nil
}

func (r *GormEmrRepository) GetLaboratoryValuesByResidentIDToday(residentID string, isLatest bool) ([]*entities.LaboratoryValue, error) {
	var labs []*entities.LaboratoryValue

	query := r.db.
		Where("resident_id = ?", residentID).
		Where("created_at >= DATE_TRUNC('day', NOW() AT TIME ZONE 'Asia/Bangkok') AT TIME ZONE 'Asia/Bangkok'").
		Where("created_at < DATE_TRUNC('day', NOW() AT TIME ZONE 'Asia/Bangkok') AT TIME ZONE 'Asia/Bangkok' + INTERVAL '1 day'").
		Order("created_at DESC")

	if isLatest {
		query = query.Limit(1)
	}

	if err := query.Find(&labs).Error; err != nil {
		return nil, err
	}
	return labs, nil
}

func (r *GormEmrRepository) GetLaboratoryValuesToday(isLatest bool) ([]*entities.LaboratoryValue, error) {
	var labs []*entities.LaboratoryValue

	query := r.db

	if isLatest {
		query = query.Table("laboratory_values").
			Select("DISTINCT ON (laboratory_values.resident_id) laboratory_values.*")
	}

	query = query.
		Where("laboratory_values.created_at >= DATE_TRUNC('day', NOW() AT TIME ZONE 'Asia/Bangkok') AT TIME ZONE 'Asia/Bangkok'").
		Where("laboratory_values.created_at < DATE_TRUNC('day', NOW() AT TIME ZONE 'Asia/Bangkok') AT TIME ZONE 'Asia/Bangkok' + INTERVAL '1 day'")

	if isLatest {
		query = query.Order("laboratory_values.resident_id, laboratory_values.created_at DESC")
	} else {
		query = query.Order("laboratory_values.created_at DESC")
	}

	if err := query.Find(&labs).Error; err != nil {
		return nil, err
	}
	return labs, nil
}

func (r *GormEmrRepository) GetLaboratoryValuesCustom(params models.LaboratoryValueQueryParams) ([]*entities.LaboratoryValue, int64, error) {
	var labs []*entities.LaboratoryValue

	buildQuery := func(withPagination bool) *gorm.DB {
		var query *gorm.DB

		if params.IsLatest {
			query = r.db.Table("laboratory_values").Select("DISTINCT ON (laboratory_values.resident_id) laboratory_values.*")
		} else {
			query = r.db.Model(&entities.LaboratoryValue{})
		}

		needResidentsJoin := false
		needRoomsJoin := false

		if params.ResidentID != nil && *params.ResidentID != "" {
			query = query.Where("laboratory_values.resident_id = ?", *params.ResidentID)
		}

		if params.RoomID != nil && *params.RoomID != "" {
			needResidentsJoin = true
			query = query.Where("residents.room_id = ?", *params.RoomID)
		}

		if params.Floor != nil {
			needResidentsJoin = true
			needRoomsJoin = true
			query = query.Where("rooms.floor = ?", *params.Floor)
		}

		if len(params.LabelIDs) > 0 {
			subQuery := r.db.Table("resident_labels").
				Select("resident_id").
				Where("label_id IN ?", params.LabelIDs).
				Group("resident_id").
				Having("COUNT(DISTINCT label_id) = ?", len(params.LabelIDs))
			query = query.Where("laboratory_values.resident_id IN (?)", subQuery)
		}

		if needResidentsJoin {
			query = query.Joins("JOIN residents ON laboratory_values.resident_id = residents.id")
		}
		if needRoomsJoin {
			query = query.Joins("JOIN rooms ON residents.room_id = rooms.id")
		}

		if params.StartDate != nil {
			query = query.Where("laboratory_values.created_at >= ?", *params.StartDate)
		} else {
			query = query.Where("laboratory_values.created_at >= DATE_TRUNC('day', NOW() AT TIME ZONE 'Asia/Bangkok') AT TIME ZONE 'Asia/Bangkok'")
		}
		if params.EndDate != nil {
			endDateInclusive := params.EndDate.AddDate(0, 0, 1)
			query = query.Where("laboratory_values.created_at < ?", endDateInclusive)
		} else {
			query = query.Where("laboratory_values.created_at < DATE_TRUNC('day', NOW() AT TIME ZONE 'Asia/Bangkok') AT TIME ZONE 'Asia/Bangkok' + INTERVAL '1 day'")
		}

		if params.IsLatest {
			query = query.Order("laboratory_values.resident_id, laboratory_values.created_at DESC")
		} else {
			query = query.Order("laboratory_values.created_at DESC")
		}

		if withPagination {
			if params.Limit > 0 {
				query = query.Limit(params.Limit)
			}
			if params.Offset > 0 {
				query = query.Offset(params.Offset)
			}
		}

		return query
	}

	countSubQuery := buildQuery(false).Select("laboratory_values.id")
	var total int64
	if err := r.db.Table("(?) AS filtered_laboratory_values", countSubQuery).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := buildQuery(true).Find(&labs).Error; err != nil {
		return nil, 0, err
	}
	return labs, total, nil
}

func (r *GormEmrRepository) UpdateLaboratoryValueByID(laboratoryValue *entities.LaboratoryValue) (*entities.LaboratoryValue, error) {
	if err := r.db.Save(&laboratoryValue).Error; err != nil {
		return nil, err
	}
	return laboratoryValue, nil
}

func (r *GormEmrRepository) LaboratoryValueSlotExists(residentID string, measurementDate time.Time, timeOfDay string) (bool, error) {
	var count int64
	r.db.Model(&entities.LaboratoryValue{}).
		Where("resident_id = ?", residentID).
		Where("measurement_date = ?::date", measurementDate.Format("2006-01-02")).
		Where("LOWER(time_of_day) = LOWER(?)", timeOfDay).
		Count(&count)
	return count > 0, nil
}

func (r *GormEmrRepository) GetLaboratoryValuesByResidentIDOnDate(residentID string, dayDate time.Time, isLatest bool) ([]*entities.LaboratoryValue, error) {
	query := r.db.
		Where("resident_id = ?", residentID).
		Where("measurement_date = ?::date", dayDate.Format("2006-01-02")).
		Order("created_at DESC")
	if isLatest {
		query = query.Limit(1)
	}
	var labs []*entities.LaboratoryValue
	err := query.Find(&labs).Error
	return labs, err
}

func (r *GormEmrRepository) GetLaboratoryValuesOnDate(dayDate time.Time, isLatest bool) ([]*entities.LaboratoryValue, error) {
	query := r.db
	if isLatest {
		query = query.Table("laboratory_values").
			Select("DISTINCT ON (laboratory_values.resident_id) laboratory_values.*")
	}
	query = query.
		Where("measurement_date = ?::date", dayDate.Format("2006-01-02")).
		Order("resident_id, created_at DESC")
	var labs []*entities.LaboratoryValue
	err := query.Find(&labs).Error
	return labs, err
}

func (r *GormEmrRepository) GetUrineOutputSumGroupByResident(params models.LaboratoryValueQueryParams, urineType string) (*models.UrineOutputSumResponse, error) {
	var result models.UrineOutputSumResponse

	query := r.db.Table("laboratory_values").
		Select("laboratory_values.resident_id, COALESCE(SUM(laboratory_values.urine_output), 0) as total_amount").
		Where("laboratory_values.urine_type = ?", urineType)

	query = applyLaboratoryValueQueryFilters(query, params)

	query = query.Group("laboratory_values.resident_id")

	if err := query.Scan(&result).Error; err != nil {
		return nil, err
	}
	return &result, nil
}

func applyLaboratoryValueQueryFilters(query *gorm.DB, params models.LaboratoryValueQueryParams) *gorm.DB {
	needResidentsJoin := false
	needRoomsJoin := false

	if params.ResidentID != nil && *params.ResidentID != "" {
		query = query.Where("laboratory_values.resident_id = ?", *params.ResidentID)
	}
	if params.RoomID != nil && *params.RoomID != "" {
		needResidentsJoin = true
		query = query.Where("residents.room_id = ?", *params.RoomID)
	}
	if params.Floor != nil {
		needResidentsJoin = true
		needRoomsJoin = true
		query = query.Where("rooms.floor = ?", *params.Floor)
	}
	if len(params.LabelIDs) > 0 {
		subQuery := query.Session(&gorm.Session{NewDB: true}).Table("resident_labels").
			Select("resident_id").
			Where("label_id IN ?", params.LabelIDs).
			Group("resident_id").
			Having("COUNT(DISTINCT label_id) = ?", len(params.LabelIDs))
		query = query.Where("laboratory_values.resident_id IN (?)", subQuery)
	}

	if needResidentsJoin {
		query = query.Joins("JOIN residents ON laboratory_values.resident_id = residents.id")
	}
	if needRoomsJoin {
		query = query.Joins("JOIN rooms ON residents.room_id = rooms.id")
	}

	if params.StartDate != nil {
		query = query.Where("laboratory_values.created_at >= ?", *params.StartDate)
	}
	if params.EndDate != nil {
		endDateInclusive := params.EndDate.AddDate(0, 0, 1)
		query = query.Where("laboratory_values.created_at < ?", endDateInclusive)
	}

	return query
}

func (r *GormEmrRepository) CreateNurseNote(note *entities.NurseNote) (*entities.NurseNote, error) {
	if err := r.db.Create(&note).Error; err != nil {
		return nil, err
	}
	return r.GetNurseNoteByID(note.ID)
}

func (r *GormEmrRepository) GetNurseNoteByID(id string) (*entities.NurseNote, error) {
	var note entities.NurseNote
	if err := r.db.Preload("Resident").Where("id = ?", id).First(&note).Error; err != nil {
		return nil, err
	}
	return &note, nil
}

func (r *GormEmrRepository) GetNurseNotesOverviewOnDate(dayDate time.Time) ([]*entities.NurseNote, error) {
	var notes []*entities.NurseNote
	if err := r.db.
		Where("DATE(created_at AT TIME ZONE 'Asia/Bangkok') = ?", dayDate.Format("2006-01-02")).
		Order("created_at DESC").
		Find(&notes).Error; err != nil {
		return nil, err
	}
	return notes, nil
}

func (r *GormEmrRepository) GetNurseNotesByResidentIDOnDate(residentID string, dayDate time.Time) ([]*entities.NurseNote, error) {
	var notes []*entities.NurseNote
	if err := r.db.
		Where("resident_id = ?", residentID).
		Where("DATE(created_at AT TIME ZONE 'Asia/Bangkok') = ?", dayDate.Format("2006-01-02")).
		Order("created_at DESC").
		Find(&notes).Error; err != nil {
		return nil, err
	}
	return notes, nil
}

func (r *GormEmrRepository) UpdateNurseNoteByID(note *entities.NurseNote) (*entities.NurseNote, error) {
	if err := r.db.Save(&note).Error; err != nil {
		return nil, err
	}
	return r.GetNurseNoteByID(note.ID)
}

func (r *GormEmrRepository) DeleteNurseNoteByID(id string) error {
	return r.db.Where("id = ?", id).Delete(&entities.NurseNote{}).Error
}

func (r *GormEmrRepository) CreateWoundCareNote(note *entities.WoundCareNote) (*entities.WoundCareNote, error) {
	if err := r.db.Create(&note).Error; err != nil {
		return nil, err
	}
	return r.GetWoundCareNoteByID(note.ID)
}

func (r *GormEmrRepository) GetWoundCareNoteByID(id string) (*entities.WoundCareNote, error) {
	var note entities.WoundCareNote
	if err := r.db.Preload("Resident").Where("id = ?", id).First(&note).Error; err != nil {
		return nil, err
	}
	return &note, nil
}

func (r *GormEmrRepository) GetWoundCareNotesOverviewOnDate(dayDate time.Time) ([]*entities.WoundCareNote, error) {
	var notes []*entities.WoundCareNote
	if err := r.db.
		Where("DATE(created_at AT TIME ZONE 'Asia/Bangkok') = ?", dayDate.Format("2006-01-02")).
		Order("created_at DESC").
		Find(&notes).Error; err != nil {
		return nil, err
	}
	return notes, nil
}

func (r *GormEmrRepository) GetWoundCareNotesByResidentIDOnDate(residentID string, dayDate time.Time) ([]*entities.WoundCareNote, error) {
	var notes []*entities.WoundCareNote
	if err := r.db.
		Where("resident_id = ?", residentID).
		Where("DATE(created_at AT TIME ZONE 'Asia/Bangkok') = ?", dayDate.Format("2006-01-02")).
		Order("created_at DESC").
		Find(&notes).Error; err != nil {
		return nil, err
	}
	return notes, nil
}

func (r *GormEmrRepository) UpdateWoundCareNoteByID(note *entities.WoundCareNote) (*entities.WoundCareNote, error) {
	if err := r.db.Save(&note).Error; err != nil {
		return nil, err
	}
	return r.GetWoundCareNoteByID(note.ID)
}

func (r *GormEmrRepository) DeleteWoundCareNoteByID(id string) error {
	return r.db.Where("id = ?", id).Delete(&entities.WoundCareNote{}).Error
}

func (r *GormEmrRepository) CreateRelativeNote(note *entities.RelativeNote) (*entities.RelativeNote, error) {
	if err := r.db.Create(&note).Error; err != nil {
		return nil, err
	}
	return r.GetRelativeNoteByID(note.ID)
}

func (r *GormEmrRepository) GetRelativeNoteByID(id string) (*entities.RelativeNote, error) {
	var note entities.RelativeNote
	if err := r.db.Preload("Resident").Where("id = ?", id).First(&note).Error; err != nil {
		return nil, err
	}
	return &note, nil
}

func (r *GormEmrRepository) GetRelativeNotesOverviewOnDate(dayDate time.Time) ([]*entities.RelativeNote, error) {
	var notes []*entities.RelativeNote
	if err := r.db.
		Where("DATE(created_at AT TIME ZONE 'Asia/Bangkok') = ?", dayDate.Format("2006-01-02")).
		Order("created_at DESC").
		Find(&notes).Error; err != nil {
		return nil, err
	}
	return notes, nil
}

func (r *GormEmrRepository) GetRelativeNotesByResidentIDOnDate(residentID string, dayDate time.Time) ([]*entities.RelativeNote, error) {
	var notes []*entities.RelativeNote
	if err := r.db.
		Where("resident_id = ?", residentID).
		Where("DATE(created_at AT TIME ZONE 'Asia/Bangkok') = ?", dayDate.Format("2006-01-02")).
		Order("created_at DESC").
		Find(&notes).Error; err != nil {
		return nil, err
	}
	return notes, nil
}

func (r *GormEmrRepository) UpdateRelativeNoteByID(note *entities.RelativeNote) (*entities.RelativeNote, error) {
	if err := r.db.Save(&note).Error; err != nil {
		return nil, err
	}
	return r.GetRelativeNoteByID(note.ID)
}

func (r *GormEmrRepository) DeleteRelativeNoteByID(id string) error {
	return r.db.Where("id = ?", id).Delete(&entities.RelativeNote{}).Error
}

func (r *GormEmrRepository) CreateDoctorOrder(order *entities.DoctorOrder) (*entities.DoctorOrder, error) {
	if err := r.db.Create(&order).Error; err != nil {
		return nil, err
	}
	return r.GetDoctorOrderByID(order.ID)
}

func (r *GormEmrRepository) GetDoctorOrderByID(id string) (*entities.DoctorOrder, error) {
	var order entities.DoctorOrder
	if err := r.db.Preload("Resident").Where("id = ?", id).First(&order).Error; err != nil {
		return nil, err
	}
	return &order, nil
}

func (r *GormEmrRepository) GetDoctorOrdersOverviewOnDate(dayDate time.Time) ([]*entities.DoctorOrder, error) {
	var orders []*entities.DoctorOrder
	if err := r.db.
		Where("order_date = ?", dayDate.Format("2006-01-02")).
		Order("created_at DESC").
		Find(&orders).Error; err != nil {
		return nil, err
	}
	return orders, nil
}

func (r *GormEmrRepository) GetDoctorOrdersByResidentIDOnDate(residentID string, dayDate time.Time) ([]*entities.DoctorOrder, error) {
	var orders []*entities.DoctorOrder
	if err := r.db.
		Where("resident_id = ?", residentID).
		Where("order_date = ?", dayDate.Format("2006-01-02")).
		Order("created_at DESC").
		Find(&orders).Error; err != nil {
		return nil, err
	}
	return orders, nil
}

func (r *GormEmrRepository) UpdateDoctorOrderByID(order *entities.DoctorOrder) (*entities.DoctorOrder, error) {
	if err := r.db.Save(&order).Error; err != nil {
		return nil, err
	}
	return r.GetDoctorOrderByID(order.ID)
}

func (r *GormEmrRepository) DeleteDoctorOrderByID(id string) error {
	return r.db.Where("id = ?", id).Delete(&entities.DoctorOrder{}).Error
}
