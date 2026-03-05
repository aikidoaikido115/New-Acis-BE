package repositories

import (
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

	// Room operations
	RoomExists(id string) (bool, error)
	GetRoomByID(id string) (*entities.Room, error)
	GetAllRooms() ([]*entities.Room, error)
	CreateRoom(room *entities.Room) (*entities.Room, error)
	UpdateRoom(room *entities.Room) (*entities.Room, error)
	RoomNumberExists(roomNumber string) (bool, error)

	// IntakeLabel operations
	CreateIntakeLabel(label *entities.IntakeLabels) (*entities.IntakeLabels, error)
	GetIntakeLabelByID(id string) (*entities.IntakeLabels, error)
	GetIntakeLabelByName(labelName string) (*entities.IntakeLabels, error)
	GetAllIntakeLabels() ([]*entities.IntakeLabels, error)
	LabelExists(labelName string) (bool, error)
	DeleteIntakeLabel(id string) error

	// ResidentLabel operations (many-to-many)
	CreateIntakeLabelByResidentID(residentLabel *entities.ResidentLabels) (*entities.ResidentLabels, error)
	GetResidentLabelsByResidentID(residentID string) ([]*entities.ResidentLabels, error)
	ResidentLabelExists(residentID, labelID string) (bool, error)
	DeleteResidentLabelsByResidentID(residentID string) error

	// VitalSign operations
	CreateVitalSign(vitalSign *entities.VitalSign) (*entities.VitalSign, error)

	GetVitalSignByID(id string) (*entities.VitalSign, error)
	GetVitalSignsByRoomIDToday(roomID string, isLatest bool) ([]*entities.VitalSign, error)
	GetVitalSignsByFloorToday(floor int16, isLatest bool) ([]*entities.VitalSign, error)
	GetVitalSignsByResidentIDToday(residentID string, isLatest bool) ([]*entities.VitalSign, error)
	GetVitalSignsHistory(residentID string) ([]*entities.VitalSign, error)
	GetVitalSignsToday(isLatest bool) ([]*entities.VitalSign, error)
	GetVitalSignsCustom(params models.VitalSignQueryParams) ([]*entities.VitalSign, error)

	UpdateVitalSignByID(vitalSign *entities.VitalSign) (*entities.VitalSign, error)

	//LaboratoryValue operations
	// CreateLaboratoryValue(laboratoryValue *entities.LaboratoryValue) (*entities.LaboratoryValue, error)

	// todo เพราะมันเจาะจงว่า ค่าไหนของ vital sign อีกทีนึง
	// GetLatestVitalSignsGreaterThanCustom(params models.VitalSignQueryParams, greaterThan float64) ([]*entities.VitalSign, error)
	// GetLatestVitalSignsLessThanCustom(params models.VitalSignQueryParams, lessThan float64) ([]*entities.VitalSign, error)

	//todo Allergy ย้ายไป Meal repository
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
	if err := r.db.Preload("Room").Preload("ResidentLabels.IntakeLabel").Where("id = ?", id).First(&resident).Error; err != nil {
		return nil, err
	}
	return &resident, nil
}

func (r *GormEmrRepository) GetResidentByRoomID(roomID string) ([]*entities.Resident, error) {
	var residents []*entities.Resident
	if err := r.db.Preload("Room").Preload("ResidentLabels.IntakeLabel").Where("room_id = ?", roomID).Find(&residents).Error; err != nil {
		return nil, err
	}
	return residents, nil
}

func (r *GormEmrRepository) GetAllResidents() ([]*entities.Resident, error) {
	var residents []*entities.Resident
	if err := r.db.Preload("Room").Preload("ResidentLabels.IntakeLabel").Find(&residents).Error; err != nil {
		return nil, err
	}
	return residents, nil
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

func (r *GormEmrRepository) CreateVitalSign(vitalSign *entities.VitalSign) (*entities.VitalSign, error) {
	if err := r.db.Create(&vitalSign).Error; err != nil {
		return nil, err
	}
	return vitalSign, nil
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
	var vitalSigns []*entities.VitalSign

	query := r.db.
		Where("resident_id = ?", residentID).
		Where("created_at >= DATE_TRUNC('day', NOW() AT TIME ZONE 'Asia/Bangkok') AT TIME ZONE 'Asia/Bangkok'").
		Where("created_at < DATE_TRUNC('day', NOW() AT TIME ZONE 'Asia/Bangkok') AT TIME ZONE 'Asia/Bangkok' + INTERVAL '1 day'").
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
	var vitalSigns []*entities.VitalSign

	query := r.db

	if isLatest {
		query = query.Table("vital_signs").
			Select("DISTINCT ON (vital_signs.resident_id) vital_signs.*")
	}

	query = query.Where("vital_signs.created_at >= DATE_TRUNC('day', NOW() AT TIME ZONE 'Asia/Bangkok') AT TIME ZONE 'Asia/Bangkok'").
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

func (r *GormEmrRepository) GetVitalSignsCustom(params models.VitalSignQueryParams) ([]*entities.VitalSign, error) {
	var vitalSigns []*entities.VitalSign

	var query *gorm.DB

	if params.IsLatest {
		query = r.db.Table("vital_signs").Select("DISTINCT ON (vital_signs.resident_id) vital_signs.*")
	} else {
		query = r.db.Model(&entities.VitalSign{})
	}

	needResidentsJoin := false
	needRoomsJoin := false
	needLabelsJoin := false

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
		needResidentsJoin = true
		needLabelsJoin = true
		query = query.Where("resident_labels.label_id IN ?", params.LabelIDs)
	}

	if needResidentsJoin {
		query = query.Joins("JOIN residents ON vital_signs.resident_id = residents.id")
	}
	if needRoomsJoin {
		query = query.Joins("JOIN rooms ON residents.room_id = rooms.id")
	}
	if needLabelsJoin {
		query = query.Joins("JOIN resident_labels ON residents.id = resident_labels.resident_id")
	}

	if params.StartDate != nil {
		query = query.Where("vital_signs.created_at >= ?", *params.StartDate)
	} else {
		query = query.Where("vital_signs.created_at >= DATE_TRUNC('day', NOW() AT TIME ZONE 'Asia/Bangkok') AT TIME ZONE 'Asia/Bangkok'")
	}
	if params.EndDate != nil {
		endDateInclusive := params.EndDate.AddDate(0, 0, 1)
		query = query.Where("vital_signs.created_at < ?", endDateInclusive)
	} else {
		query = query.Where("vital_signs.created_at < DATE_TRUNC('day', NOW() AT TIME ZONE 'Asia/Bangkok') AT TIME ZONE 'Asia/Bangkok' + INTERVAL '1 day'")
	}

	// Order: สำหรับ Latest จะเรียง resident_id ก่อน, แล้ว created_at DESC (สำคัญสำหรับ DISTINCT ON)
	if params.IsLatest {
		query = query.Order("vital_signs.resident_id, vital_signs.created_at DESC")
	} else {
		query = query.Order("vital_signs.created_at DESC")
	}

	// Pagination
	if params.Limit > 0 {
		query = query.Limit(params.Limit)
	}
	if params.Offset > 0 {
		query = query.Offset(params.Offset)
	}

	err := query.Find(&vitalSigns).Error
	if err != nil {
		return nil, err
	}

	return vitalSigns, nil
}

func (r *GormEmrRepository) UpdateVitalSignByID(vitalSign *entities.VitalSign) (*entities.VitalSign, error) {
	if err := r.db.Save(&vitalSign).Error; err != nil {
		return nil, err
	}
	return vitalSign, nil
}
