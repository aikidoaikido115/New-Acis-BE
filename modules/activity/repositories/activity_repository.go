package repositories

import (
	"time"

	activityModels "github.com/aikidoaikido115/New-Acis-BE/modules/activity/models"
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

	// Activity operations
	CreateActivity(activity *entities.Activity) (*entities.Activity, error)
	GetActivityByID(id string) (*entities.Activity, error)
	GetActivityByName(activityName string) (*entities.Activity, error)
	GetAllActivities() ([]*entities.Activity, error)
	UpdateActivity(activity *entities.Activity) (*entities.Activity, error)
	DeleteActivity(id string) error

	// Activity Schedule operations
	CreateActivitySchedule(activitySchedule *entities.ActivitySchedule) (*entities.ActivitySchedule, error)
	CreateActivityScheduleWithDefaultParticipations(activitySchedule *entities.ActivitySchedule) (*entities.ActivitySchedule, error)
	GetActivityScheduleByID(id string) (*entities.ActivitySchedule, error)
	GetAllActivitySchedules() ([]*entities.ActivitySchedule, error)
	GetActivitySchedulesWithActivitySyncByDate(date *time.Time) ([]*entities.ActivitySchedule, error)
	UpdateActivitySchedule(activitySchedule *entities.ActivitySchedule) (*entities.ActivitySchedule, error)
	DeleteActivitySchedule(id string) error

	// Participation operations
	CreateParticipation(participation *entities.Participation) (*entities.Participation, error)
	GetParticipationByResidentIDAndASID(residentID, asID string) (*entities.Participation, error)
	GetParticipationsByASIDAndResidentIDs(asID string, residentIDs []string) ([]*entities.Participation, error)
	GetAllParticipations() ([]*entities.Participation, error)
	GetResidentsByScheduleIDCustom(asID string, params activityModels.ResidentsByScheduleQueryParams) ([]*entities.Participation, error)
	UpdateParticipation(participation *entities.Participation) (*entities.Participation, error)
	UpdateParticipationIsParticipatingByResidentIDs(asID string, residentIDs []string, isParticipating bool) error
	DeleteParticipation(residentID, asID string) error
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

func (r *GormActivityRepository) GetActivityByName(activityName string) (*entities.Activity, error) {
	var activity entities.Activity
	if err := r.db.Preload("Staff").Where("LOWER(activity_name) = LOWER(?)", activityName).First(&activity).Error; err != nil {
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

func (r *GormActivityRepository) CreateActivitySchedule(activitySchedule *entities.ActivitySchedule) (*entities.ActivitySchedule, error) {
	if err := r.db.Create(&activitySchedule).Error; err != nil {
		return nil, err
	}

	return r.GetActivityScheduleByID(activitySchedule.ID)
}

func (r *GormActivityRepository) CreateActivityScheduleWithDefaultParticipations(activitySchedule *entities.ActivitySchedule) (*entities.ActivitySchedule, error) {
	err := r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&activitySchedule).Error; err != nil {
			return err
		}

		var residentIDs []string
		if err := tx.Model(&entities.Resident{}).Pluck("id", &residentIDs).Error; err != nil {
			return err
		}

		if len(residentIDs) == 0 {
			return nil
		}

		defaultParticipations := make([]entities.Participation, 0, len(residentIDs))
		for _, residentID := range residentIDs {
			defaultParticipations = append(defaultParticipations, entities.Participation{
				ResidentID:      residentID,
				ASID:            activitySchedule.ID,
				IsParticipating: false,
				ImgURLs:         []entities.ImageURL{},
			})
		}

		if err := tx.Create(&defaultParticipations).Error; err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return r.GetActivityScheduleByID(activitySchedule.ID)
}

func (r *GormActivityRepository) GetActivityScheduleByID(id string) (*entities.ActivitySchedule, error) {
	var activitySchedule entities.ActivitySchedule
	if err := r.db.Preload("Activity").Where("id = ?", id).First(&activitySchedule).Error; err != nil {
		return nil, err
	}

	return &activitySchedule, nil
}

func (r *GormActivityRepository) GetAllActivitySchedules() ([]*entities.ActivitySchedule, error) {
	var activitySchedules []*entities.ActivitySchedule
	if err := r.db.Preload("Activity").Find(&activitySchedules).Error; err != nil {
		return nil, err
	}

	return activitySchedules, nil
}

func (r *GormActivityRepository) GetActivitySchedulesWithActivitySyncByDate(date *time.Time) ([]*entities.ActivitySchedule, error) {
	var activitySchedules []*entities.ActivitySchedule

	query := r.db.Model(&entities.ActivitySchedule{}).Preload("Activity")
	if date != nil {
		loc := time.FixedZone("ICT", 7*60*60)
		dayStart := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, loc)
		dayEnd := dayStart.AddDate(0, 0, 1)

		query = query.Where("activity_schedules.date >= ? AND activity_schedules.date < ?", dayStart, dayEnd)
	}

	if err := query.Order("activity_schedules.date ASC").Order("activity_schedules.start_time ASC").Find(&activitySchedules).Error; err != nil {
		return nil, err
	}

	return activitySchedules, nil
}

func (r *GormActivityRepository) UpdateActivitySchedule(activitySchedule *entities.ActivitySchedule) (*entities.ActivitySchedule, error) {
	if err := r.db.Save(&activitySchedule).Error; err != nil {
		return nil, err
	}

	return r.GetActivityScheduleByID(activitySchedule.ID)
}

func (r *GormActivityRepository) DeleteActivitySchedule(id string) error {
	if err := r.db.Delete(&entities.ActivitySchedule{}, "id = ?", id).Error; err != nil {
		return err
	}

	return nil
}

func (r *GormActivityRepository) CreateParticipation(participation *entities.Participation) (*entities.Participation, error) {
	if err := r.db.Create(&participation).Error; err != nil {
		return nil, err
	}

	return r.GetParticipationByResidentIDAndASID(participation.ResidentID, participation.ASID)
}

func (r *GormActivityRepository) GetParticipationByResidentIDAndASID(residentID, asID string) (*entities.Participation, error) {
	var participation entities.Participation
	if err := r.db.Preload("Resident").Preload("ActivitySchedule").
		Where("resident_id = ? AND as_id = ?", residentID, asID).
		First(&participation).Error; err != nil {
		return nil, err
	}

	return &participation, nil
}

func (r *GormActivityRepository) GetParticipationsByASIDAndResidentIDs(asID string, residentIDs []string) ([]*entities.Participation, error) {
	var participations []*entities.Participation
	if err := r.db.Preload("Resident").Preload("ActivitySchedule").
		Where("as_id = ? AND resident_id IN ?", asID, residentIDs).
		Find(&participations).Error; err != nil {
		return nil, err
	}

	return participations, nil
}

func (r *GormActivityRepository) GetAllParticipations() ([]*entities.Participation, error) {
	var participations []*entities.Participation
	if err := r.db.Preload("Resident").Preload("ActivitySchedule").Find(&participations).Error; err != nil {
		return nil, err
	}

	return participations, nil
}

func (r *GormActivityRepository) GetResidentsByScheduleIDCustom(asID string, params activityModels.ResidentsByScheduleQueryParams) ([]*entities.Participation, error) {
	var participations []*entities.Participation

	query := r.db.Model(&entities.Participation{}).
		Joins("JOIN residents ON residents.id = participations.resident_id").
		Joins("JOIN rooms ON rooms.id = residents.room_id").
		Where("participations.as_id = ?", asID)

	if params.Search != nil && *params.Search != "" {
		like := "%" + *params.Search + "%"
		query = query.Where(
			"residents.first_name ILIKE ? OR residents.last_name ILIKE ? OR residents.nickname ILIKE ?",
			like, like, like,
		)
	}

	if params.Floor != nil {
		query = query.Where("rooms.floor = ?", *params.Floor)
	}

	if len(params.LabelIDs) > 0 {
		residentIDsSubQuery := r.db.
			Table("resident_labels").
			Select("resident_labels.resident_id").
			Where("resident_labels.label_id IN ?", params.LabelIDs).
			Group("resident_labels.resident_id").
			Having("COUNT(DISTINCT resident_labels.label_id) = ?", len(params.LabelIDs))

		query = query.Where("participations.resident_id IN (?)", residentIDsSubQuery)
	}

	if err := query.
		Preload("Resident.Room").
		Preload("Resident.ResidentLabels.IntakeLabel").
		Order("rooms.floor ASC").
		Order("rooms.room_number ASC").
		Order("residents.first_name ASC").
		Find(&participations).Error; err != nil {
		return nil, err
	}

	return participations, nil
}

func (r *GormActivityRepository) UpdateParticipation(participation *entities.Participation) (*entities.Participation, error) {
	updatePayload := entities.Participation{
		IsParticipating: participation.IsParticipating,
		ImgURLs:         participation.ImgURLs,
	}

	if err := r.db.Model(&entities.Participation{}).
		Where("resident_id = ? AND as_id = ?", participation.ResidentID, participation.ASID).
		Select("is_participating", "img_urls").
		Updates(&updatePayload).Error; err != nil {
		return nil, err
	}

	return r.GetParticipationByResidentIDAndASID(participation.ResidentID, participation.ASID)
}

func (r *GormActivityRepository) UpdateParticipationIsParticipatingByResidentIDs(asID string, residentIDs []string, isParticipating bool) error {
	if err := r.db.Model(&entities.Participation{}).
		Where("as_id = ? AND resident_id IN ?", asID, residentIDs).
		Update("is_participating", isParticipating).Error; err != nil {
		return err
	}

	return nil
}

func (r *GormActivityRepository) DeleteParticipation(residentID, asID string) error {
	if err := r.db.Where("resident_id = ? AND as_id = ?", residentID, asID).
		Delete(&entities.Participation{}).Error; err != nil {
		return err
	}

	return nil
}
