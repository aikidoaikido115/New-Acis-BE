package usecases

import (
	"errors"
	"strings"

	activityModels "github.com/aikidoaikido115/New-Acis-BE/modules/activity/models"
	activityRepo "github.com/aikidoaikido115/New-Acis-BE/modules/activity/repositories"
	"github.com/aikidoaikido115/New-Acis-BE/modules/entities"
	userRepo "github.com/aikidoaikido115/New-Acis-BE/modules/user/repositories"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

var ErrActivityNotFound = errors.New("activity not found")
var ErrStaffProfileNotFound = errors.New("staff profile not found for this user")

type ActivityUsecase interface {
	CreateActivity(req activityModels.CreateActivityRequest, userID string) (*entities.Activity, error)
	GetActivityByID(id string) (*entities.Activity, error)
	GetAllActivities() ([]*entities.Activity, error)
	UpdateActivityByID(id string, req activityModels.UpdateActivityRequest) (*entities.Activity, error)
	DeleteActivityByID(id string) error
}

type ActivityUseCaseImpl struct {
	repo     activityRepo.ActivityRepository
	userRepo userRepo.UserRepository
}

func NewActivityUseCase(repo activityRepo.ActivityRepository, userRepo userRepo.UserRepository) ActivityUsecase {
	return &ActivityUseCaseImpl{
		repo:     repo,
		userRepo: userRepo,
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

	return uc.repo.CreateActivity(activity)
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

func (uc *ActivityUseCaseImpl) UpdateActivityByID(id string, req activityModels.UpdateActivityRequest) (*entities.Activity, error) {
	id = strings.TrimSpace(id)
	if id == "" {
		return nil, errors.New("activity id is required")
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

	return uc.repo.UpdateActivity(existingActivity)
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
