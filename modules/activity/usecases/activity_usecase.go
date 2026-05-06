package usecases

import (
	"encoding/json"
	"errors"
	"io"
	"log"
	"mime/multipart"
	"strings"
	"time"

	"github.com/aikidoaikido115/New-Acis-BE/configs"
	activityModels "github.com/aikidoaikido115/New-Acis-BE/modules/activity/models"
	activityRepo "github.com/aikidoaikido115/New-Acis-BE/modules/activity/repositories"
	audit_constants "github.com/aikidoaikido115/New-Acis-BE/modules/audit_logs/constants"
	auditRepo "github.com/aikidoaikido115/New-Acis-BE/modules/audit_logs/repositories"
	"github.com/aikidoaikido115/New-Acis-BE/modules/entities"
	userRepo "github.com/aikidoaikido115/New-Acis-BE/modules/user/repositories"
	"github.com/aikidoaikido115/New-Acis-BE/pkg/utils"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

var ErrActivityNotFound = errors.New("activity not found")
var ErrActivityScheduleNotFound = errors.New("activity schedule not found")
var ErrParticipationNotFound = errors.New("participation not found")
var ErrActivityAlreadyExists = errors.New("activity already exists, please do not duplicate")
var ErrParticipationAlreadyExists = errors.New("participation already exists")
var ErrStaffProfileNotFound = errors.New("staff profile not found for this user")

type ActivityUsecase interface {

	// Activity operations
	CreateActivity(req activityModels.CreateActivityRequest, userID string) (*entities.Activity, error)
	GetActivityByID(id string) (*entities.Activity, error)
	GetAllActivities() ([]*entities.Activity, error)
	UpdateActivityByID(id string, req activityModels.UpdateActivityRequest, userID string) (*entities.Activity, error)
	DeleteActivityByID(id string) error
	CreateActivityScheduleWithActivitySync(req activityModels.CreateActivityScheduleWithActivitySyncRequest, userID string) (*entities.ActivitySchedule, error)
	UpdateActivityScheduleWithActivitySyncByID(id string, req activityModels.UpdateActivityScheduleWithActivitySyncRequest, userID string) (*entities.ActivitySchedule, error)
	GetActivityScheduleWithActivitySyncByID(id string) (*activityModels.ActivityScheduleWithActivitySyncResponse, error)

	// Activity Schedule operations
	CreateActivitySchedule(req activityModels.CreateActivityScheduleRequest, userID string) (*entities.ActivitySchedule, error)
	GetActivityScheduleByID(id string) (*entities.ActivitySchedule, error)
	GetAllActivitySchedules() ([]*entities.ActivitySchedule, error)
	GetAllActivitySchedulesWithActivitySync(date *time.Time) ([]*activityModels.ActivityScheduleWithActivitySyncResponse, error)
	UpdateActivityScheduleByID(id string, req activityModels.UpdateActivityScheduleRequest, userID string) (*entities.ActivitySchedule, error)
	DeleteActivityScheduleByID(id string) error
	// Participation operations
	CreateParticipation(req activityModels.CreateParticipationRequest, userID string, files []*multipart.FileHeader) (*entities.Participation, error)
	GetParticipationByResidentIDAndASID(residentID, asID string) (*entities.Participation, error)
	GetAllParticipations() ([]*entities.Participation, error)
	GetResidentsByScheduleIDCustom(asID string, params activityModels.ResidentsByScheduleQueryParams) (*activityModels.ResidentsByScheduleListResponse, error)
	UpdateParticipationByResidentIDAndASID(residentID, asID string, req activityModels.UpdateParticipationRequest, userID string, files []*multipart.FileHeader) (*entities.Participation, error)
	BulkUpdateParticipationIsParticipatingByResidentIDs(req activityModels.BulkUpdateParticipationIsParticipatingByResidentIDsRequest, userID string) ([]*entities.Participation, error)
	DeleteParticipationByResidentIDAndASID(residentID, asID string) error
}

type ActivityUseCaseImpl struct {
	repo         activityRepo.ActivityRepository
	userRepo     userRepo.UserRepository
	auditlogrepo auditRepo.AuditLogRepository
	supa         configs.Supabase
}

func NewActivityUseCase(repo activityRepo.ActivityRepository, userRepo userRepo.UserRepository, auditlogrepo auditRepo.AuditLogRepository, supa configs.Supabase) ActivityUsecase {
	return &ActivityUseCaseImpl{
		repo:         repo,
		userRepo:     userRepo,
		auditlogrepo: auditlogrepo,
		supa:         supa,
	}
}

func (uc *ActivityUseCaseImpl) CreateActivity(req activityModels.CreateActivityRequest, userID string) (*entities.Activity, error) {
	userID = strings.TrimSpace(userID)
	if userID == "" {
		return nil, errors.New("user id is required")
	}

	staff, err := uc.userRepo.GetStaffByUserID(userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrStaffProfileNotFound
		}
		return nil, err
	}

	activityName := strings.TrimSpace(req.ActivityName)
	activityType := strings.TrimSpace(req.ActivityType)

	if activityName == "" {
		return nil, errors.New("activity_name is required")
	}
	if activityType == "" {
		return nil, errors.New("activity_type is required")
	}

	activity := &entities.Activity{
		ID:           uuid.New().String(),
		StaffID:      staff.ID,
		ActivityName: activityName,
		ActivityType: activityType,
		Description:  normalizeOptionalString(req.Description),
		Location:     normalizeOptionalString(req.Location),
	}

	createdActivity, err := uc.repo.CreateActivity(activity)
	if err != nil {
		return nil, err
	}

	newValue, _ := json.Marshal(createdActivity)
	uc.createAuditLog(userID, audit_constants.AuditActionInsert, "activities", createdActivity.ID, "", string(newValue))

	return createdActivity, nil
}

func (uc *ActivityUseCaseImpl) GetActivityByID(id string) (*entities.Activity, error) {
	id = strings.TrimSpace(id)
	if id == "" {
		return nil, errors.New("activity id is required")
	}

	activity, err := uc.repo.GetActivityByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrActivityNotFound
		}
		return nil, err
	}

	return activity, nil
}

func (uc *ActivityUseCaseImpl) GetAllActivities() ([]*entities.Activity, error) {
	return uc.repo.GetAllActivities()
}

func (uc *ActivityUseCaseImpl) UpdateActivityByID(id string, req activityModels.UpdateActivityRequest, userID string) (*entities.Activity, error) {
	id = strings.TrimSpace(id)
	if id == "" {
		return nil, errors.New("activity id is required")
	}

	userID = strings.TrimSpace(userID)
	if userID == "" {
		return nil, errors.New("user id is required")
	}

	if req.StaffID == nil && req.ActivityName == nil && req.ActivityType == nil && req.Description == nil && req.Location == nil {
		return nil, errors.New("at least one field must be provided")
	}

	existingActivity, err := uc.repo.GetActivityByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrActivityNotFound
		}
		return nil, err
	}

	oldValue, _ := json.Marshal(existingActivity)

	staff, staffErr := uc.userRepo.GetStaffByUserID(userID)
	if staffErr == nil {
		existingActivity.StaffID = staff.ID
	}

	if req.StaffID != nil {
		staffID := strings.TrimSpace(*req.StaffID)
		if staffID == "" {
			return nil, errors.New("staff_id cannot be empty")
		}
		existingActivity.StaffID = staffID
	}

	if req.ActivityName != nil {
		activityName := strings.TrimSpace(*req.ActivityName)
		if activityName == "" {
			return nil, errors.New("activity_name cannot be empty")
		}
		existingActivity.ActivityName = activityName
	}

	if req.ActivityType != nil {
		activityType := strings.TrimSpace(*req.ActivityType)
		if activityType == "" {
			return nil, errors.New("activity_type cannot be empty")
		}
		existingActivity.ActivityType = activityType
	}

	if req.Description != nil {
		existingActivity.Description = normalizeOptionalString(req.Description)
	}

	if req.Location != nil {
		existingActivity.Location = normalizeOptionalString(req.Location)
	}

	updatedActivity, err := uc.repo.UpdateActivity(existingActivity)
	if err != nil {
		return nil, err
	}

	newValue, _ := json.Marshal(updatedActivity)
	uc.createAuditLog(userID, audit_constants.AuditActionUpdate, "activities", updatedActivity.ID, string(oldValue), string(newValue))

	return uc.repo.GetActivityByID(id)
}

func (uc *ActivityUseCaseImpl) DeleteActivityByID(id string) error {
	id = strings.TrimSpace(id)
	if id == "" {
		return errors.New("activity id is required")
	}

	if _, err := uc.repo.GetActivityByID(id); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrActivityNotFound
		}
		return err
	}

	return uc.repo.DeleteActivity(id)
}

func (uc *ActivityUseCaseImpl) CreateActivityScheduleWithActivitySync(req activityModels.CreateActivityScheduleWithActivitySyncRequest, userID string) (*entities.ActivitySchedule, error) {
	userID = strings.TrimSpace(userID)
	if userID == "" {
		return nil, errors.New("user id is required")
	}

	activityName := strings.TrimSpace(req.ActivityName)
	activityType := strings.TrimSpace(req.ActivityType)
	if activityName == "" {
		return nil, errors.New("activity_name is required")
	}
	if activityType == "" {
		return nil, errors.New("activity_type is required")
	}

	if req.Date.IsZero() {
		return nil, errors.New("date is required")
	}
	if req.StartTime.IsZero() {
		return nil, errors.New("start_time is required")
	}
	if req.EndTime.IsZero() {
		return nil, errors.New("end_time is required")
	}
	if !req.EndTime.After(req.StartTime) {
		return nil, errors.New("end_time must be after start_time")
	}

	existingActivity, err := uc.repo.GetActivityByName(activityName)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	if errors.Is(err, gorm.ErrRecordNotFound) {
		createdActivity, createErr := uc.CreateActivity(activityModels.CreateActivityRequest{
			ActivityName: activityName,
			ActivityType: activityType,
			Description:  req.Description,
			Location:     req.Location,
		}, userID)
		if createErr != nil {
			return nil, createErr
		}

		return uc.CreateActivitySchedule(activityModels.CreateActivityScheduleRequest{
			ActivityID: createdActivity.ID,
			Date:       req.Date,
			StartTime:  req.StartTime,
			EndTime:    req.EndTime,
		}, userID)
	}

	newDescription := normalizeOptionalString(req.Description)
	currentDescription := normalizeOptionalString(existingActivity.Description)
	newLocation := normalizeOptionalString(req.Location)
	currentLocation := normalizeOptionalString(existingActivity.Location)

	updateReq := activityModels.UpdateActivityRequest{}
	needUpdate := false

	staff, staffErr := uc.userRepo.GetStaffByUserID(userID)
	if staffErr == nil {
		staffIDStr := staff.ID
		updateReq.StaffID = &staffIDStr
		needUpdate = true
	}

	if existingActivity.ActivityType != activityType {
		updateReq.ActivityType = &activityType
		needUpdate = true
	}

	if !optionalStringEqual(currentDescription, newDescription) {
		if newDescription == nil {
			empty := ""
			updateReq.Description = &empty
		} else {
			updateReq.Description = newDescription
		}
		needUpdate = true
	}

	if !optionalStringEqual(currentLocation, newLocation) {
		if newLocation == nil {
			empty := ""
			updateReq.Location = &empty
		} else {
			updateReq.Location = newLocation
		}
		needUpdate = true
	}

	resolvedActivityID := existingActivity.ID
	if needUpdate {
		updatedActivity, updateErr := uc.UpdateActivityByID(existingActivity.ID, updateReq, userID)
		if updateErr != nil {
			return nil, updateErr
		}
		resolvedActivityID = updatedActivity.ID
	}

	return uc.CreateActivitySchedule(activityModels.CreateActivityScheduleRequest{
		ActivityID: resolvedActivityID,
		Date:       req.Date,
		StartTime:  req.StartTime,
		EndTime:    req.EndTime,
	}, userID)
}

func (uc *ActivityUseCaseImpl) UpdateActivityScheduleWithActivitySyncByID(id string, req activityModels.UpdateActivityScheduleWithActivitySyncRequest, userID string) (*entities.ActivitySchedule, error) {
	id = strings.TrimSpace(id)
	if id == "" {
		return nil, errors.New("activity schedule id is required")
	}

	userID = strings.TrimSpace(userID)
	if userID == "" {
		return nil, errors.New("user id is required")
	}

	if req.ActivityName == nil && req.ActivityType == nil && req.Date == nil && req.StartTime == nil && req.EndTime == nil && req.Location == nil && req.Description == nil {
		return nil, errors.New("at least one field must be provided")
	}

	existingSchedule, err := uc.repo.GetActivityScheduleByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrActivityScheduleNotFound
		}
		return nil, err
	}

	existingActivity, err := uc.repo.GetActivityByID(existingSchedule.ActivityID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrActivityNotFound
		}
		return nil, err
	}

	activityUpdateReq := activityModels.UpdateActivityRequest{}
	needUpdateActivity := false

	staff, staffErr := uc.userRepo.GetStaffByUserID(userID)
	if staffErr == nil {
		staffIDStr := staff.ID
		activityUpdateReq.StaffID = &staffIDStr
		needUpdateActivity = true
	}

	if req.ActivityName != nil {
		activityName := strings.TrimSpace(*req.ActivityName)
		if activityName == "" {
			return nil, errors.New("activity_name cannot be empty")
		}

		if !strings.EqualFold(existingActivity.ActivityName, activityName) {
			duplicated, duplicateErr := uc.repo.GetActivityByName(activityName)
			if duplicateErr == nil && duplicated.ID != existingActivity.ID {
				return nil, ErrActivityAlreadyExists
			}
			if duplicateErr != nil && !errors.Is(duplicateErr, gorm.ErrRecordNotFound) {
				return nil, duplicateErr
			}
		}

		activityUpdateReq.ActivityName = &activityName
		needUpdateActivity = true
	}

	if req.ActivityType != nil {
		activityType := strings.TrimSpace(*req.ActivityType)
		if activityType == "" {
			return nil, errors.New("activity_type cannot be empty")
		}
		activityUpdateReq.ActivityType = &activityType
		needUpdateActivity = true
	}

	if req.Description != nil {
		activityUpdateReq.Description = req.Description
		needUpdateActivity = true
	}

	if req.Location != nil {
		activityUpdateReq.Location = req.Location
		needUpdateActivity = true
	}

	if needUpdateActivity {
		if _, err := uc.UpdateActivityByID(existingActivity.ID, activityUpdateReq, userID); err != nil {
			return nil, err
		}
	}

	scheduleUpdateReq := activityModels.UpdateActivityScheduleRequest{}
	needUpdateSchedule := false

	if req.Date != nil {
		scheduleUpdateReq.Date = req.Date
		needUpdateSchedule = true
	}
	if req.StartTime != nil {
		scheduleUpdateReq.StartTime = req.StartTime
		needUpdateSchedule = true
	}
	if req.EndTime != nil {
		scheduleUpdateReq.EndTime = req.EndTime
		needUpdateSchedule = true
	}

	if needUpdateSchedule {
		return uc.UpdateActivityScheduleByID(existingSchedule.ID, scheduleUpdateReq, userID)
	}

	return uc.repo.GetActivityScheduleByID(existingSchedule.ID)
}

func (uc *ActivityUseCaseImpl) GetActivityScheduleWithActivitySyncByID(id string) (*activityModels.ActivityScheduleWithActivitySyncResponse, error) {
	id = strings.TrimSpace(id)
	if id == "" {
		return nil, errors.New("activity schedule id is required")
	}

	activitySchedule, err := uc.repo.GetActivityScheduleByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrActivityScheduleNotFound
		}
		return nil, err
	}

	return &activityModels.ActivityScheduleWithActivitySyncResponse{
		ActivityName: activitySchedule.Activity.ActivityName,
		ActivityType: activitySchedule.Activity.ActivityType,
		Date:         activitySchedule.Date,
		StartTime:    activitySchedule.StartTime,
		EndTime:      activitySchedule.EndTime,
		Location:     activitySchedule.Activity.Location,
		Description:  activitySchedule.Activity.Description,
	}, nil
}

func (uc *ActivityUseCaseImpl) CreateActivitySchedule(req activityModels.CreateActivityScheduleRequest, userID string) (*entities.ActivitySchedule, error) {
	userID = strings.TrimSpace(userID)
	if userID == "" {
		return nil, errors.New("user id is required")
	}

	activityID := strings.TrimSpace(req.ActivityID)
	if activityID == "" {
		return nil, errors.New("activity_id is required")
	}

	if req.Date.IsZero() {
		return nil, errors.New("date is required")
	}
	if req.StartTime.IsZero() {
		return nil, errors.New("start_time is required")
	}
	if req.EndTime.IsZero() {
		return nil, errors.New("end_time is required")
	}
	if !req.EndTime.After(req.StartTime) {
		return nil, errors.New("end_time must be after start_time")
	}

	if _, err := uc.repo.GetActivityByID(activityID); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrActivityNotFound
		}
		return nil, err
	}

	now := time.Now()
	activitySchedule := &entities.ActivitySchedule{
		ID:         uuid.New().String(),
		ActivityID: activityID,
		Date:       req.Date,
		StartTime:  req.StartTime,
		EndTime:    req.EndTime,
		CreatedAt:  now,
		UpdatedAt:  now,
	}

	createdActivitySchedule, err := uc.repo.CreateActivityScheduleWithDefaultParticipations(activitySchedule)
	if err != nil {
		return nil, err
	}

	newValue, _ := json.Marshal(createdActivitySchedule)
	uc.createAuditLog(userID, audit_constants.AuditActionInsert, "activity_schedules", createdActivitySchedule.ID, "", string(newValue))

	return createdActivitySchedule, nil
}

func (uc *ActivityUseCaseImpl) GetActivityScheduleByID(id string) (*entities.ActivitySchedule, error) {
	id = strings.TrimSpace(id)
	if id == "" {
		return nil, errors.New("activity schedule id is required")
	}

	activitySchedule, err := uc.repo.GetActivityScheduleByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrActivityScheduleNotFound
		}
		return nil, err
	}

	return activitySchedule, nil
}

func (uc *ActivityUseCaseImpl) GetAllActivitySchedules() ([]*entities.ActivitySchedule, error) {
	return uc.repo.GetAllActivitySchedules()
}

func (uc *ActivityUseCaseImpl) GetAllActivitySchedulesWithActivitySync(date *time.Time) ([]*activityModels.ActivityScheduleWithActivitySyncResponse, error) {
	activitySchedules, err := uc.repo.GetActivitySchedulesWithActivitySyncByDate(date)
	if err != nil {
		return nil, err
	}

	responses := make([]*activityModels.ActivityScheduleWithActivitySyncResponse, 0, len(activitySchedules))
	for _, schedule := range activitySchedules {
		responses = append(responses, &activityModels.ActivityScheduleWithActivitySyncResponse{
			ActivityName: schedule.Activity.ActivityName,
			ActivityType: schedule.Activity.ActivityType,
			Date:         schedule.Date,
			StartTime:    schedule.StartTime,
			EndTime:      schedule.EndTime,
			Location:     schedule.Activity.Location,
			Description:  schedule.Activity.Description,
		})
	}

	return responses, nil
}

func (uc *ActivityUseCaseImpl) UpdateActivityScheduleByID(id string, req activityModels.UpdateActivityScheduleRequest, userID string) (*entities.ActivitySchedule, error) {
	id = strings.TrimSpace(id)
	if id == "" {
		return nil, errors.New("activity schedule id is required")
	}

	userID = strings.TrimSpace(userID)
	if userID == "" {
		return nil, errors.New("user id is required")
	}

	if req.ActivityID == nil && req.Date == nil && req.StartTime == nil && req.EndTime == nil {
		return nil, errors.New("at least one field must be provided")
	}

	existingActivitySchedule, err := uc.repo.GetActivityScheduleByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrActivityScheduleNotFound
		}
		return nil, err
	}

	oldValue, _ := json.Marshal(existingActivitySchedule)

	if req.ActivityID != nil {
		activityID := strings.TrimSpace(*req.ActivityID)
		if activityID == "" {
			return nil, errors.New("activity_id cannot be empty")
		}

		if _, err := uc.repo.GetActivityByID(activityID); err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return nil, ErrActivityNotFound
			}
			return nil, err
		}

		existingActivitySchedule.ActivityID = activityID
	}

	if req.Date != nil {
		if req.Date.IsZero() {
			return nil, errors.New("date cannot be empty")
		}
		existingActivitySchedule.Date = *req.Date
	}

	candidateStartTime := existingActivitySchedule.StartTime
	if req.StartTime != nil {
		if req.StartTime.IsZero() {
			return nil, errors.New("start_time cannot be empty")
		}
		candidateStartTime = *req.StartTime
	}

	candidateEndTime := existingActivitySchedule.EndTime
	if req.EndTime != nil {
		if req.EndTime.IsZero() {
			return nil, errors.New("end_time cannot be empty")
		}
		candidateEndTime = *req.EndTime
	}

	if !candidateEndTime.After(candidateStartTime) {
		return nil, errors.New("end_time must be after start_time")
	}

	staff, staffErr := uc.userRepo.GetStaffByUserID(userID)
	if staffErr == nil {
		staffIDStr := staff.ID
		activityUpdateReq := activityModels.UpdateActivityRequest{StaffID: &staffIDStr}
		if _, err := uc.UpdateActivityByID(existingActivitySchedule.ActivityID, activityUpdateReq, userID); err != nil {
			return nil, err
		}
	}

	existingActivitySchedule.StartTime = candidateStartTime
	existingActivitySchedule.EndTime = candidateEndTime
	existingActivitySchedule.UpdatedAt = time.Now()

	updatedActivitySchedule, err := uc.repo.UpdateActivitySchedule(existingActivitySchedule)
	if err != nil {
		return nil, err
	}

	newValue, _ := json.Marshal(updatedActivitySchedule)
	uc.createAuditLog(userID, audit_constants.AuditActionUpdate, "activity_schedules", updatedActivitySchedule.ID, string(oldValue), string(newValue))

	return updatedActivitySchedule, nil
}

func (uc *ActivityUseCaseImpl) DeleteActivityScheduleByID(id string) error {
	id = strings.TrimSpace(id)
	if id == "" {
		return errors.New("activity schedule id is required")
	}

	if _, err := uc.repo.GetActivityScheduleByID(id); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrActivityScheduleNotFound
		}
		return err
	}

	return uc.repo.DeleteActivitySchedule(id)
}

func (uc *ActivityUseCaseImpl) CreateParticipation(req activityModels.CreateParticipationRequest, userID string, files []*multipart.FileHeader) (*entities.Participation, error) {
	userID = strings.TrimSpace(userID)
	if userID == "" {
		return nil, errors.New("user id is required")
	}

	residentID := strings.TrimSpace(req.ResidentID)
	asID := strings.TrimSpace(req.ASID)
	if residentID == "" {
		return nil, errors.New("resident_id is required")
	}
	if asID == "" {
		return nil, errors.New("as_id is required")
	}

	if _, err := uc.repo.GetParticipationByResidentIDAndASID(residentID, asID); err == nil {
		return nil, ErrParticipationAlreadyExists
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	imgURLs, err := uc.uploadParticipationImages(files)
	if err != nil {
		return nil, err
	}

	participation := &entities.Participation{
		ResidentID:      residentID,
		ASID:            asID,
		IsParticipating: req.IsParticipating,
		ImgURLs:         imgURLs,
	}

	createdParticipation, err := uc.repo.CreateParticipation(participation)
	if err != nil {
		return nil, err
	}

	existingSchedule, err := uc.repo.GetActivityScheduleByID(asID)
	if err == nil {
		_, _ = uc.UpdateActivityByID(existingSchedule.ActivityID, activityModels.UpdateActivityRequest{}, userID)
	}

	newValue, _ := json.Marshal(createdParticipation)
	uc.createAuditLog(userID, audit_constants.AuditActionInsert, "participations", residentID+"-"+asID, "", string(newValue))

	return createdParticipation, nil
}

func (uc *ActivityUseCaseImpl) GetParticipationByResidentIDAndASID(residentID, asID string) (*entities.Participation, error) {
	residentID = strings.TrimSpace(residentID)
	asID = strings.TrimSpace(asID)
	if residentID == "" {
		return nil, errors.New("resident_id is required")
	}
	if asID == "" {
		return nil, errors.New("as_id is required")
	}

	participation, err := uc.repo.GetParticipationByResidentIDAndASID(residentID, asID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrParticipationNotFound
		}
		return nil, err
	}

	return participation, nil
}

func (uc *ActivityUseCaseImpl) GetAllParticipations() ([]*entities.Participation, error) {
	return uc.repo.GetAllParticipations()
}

func (uc *ActivityUseCaseImpl) GetResidentsByScheduleIDCustom(asID string, params activityModels.ResidentsByScheduleQueryParams) (*activityModels.ResidentsByScheduleListResponse, error) {
	asID = strings.TrimSpace(asID)
	if asID == "" {
		return nil, errors.New("activity schedule id is required")
	}

	page := 1
	if params.Page != nil {
		if *params.Page <= 0 {
			return nil, errors.New("page must be greater than 0")
		}
		page = *params.Page
	}

	pageSize := 20
	if params.PageSize != nil {
		if *params.PageSize <= 0 {
			return nil, errors.New("page_size must be greater than 0")
		}
		pageSize = *params.PageSize
	} else if params.Limit > 0 {
		pageSize = params.Limit
	}
	if pageSize > 100 {
		pageSize = 100
	}

	hasPagination := params.Page != nil || params.PageSize != nil || params.Limit > 0 || params.Offset > 0
	if params.Page == nil && params.Offset > 0 {
		page = (params.Offset / pageSize) + 1
	}

	if hasPagination {
		params.Limit = pageSize
		params.Offset = (page - 1) * pageSize
	}

	participations, total, err := uc.repo.GetResidentsByScheduleIDCustom(asID, params)
	if err != nil {
		return nil, err
	}

	result := make([]*activityModels.ResidentByScheduleResponse, 0, len(participations))
	for _, p := range participations {
		labels := make([]string, 0, len(p.Resident.ResidentLabels))
		seen := make(map[string]struct{})
		for _, residentLabel := range p.Resident.ResidentLabels {
			labelName := strings.TrimSpace(residentLabel.IntakeLabel.LabelName)
			if labelName == "" {
				continue
			}
			if _, exists := seen[labelName]; exists {
				continue
			}
			seen[labelName] = struct{}{}
			labels = append(labels, labelName)
		}

		result = append(result, &activityModels.ResidentByScheduleResponse{
			ResidentID:      p.ResidentID,
			FirstName:       p.Resident.FirstName,
			LastName:        p.Resident.LastName,
			Nickname:        p.Resident.Nickname,
			RoomNumber:      p.Resident.Room.RoomNumber,
			Floor:           p.Resident.Room.Floor,
			IntakeLabels:    labels,
			IsParticipating: p.IsParticipating,
		})
	}

	totalPages := 0
	if total > 0 {
		totalPages = int((total + int64(pageSize) - 1) / int64(pageSize))
	}

	return &activityModels.ResidentsByScheduleListResponse{
		Items: result,
		Pagination: activityModels.ActivityPagination{
			Page:       page,
			PageSize:   pageSize,
			TotalItems: int(total),
			TotalPages: totalPages,
		},
	}, nil
}

func (uc *ActivityUseCaseImpl) UpdateParticipationByResidentIDAndASID(residentID, asID string, req activityModels.UpdateParticipationRequest, userID string, files []*multipart.FileHeader) (*entities.Participation, error) {
	residentID = strings.TrimSpace(residentID)
	asID = strings.TrimSpace(asID)
	if residentID == "" {
		return nil, errors.New("resident_id is required")
	}
	if asID == "" {
		return nil, errors.New("as_id is required")
	}

	userID = strings.TrimSpace(userID)
	if userID == "" {
		return nil, errors.New("user id is required")
	}

	if req.IsParticipating == nil && len(files) == 0 {
		return nil, errors.New("at least one field must be provided")
	}

	existingParticipation, err := uc.repo.GetParticipationByResidentIDAndASID(residentID, asID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrParticipationNotFound
		}
		return nil, err
	}

	if existingParticipation.ResidentID != residentID || existingParticipation.ASID != asID {
		return nil, errors.New("resident_id and as_id are immutable")
	}

	oldValue, _ := json.Marshal(existingParticipation)

	if req.IsParticipating != nil {
		existingParticipation.IsParticipating = *req.IsParticipating
	}

	if len(files) > 0 {
		imgURLs, err := uc.uploadParticipationImages(files)
		if err != nil {
			return nil, err
		}
		existingParticipation.ImgURLs = imgURLs
	}

	// Enforce composite key from path even if entity is mutated elsewhere.
	existingParticipation.ResidentID = residentID
	existingParticipation.ASID = asID

	updatedParticipation, err := uc.repo.UpdateParticipation(existingParticipation)
	if err != nil {
		return nil, err
	}

	existingSchedule, err := uc.repo.GetActivityScheduleByID(asID)
	if err == nil {
		_, _ = uc.UpdateActivityByID(existingSchedule.ActivityID, activityModels.UpdateActivityRequest{}, userID)
	}

	newValue, _ := json.Marshal(updatedParticipation)
	uc.createAuditLog(userID, audit_constants.AuditActionUpdate, "participations", residentID+"-"+asID, string(oldValue), string(newValue))

	return updatedParticipation, nil
}

func (uc *ActivityUseCaseImpl) DeleteParticipationByResidentIDAndASID(residentID, asID string) error {
	residentID = strings.TrimSpace(residentID)
	asID = strings.TrimSpace(asID)
	if residentID == "" {
		return errors.New("resident_id is required")
	}
	if asID == "" {
		return errors.New("as_id is required")
	}

	if _, err := uc.repo.GetParticipationByResidentIDAndASID(residentID, asID); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrParticipationNotFound
		}
		return err
	}

	return uc.repo.DeleteParticipation(residentID, asID)
}

func (uc *ActivityUseCaseImpl) BulkUpdateParticipationIsParticipatingByResidentIDs(req activityModels.BulkUpdateParticipationIsParticipatingByResidentIDsRequest, userID string) ([]*entities.Participation, error) {
	userID = strings.TrimSpace(userID)
	if userID == "" {
		return nil, errors.New("user id is required")
	}

	asID := strings.TrimSpace(req.ASID)
	if asID == "" {
		return nil, errors.New("as_id is required")
	}

	if req.IsParticipating == nil {
		return nil, errors.New("is_participating is required")
	}

	residentIDs := normalizeUniqueResidentIDs(req.ResidentIDs)
	if len(residentIDs) == 0 {
		return nil, errors.New("resident_ids is required")
	}

	existingParticipations, err := uc.repo.GetParticipationsByASIDAndResidentIDs(asID, residentIDs)
	if err != nil {
		return nil, err
	}
	if len(existingParticipations) == 0 {
		return nil, ErrParticipationNotFound
	}

	oldValueByKey := make(map[string]string, len(existingParticipations))
	for _, participation := range existingParticipations {
		key := participation.ResidentID + "-" + participation.ASID
		oldValue, _ := json.Marshal(participation)
		oldValueByKey[key] = string(oldValue)
	}

	if err := uc.repo.UpdateParticipationIsParticipatingByResidentIDs(asID, residentIDs, *req.IsParticipating); err != nil {
		return nil, err
	}

	updatedParticipations, err := uc.repo.GetParticipationsByASIDAndResidentIDs(asID, residentIDs)
	if err != nil {
		return nil, err
	}

	for _, participation := range updatedParticipations {
		key := participation.ResidentID + "-" + participation.ASID
		newValue, _ := json.Marshal(participation)
		uc.createAuditLog(userID, audit_constants.AuditActionUpdate, "participations", key, oldValueByKey[key], string(newValue))
	}

	return updatedParticipations, nil
}

func normalizeOptionalString(value *string) *string {
	if value == nil {
		return nil
	}

	trimmed := strings.TrimSpace(*value)
	if trimmed == "" {
		return nil
	}

	return &trimmed
}

func optionalStringEqual(a *string, b *string) bool {
	if a == nil && b == nil {
		return true
	}
	if a == nil || b == nil {
		return false
	}
	return strings.TrimSpace(*a) == strings.TrimSpace(*b)
}

func (uc *ActivityUseCaseImpl) uploadParticipationImages(files []*multipart.FileHeader) ([]entities.ImageURL, error) {
	if len(files) == 0 {
		return []entities.ImageURL{}, nil
	}

	result := make([]entities.ImageURL, 0, len(files))
	for _, fileHeader := range files {
		file, err := fileHeader.Open()
		if err != nil {
			return nil, errors.New("failed to open file: " + err.Error())
		}

		fileExtension, err := utils.DetectFileType(file)
		if err != nil {
			file.Close()
			return nil, errors.New("invalid file: " + err.Error())
		}

		if _, err = file.Seek(0, io.SeekStart); err != nil {
			file.Close()
			return nil, errors.New("failed to reset file pointer: " + err.Error())
		}

		fileName := uuid.New().String() + fileExtension
		imageURL, err := utils.UploadFile2Supa(file, fileName, "participation/", uc.supa)
		file.Close()
		if err != nil {
			return nil, errors.New("failed to upload participation image: " + err.Error())
		}

		result = append(result, entities.ImageURL{URL: imageURL})
	}

	return result, nil
}

func normalizeUniqueResidentIDs(residentIDs []string) []string {
	seen := make(map[string]struct{}, len(residentIDs))
	result := make([]string, 0, len(residentIDs))
	for _, residentID := range residentIDs {
		trimmed := strings.TrimSpace(residentID)
		if trimmed == "" {
			continue
		}
		if _, exists := seen[trimmed]; exists {
			continue
		}
		seen[trimmed] = struct{}{}
		result = append(result, trimmed)
	}
	return result
}

func (uc *ActivityUseCaseImpl) createAuditLog(userID string, action string, tableName string, recordID string, oldValue string, newValue string) {
	auditLog := &entities.AuditLogs{
		ID:        uuid.New().String(),
		TableName: tableName,
		RecordID:  recordID,
		UserID:    userID,
		Action:    action,
		OldValue:  oldValue,
		NewValue:  newValue,
	}

	if _, err := uc.auditlogrepo.CreateAuditLog(auditLog); err != nil {
		log.Printf("[ERROR] Failed to create audit log for %s %s: %v", tableName, recordID, err)
	}
}
