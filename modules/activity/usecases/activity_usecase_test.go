package usecases_test

import (
	"testing"
	"time"

	"github.com/aikidoaikido115/New-Acis-BE/configs"
	activityModels "github.com/aikidoaikido115/New-Acis-BE/modules/activity/models"
	activityRepo "github.com/aikidoaikido115/New-Acis-BE/modules/activity/repositories"
	"github.com/aikidoaikido115/New-Acis-BE/modules/activity/usecases"
	auditRepo "github.com/aikidoaikido115/New-Acis-BE/modules/audit_logs/repositories"
	"github.com/aikidoaikido115/New-Acis-BE/modules/entities"
	userRepo "github.com/aikidoaikido115/New-Acis-BE/modules/user/repositories"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

type fakeActivityRepo struct {
	activityRepo.ActivityRepository

	existingActivity                 *entities.Activity
	existingActivitySchedule         *entities.ActivitySchedule
	updatedActivity                  *entities.Activity
	updatedActivitySchedule          *entities.ActivitySchedule
	createdActivitySchedule          *entities.ActivitySchedule
	capturedCreateActivitySchedule   *entities.ActivitySchedule
	capturedUpdateResidentIDs        []string
	capturedUpdateASID               string
	capturedUpdateIsParticipating    bool
	capturedResidentsByScheduleASID  string
	participationsFirstFetch         []*entities.Participation
	participationsSecondFetch        []*entities.Participation
	residentsByScheduleResponse      []*entities.Participation
	getActivityByNameCallCount       int
	getActivityByIDCallCount         int
	getActivityScheduleByIDCallCount int
	updateActivityCallCount          int
	updateActivityScheduleCallCount  int
	createActivityScheduleCallCount  int
	getResidentsByScheduleCallCount  int
	getParticipationsCallCount       int
	updateParticipationsCallCount    int
}

func newFakeActivityRepo(existingActivity *entities.Activity, createdActivitySchedule *entities.ActivitySchedule) *fakeActivityRepo {
	return &fakeActivityRepo{
		ActivityRepository:      activityRepo.NewGormActivityRepository(nil),
		existingActivity:        existingActivity,
		createdActivitySchedule: createdActivitySchedule,
	}
}

func (f *fakeActivityRepo) GetActivityByName(activityName string) (*entities.Activity, error) {
	f.getActivityByNameCallCount++
	return f.existingActivity, nil
}

func (f *fakeActivityRepo) GetActivityByID(id string) (*entities.Activity, error) {
	f.getActivityByIDCallCount++
	return f.existingActivity, nil
}

func (f *fakeActivityRepo) CreateActivityScheduleWithDefaultParticipations(activitySchedule *entities.ActivitySchedule) (*entities.ActivitySchedule, error) {
	f.createActivityScheduleCallCount++
	copied := *activitySchedule
	f.capturedCreateActivitySchedule = &copied
	return f.createdActivitySchedule, nil
}

func (f *fakeActivityRepo) UpdateActivity(activity *entities.Activity) (*entities.Activity, error) {
	f.updateActivityCallCount++
	copied := *activity
	f.updatedActivity = &copied
	return activity, nil
}

func (f *fakeActivityRepo) GetActivityScheduleByID(id string) (*entities.ActivitySchedule, error) {
	f.getActivityScheduleByIDCallCount++
	return f.existingActivitySchedule, nil
}

func (f *fakeActivityRepo) UpdateActivitySchedule(activitySchedule *entities.ActivitySchedule) (*entities.ActivitySchedule, error) {
	f.updateActivityScheduleCallCount++
	copied := *activitySchedule
	f.updatedActivitySchedule = &copied
	return activitySchedule, nil
}

func (f *fakeActivityRepo) GetResidentsByScheduleIDCustom(asID string, params activityModels.ResidentsByScheduleQueryParams) ([]*entities.Participation, int64, error) {
	f.getResidentsByScheduleCallCount++
	f.capturedResidentsByScheduleASID = asID
	return f.residentsByScheduleResponse, int64(len(f.residentsByScheduleResponse)), nil
}

func (f *fakeActivityRepo) GetParticipationsByASIDAndResidentIDs(asID string, residentIDs []string) ([]*entities.Participation, error) {
	f.getParticipationsCallCount++
	if f.getParticipationsCallCount == 1 {
		return f.participationsFirstFetch, nil
	}
	return f.participationsSecondFetch, nil
}

func (f *fakeActivityRepo) UpdateParticipationIsParticipatingByResidentIDs(asID string, residentIDs []string, isParticipating bool) error {
	f.updateParticipationsCallCount++
	f.capturedUpdateASID = asID
	f.capturedUpdateResidentIDs = append([]string{}, residentIDs...)
	f.capturedUpdateIsParticipating = isParticipating
	return nil
}

type fakeAuditLogRepo struct {
	auditRepo.AuditLogRepository
	createAuditLogCallCount int
}

type fakeActivityUserRepo struct {
	userRepo.UserRepository
	staffByUserID map[string]*entities.Staff
}

func newFakeActivityUserRepo() *fakeActivityUserRepo {
	return &fakeActivityUserRepo{
		UserRepository: userRepo.NewGormUserRepository(nil),
		staffByUserID: map[string]*entities.Staff{
			"user-1": {
				ID:     "staff-1",
				UserID: "user-1",
			},
		},
	}
}

func (f *fakeActivityUserRepo) GetStaffByUserID(userID string) (*entities.Staff, error) {
	if staff, ok := f.staffByUserID[userID]; ok {
		return staff, nil
	}
	return nil, gorm.ErrRecordNotFound
}

func newFakeAuditLogRepo() *fakeAuditLogRepo {
	return &fakeAuditLogRepo{AuditLogRepository: auditRepo.NewGormAuditLogRepository(nil)}
}

func (f *fakeAuditLogRepo) CreateAuditLog(auditLog *entities.AuditLogs) (*entities.AuditLogs, error) {
	f.createAuditLogCallCount++
	return auditLog, nil
}

func TestCreateActivityScheduleWithActivitySync_Success_ReuseExistingActivity(t *testing.T) {
	location := "Room A"
	description := "Simple stretching"

	date := time.Date(2026, 4, 23, 0, 0, 0, 0, time.UTC)
	startTime := time.Date(2026, 4, 23, 9, 0, 0, 0, time.UTC)
	endTime := time.Date(2026, 4, 23, 10, 0, 0, 0, time.UTC)

	existingActivity := &entities.Activity{
		ID:           "act-1",
		StaffID:      "staff-1",
		ActivityName: "Morning Exercise",
		ActivityType: "Wellness",
		Location:     &location,
		Description:  &description,
	}

	createdSchedule := &entities.ActivitySchedule{
		ID:         "as-1",
		ActivityID: "act-1",
		Date:       date,
		StartTime:  startTime,
		EndTime:    endTime,
	}

	activityRepository := newFakeActivityRepo(existingActivity, createdSchedule)
	auditRepository := newFakeAuditLogRepo()
	activityUserRepo := newFakeActivityUserRepo()
	uc := usecases.NewActivityUseCase(
		activityRepository,
		activityUserRepo,
		auditRepository,
		configs.Supabase{},
	)

	result, err := uc.CreateActivityScheduleWithActivitySync(activityModels.CreateActivityScheduleWithActivitySyncRequest{
		ActivityName: "Morning Exercise",
		ActivityType: "Wellness",
		Date:         date,
		StartTime:    startTime,
		EndTime:      endTime,
		Location:     &location,
		Description:  &description,
	}, "user-1")

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "as-1", result.ID)
	assert.Equal(t, "act-1", result.ActivityID)

	assert.Equal(t, 1, activityRepository.getActivityByNameCallCount)
	assert.Equal(t, 3, activityRepository.getActivityByIDCallCount)
	assert.Equal(t, 1, activityRepository.createActivityScheduleCallCount)

	if assert.NotNil(t, activityRepository.capturedCreateActivitySchedule) {
		assert.Equal(t, "act-1", activityRepository.capturedCreateActivitySchedule.ActivityID)
		assert.True(t, date.Equal(activityRepository.capturedCreateActivitySchedule.Date))
		assert.True(t, startTime.Equal(activityRepository.capturedCreateActivitySchedule.StartTime))
		assert.True(t, endTime.Equal(activityRepository.capturedCreateActivitySchedule.EndTime))
	}

	assert.Equal(t, 2, auditRepository.createAuditLogCallCount)
}

func TestCreateActivityScheduleWithActivitySync_Success_UpdateExistingActivityWhenFieldsChanged(t *testing.T) {
	oldLocation := "Room A"
	oldDescription := "Simple stretching"
	newLocation := "Room B"

	date := time.Date(2026, 4, 23, 0, 0, 0, 0, time.UTC)
	startTime := time.Date(2026, 4, 23, 9, 0, 0, 0, time.UTC)
	endTime := time.Date(2026, 4, 23, 10, 0, 0, 0, time.UTC)

	existingActivity := &entities.Activity{
		ID:           "act-1",
		StaffID:      "staff-1",
		ActivityName: "Morning Exercise",
		ActivityType: "Wellness",
		Location:     &oldLocation,
		Description:  &oldDescription,
	}

	createdSchedule := &entities.ActivitySchedule{
		ID:         "as-2",
		ActivityID: "act-1",
		Date:       date,
		StartTime:  startTime,
		EndTime:    endTime,
	}

	activityRepository := newFakeActivityRepo(existingActivity, createdSchedule)
	auditRepository := newFakeAuditLogRepo()
	activityUserRepo := newFakeActivityUserRepo()
	uc := usecases.NewActivityUseCase(
		activityRepository,
		activityUserRepo,
		auditRepository,
		configs.Supabase{},
	)

	result, err := uc.CreateActivityScheduleWithActivitySync(activityModels.CreateActivityScheduleWithActivitySyncRequest{
		ActivityName: "Morning Exercise",
		ActivityType: "Rehab",
		Date:         date,
		StartTime:    startTime,
		EndTime:      endTime,
		Location:     &newLocation,
		Description:  nil,
	}, "user-1")

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "as-2", result.ID)
	assert.Equal(t, "act-1", result.ActivityID)

	assert.Equal(t, 1, activityRepository.updateActivityCallCount)
	if assert.NotNil(t, activityRepository.updatedActivity) {
		assert.Equal(t, "Rehab", activityRepository.updatedActivity.ActivityType)
		assert.Nil(t, activityRepository.updatedActivity.Description)
		if assert.NotNil(t, activityRepository.updatedActivity.Location) {
			assert.Equal(t, "Room B", *activityRepository.updatedActivity.Location)
		}
	}

	assert.Equal(t, 2, auditRepository.createAuditLogCallCount)
}

func TestBulkUpdateParticipationIsParticipatingByResidentIDs_Success_DedupResidentIDs(t *testing.T) {
	isParticipating := true

	before := []*entities.Participation{
		{ResidentID: "r1", ASID: "as-1", IsParticipating: false},
		{ResidentID: "r2", ASID: "as-1", IsParticipating: false},
	}
	after := []*entities.Participation{
		{ResidentID: "r1", ASID: "as-1", IsParticipating: true},
		{ResidentID: "r2", ASID: "as-1", IsParticipating: true},
	}

	activityRepository := newFakeActivityRepo(nil, nil)
	activityRepository.existingActivity = &entities.Activity{
		ID:           "act-1",
		StaffID:      "staff-old",
		ActivityName: "Morning Exercise",
		ActivityType: "Wellness",
	}
	activityRepository.existingActivitySchedule = &entities.ActivitySchedule{
		ID:         "as-1",
		ActivityID: "act-1",
	}
	activityRepository.participationsFirstFetch = before
	activityRepository.participationsSecondFetch = after

	auditRepository := newFakeAuditLogRepo()
	activityUserRepo := newFakeActivityUserRepo()
	uc := usecases.NewActivityUseCase(
		activityRepository,
		activityUserRepo,
		auditRepository,
		configs.Supabase{},
	)

	result, err := uc.BulkUpdateParticipationIsParticipatingByResidentIDs(activityModels.BulkUpdateParticipationIsParticipatingByResidentIDsRequest{
		ASID:            "as-1",
		ResidentIDs:     []string{" r1 ", "r2", "r1", "", "   "},
		IsParticipating: &isParticipating,
	}, "user-1")

	assert.NoError(t, err)
	assert.Len(t, result, 2)
	assert.Equal(t, 1, activityRepository.updateActivityCallCount)
	if assert.NotNil(t, activityRepository.updatedActivity) {
		assert.Equal(t, "staff-1", activityRepository.updatedActivity.StaffID)
	}
	assert.Equal(t, 1, activityRepository.updateParticipationsCallCount)
	assert.Equal(t, "as-1", activityRepository.capturedUpdateASID)
	assert.Equal(t, []string{"r1", "r2"}, activityRepository.capturedUpdateResidentIDs)
	assert.True(t, activityRepository.capturedUpdateIsParticipating)
	assert.Equal(t, 3, auditRepository.createAuditLogCallCount)
}

func TestUpdateActivityScheduleWithActivitySyncByID_Success_UpdateActivityAndSchedule(t *testing.T) {
	oldLocation := "Room A"
	oldDescription := "old description"
	newLocation := "Room C"
	newDescription := "new description"

	oldStartTime := time.Date(2026, 4, 23, 9, 0, 0, 0, time.UTC)
	oldEndTime := time.Date(2026, 4, 23, 10, 0, 0, 0, time.UTC)
	newDate := time.Date(2026, 4, 25, 0, 0, 0, 0, time.UTC)
	newStartTime := time.Date(2026, 4, 25, 14, 0, 0, 0, time.UTC)
	newEndTime := time.Date(2026, 4, 25, 15, 30, 0, 0, time.UTC)

	existingActivity := &entities.Activity{
		ID:           "act-1",
		StaffID:      "staff-1",
		ActivityName: "Morning Exercise",
		ActivityType: "Wellness",
		Location:     &oldLocation,
		Description:  &oldDescription,
	}
	existingSchedule := &entities.ActivitySchedule{
		ID:         "as-10",
		ActivityID: "act-1",
		Date:       time.Date(2026, 4, 23, 0, 0, 0, 0, time.UTC),
		StartTime:  oldStartTime,
		EndTime:    oldEndTime,
	}

	activityRepository := newFakeActivityRepo(existingActivity, nil)
	activityRepository.existingActivitySchedule = existingSchedule

	auditRepository := newFakeAuditLogRepo()
	activityUserRepo := newFakeActivityUserRepo()
	uc := usecases.NewActivityUseCase(
		activityRepository,
		activityUserRepo,
		auditRepository,
		configs.Supabase{},
	)

	newType := "Rehab"
	result, err := uc.UpdateActivityScheduleWithActivitySyncByID("as-10", activityModels.UpdateActivityScheduleWithActivitySyncRequest{
		ActivityType: &newType,
		Location:     &newLocation,
		Description:  &newDescription,
		Date:         &newDate,
		StartTime:    &newStartTime,
		EndTime:      &newEndTime,
	}, "user-1")

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "as-10", result.ID)
	assert.True(t, newDate.Equal(result.Date))
	assert.True(t, newStartTime.Equal(result.StartTime))
	assert.True(t, newEndTime.Equal(result.EndTime))

	assert.Equal(t, 2, activityRepository.updateActivityCallCount)
	if assert.NotNil(t, activityRepository.updatedActivity) {
		assert.Equal(t, "Rehab", activityRepository.updatedActivity.ActivityType)
		if assert.NotNil(t, activityRepository.updatedActivity.Location) {
			assert.Equal(t, "Room C", *activityRepository.updatedActivity.Location)
		}
		if assert.NotNil(t, activityRepository.updatedActivity.Description) {
			assert.Equal(t, "new description", *activityRepository.updatedActivity.Description)
		}
	}

	assert.Equal(t, 1, activityRepository.updateActivityScheduleCallCount)
	if assert.NotNil(t, activityRepository.updatedActivitySchedule) {
		assert.True(t, newDate.Equal(activityRepository.updatedActivitySchedule.Date))
		assert.True(t, newStartTime.Equal(activityRepository.updatedActivitySchedule.StartTime))
		assert.True(t, newEndTime.Equal(activityRepository.updatedActivitySchedule.EndTime))
	}

	assert.Equal(t, 3, auditRepository.createAuditLogCallCount)
}

func TestGetResidentsByScheduleIDCustom_Success_MapAndDedupeIntakeLabels(t *testing.T) {
	nickname := "Bobby"
	activityRepository := newFakeActivityRepo(nil, nil)
	activityRepository.residentsByScheduleResponse = []*entities.Participation{
		{
			ResidentID:      "r1",
			ASID:            "as-1",
			IsParticipating: true,
			Resident: entities.Resident{
				FirstName: "Bob",
				LastName:  "Lee",
				Nickname:  &nickname,
				Room: entities.Room{
					RoomNumber: "A-101",
					Floor:      2,
				},
				ResidentLabels: []entities.ResidentLabels{
					{IntakeLabel: entities.IntakeLabels{LabelName: "Diabetes"}},
					{IntakeLabel: entities.IntakeLabels{LabelName: "Diabetes"}},
					{IntakeLabel: entities.IntakeLabels{LabelName: " "}},
					{IntakeLabel: entities.IntakeLabels{LabelName: "Low Sodium"}},
				},
			},
		},
	}

	activityUserRepo := newFakeActivityUserRepo()
	uc := usecases.NewActivityUseCase(
		activityRepository,
		activityUserRepo,
		newFakeAuditLogRepo(),
		configs.Supabase{},
	)

	result, err := uc.GetResidentsByScheduleIDCustom("as-1", activityModels.ResidentsByScheduleQueryParams{})

	assert.NoError(t, err)
	if assert.NotNil(t, result) {
		assert.Len(t, result.Items, 1)
	}
	assert.Equal(t, 1, activityRepository.getResidentsByScheduleCallCount)
	assert.Equal(t, "as-1", activityRepository.capturedResidentsByScheduleASID)

	if assert.NotNil(t, result) && assert.Len(t, result.Items, 1) {
		item := result.Items[0]
		assert.Equal(t, "r1", item.ResidentID)
		assert.Equal(t, "Bob", item.FirstName)
		assert.Equal(t, "Lee", item.LastName)
		if assert.NotNil(t, item.Nickname) {
			assert.Equal(t, "Bobby", *item.Nickname)
		}
		assert.Equal(t, "A-101", item.RoomNumber)
		assert.Equal(t, int16(2), item.Floor)
		assert.True(t, item.IsParticipating)
		assert.Equal(t, []string{"Diabetes", "Low Sodium"}, item.IntakeLabels)
	}
}
