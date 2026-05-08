package usecases_test

import (
	"errors"
	"testing"
	"time"

	auditRepositories "github.com/aikidoaikido115/New-Acis-BE/modules/audit_logs/repositories"
	"github.com/aikidoaikido115/New-Acis-BE/modules/entities"
	medicineModels "github.com/aikidoaikido115/New-Acis-BE/modules/medicine/models"
	medicineRepositories "github.com/aikidoaikido115/New-Acis-BE/modules/medicine/repositories"
	medicineUsecases "github.com/aikidoaikido115/New-Acis-BE/modules/medicine/usecases"
	userConstants "github.com/aikidoaikido115/New-Acis-BE/modules/user/constants"
	userRepositories "github.com/aikidoaikido115/New-Acis-BE/modules/user/repositories"
	"github.com/stretchr/testify/assert"
)

type fakeDrugPlanUserRepo struct {
	*userRepositories.GormUserRepository

	actingUser *entities.User
	actingRole *entities.Role

	staffSearchResult []*entities.User
	staffByUserID     map[string]*entities.Staff

	getUserCalls          int
	getRoleCalls          int
	getUsersByNameCalls   int
	getStaffByUserIDCalls int
}

func newFakeDrugPlanUserRepo() *fakeDrugPlanUserRepo {
	medicalUser := &entities.User{ID: "user-1", RoleID: "role-medical"}
	return &fakeDrugPlanUserRepo{
		GormUserRepository: userRepositories.NewGormUserRepository(nil),
		actingUser:         medicalUser,
		actingRole:         &entities.Role{ID: "role-medical", Name: userConstants.RoleMedicalStaff},
		staffSearchResult: []*entities.User{
			{ID: "staff-user-1", Role: entities.Role{Name: userConstants.RoleMedicalStaff}},
		},
		staffByUserID: map[string]*entities.Staff{
			"staff-user-1": {ID: "staff-1", UserID: "staff-user-1"},
		},
	}
}

func (f *fakeDrugPlanUserRepo) GetUserByID(id string) (*entities.User, error) {
	f.getUserCalls++
	return f.actingUser, nil
}

func (f *fakeDrugPlanUserRepo) GetRoleByID(roleID string) (*entities.Role, error) {
	f.getRoleCalls++
	return f.actingRole, nil
}

func (f *fakeDrugPlanUserRepo) GetUsersByFirstAndLastName(firstName string, lastName string) ([]*entities.User, error) {
	f.getUsersByNameCalls++
	return f.staffSearchResult, nil
}

func (f *fakeDrugPlanUserRepo) GetStaffByUserID(userID string) (*entities.Staff, error) {
	f.getStaffByUserIDCalls++
	staff, ok := f.staffByUserID[userID]
	if !ok {
		return nil, errors.New("staff not found")
	}
	return staff, nil
}

type fakeDrugPlanRepo struct {
	*medicineRepositories.GormDrugRepository

	residentExists bool
	plansByID      map[string]*entities.DrugPlan
	plansToday     []*entities.DrugPlan

	updateDrugPlanCalls              int
	getDrugPlanByIDCalls             int
	residentExistsByIDCalls          int
	getDrugPlansByResidentTodayCalls int

	capturedUpdatedPlans []*entities.DrugPlan
}

func newFakeDrugPlanRepo() *fakeDrugPlanRepo {
	return &fakeDrugPlanRepo{
		GormDrugRepository: medicineRepositories.NewGormDrugRepository(nil),
		residentExists:     true,
		plansByID:          map[string]*entities.DrugPlan{},
	}
}

func (f *fakeDrugPlanRepo) GetDrugPlanByID(id string) (*entities.DrugPlan, error) {
	f.getDrugPlanByIDCalls++
	plan, ok := f.plansByID[id]
	if !ok {
		return nil, errors.New("not found")
	}
	return plan, nil
}

func (f *fakeDrugPlanRepo) UpdateDrugPlan(drugPlan *entities.DrugPlan) (*entities.DrugPlan, error) {
	f.updateDrugPlanCalls++
	copied := *drugPlan
	f.capturedUpdatedPlans = append(f.capturedUpdatedPlans, &copied)
	f.plansByID[drugPlan.ID] = &copied
	return &copied, nil
}

func (f *fakeDrugPlanRepo) ResidentExistsByID(id string) (bool, error) {
	f.residentExistsByIDCalls++
	return f.residentExists, nil
}

func (f *fakeDrugPlanRepo) GetDrugPlansByResidentIDToday(residentID string) ([]*entities.DrugPlan, error) {
	f.getDrugPlansByResidentTodayCalls++
	return f.plansToday, nil
}

func (f *fakeDrugPlanRepo) GetExpiredAsNeededPersonalDrugs(date time.Time, residentID *string) ([]*entities.PersonalDrug, error) {
	return []*entities.PersonalDrug{}, nil
}

func (f *fakeDrugPlanRepo) GetActivePersonalDrugsForDate(date time.Time, residentID *string) ([]*entities.PersonalDrug, error) {
	return []*entities.PersonalDrug{}, nil
}

type fakeDrugPlanAuditRepo struct {
	*auditRepositories.GormAuditLogRepository
	createAuditLogCalls int
}

func newFakeDrugPlanAuditRepo() *fakeDrugPlanAuditRepo {
	return &fakeDrugPlanAuditRepo{GormAuditLogRepository: auditRepositories.NewGormAuditLogRepository(nil)}
}

func (f *fakeDrugPlanAuditRepo) CreateAuditLog(auditLog *entities.AuditLogs) (*entities.AuditLogs, error) {
	f.createAuditLogCalls++
	return auditLog, nil
}

func newDrugPlanUsecase() (*medicineUsecases.DrugUseCaseImpl, *fakeDrugPlanRepo, *fakeDrugPlanUserRepo, *fakeDrugPlanAuditRepo) {
	drugRepo := newFakeDrugPlanRepo()
	userRepo := newFakeDrugPlanUserRepo()
	auditRepo := newFakeDrugPlanAuditRepo()

	uc := medicineUsecases.NewDrugUseCase(drugRepo, auditRepo, userRepo)
	return uc, drugRepo, userRepo, auditRepo
}

func TestTakeDrugPlanByID_Success(t *testing.T) {
	uc, drugRepo, userRepo, auditRepo := newDrugPlanUsecase()
	now := time.Now()

	drugRepo.plansByID["dp-1"] = &entities.DrugPlan{
		ID:             "dp-1",
		CreatedAt:      now,
		IsTaken:        false,
		IsOmitted:      func() *bool { v := false; return &v }(),
		GivenByStaffID: "old-staff",
	}

	note := "  given after meal  "
	result, err := uc.TakeDrugPlanByID("dp-1", medicineModels.TakeDrugPlanByIDRequest{
		StaffFirstName: "John",
		StaffLastName:  "Doe",
		Note:           &note,
	}, "user-1")

	assert.NoError(t, err)
	if assert.NotNil(t, result) {
		assert.True(t, result.IsTaken)
		if assert.NotNil(t, result.IsOmitted) {
			assert.False(t, *result.IsOmitted)
		}
		assert.Nil(t, result.OmittedReason)
		if assert.NotNil(t, result.Notes) {
			assert.Equal(t, "given after meal", *result.Notes)
		}
		assert.Equal(t, "staff-1", result.GivenByStaffID)
		assert.NotNil(t, result.TakenAt)
	}
	assert.Equal(t, 1, userRepo.getUserCalls)
	assert.Equal(t, 1, userRepo.getRoleCalls)
	assert.Equal(t, 1, userRepo.getUsersByNameCalls)
	assert.Equal(t, 1, userRepo.getStaffByUserIDCalls)
	assert.Equal(t, 1, drugRepo.getDrugPlanByIDCalls)
	assert.Equal(t, 1, drugRepo.updateDrugPlanCalls)
	assert.Equal(t, 1, auditRepo.createAuditLogCalls)
}

func TestOmitDrugPlanByID_Error_WhenOmittedReasonMissing(t *testing.T) {
	uc, drugRepo, _, auditRepo := newDrugPlanUsecase()
	now := time.Now()
	drugRepo.plansByID["dp-2"] = &entities.DrugPlan{ID: "dp-2", CreatedAt: now}

	_, err := uc.OmitDrugPlanByID("dp-2", medicineModels.OmitDrugPlanByIDRequest{
		StaffFirstName: "John",
		StaffLastName:  "Doe",
		OmittedReason:  "   ",
	}, "user-1")

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "omitted_reason is required")
	assert.Equal(t, 0, drugRepo.getDrugPlanByIDCalls)
	assert.Equal(t, 0, drugRepo.updateDrugPlanCalls)
	assert.Equal(t, 0, auditRepo.createAuditLogCalls)
}

func TestOmitDrugPlanByID_Error_WhenPlanNotToday(t *testing.T) {
	uc, drugRepo, _, auditRepo := newDrugPlanUsecase()
	yesterday := time.Now().Add(-24 * time.Hour)
	drugRepo.plansByID["dp-3"] = &entities.DrugPlan{
		ID:        "dp-3",
		CreatedAt: yesterday,
	}

	_, err := uc.OmitDrugPlanByID("dp-3", medicineModels.OmitDrugPlanByIDRequest{
		StaffFirstName: "John",
		StaffLastName:  "Doe",
		OmittedReason:  "resident refused",
	}, "user-1")

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "today records only")
	assert.Equal(t, 1, drugRepo.getDrugPlanByIDCalls)
	assert.Equal(t, 0, drugRepo.updateDrugPlanCalls)
	assert.Equal(t, 0, auditRepo.createAuditLogCalls)
}

func TestTakeDrugPlansByResidentIDToday_Success(t *testing.T) {
	uc, drugRepo, _, auditRepo := newDrugPlanUsecase()
	now := time.Now()

	drugRepo.plansToday = []*entities.DrugPlan{
		{ID: "dp-a"},
		{ID: "dp-b"},
	}
	drugRepo.plansByID["dp-a"] = &entities.DrugPlan{ID: "dp-a", CreatedAt: now, IsTaken: false, IsOmitted: func() *bool { v := false; return &v }()}
	drugRepo.plansByID["dp-b"] = &entities.DrugPlan{ID: "dp-b", CreatedAt: now, IsTaken: false, IsOmitted: func() *bool { v := false; return &v }()}

	note := "done"
	result, err := uc.TakeDrugPlansByResidentIDToday("resident-1", medicineModels.TakeDrugPlansByResidentRequest{
		StaffFirstName: "John",
		StaffLastName:  "Doe",
		Note:           &note,
	}, "user-1")

	assert.NoError(t, err)
	assert.Len(t, result, 2)
	assert.Equal(t, 1, drugRepo.residentExistsByIDCalls)
	assert.Equal(t, 1, drugRepo.getDrugPlansByResidentTodayCalls)
	assert.Equal(t, 2, drugRepo.getDrugPlanByIDCalls)
	assert.Equal(t, 2, drugRepo.updateDrugPlanCalls)
	assert.Equal(t, 2, auditRepo.createAuditLogCalls)
	if assert.Len(t, drugRepo.capturedUpdatedPlans, 2) {
		assert.True(t, drugRepo.capturedUpdatedPlans[0].IsTaken)
		assert.True(t, drugRepo.capturedUpdatedPlans[1].IsTaken)
	}
}

func TestOmitDrugPlansByResidentIDToday_Success(t *testing.T) {
	uc, drugRepo, _, auditRepo := newDrugPlanUsecase()
	now := time.Now()

	drugRepo.plansToday = []*entities.DrugPlan{{ID: "dp-c"}}
	drugRepo.plansByID["dp-c"] = &entities.DrugPlan{ID: "dp-c", CreatedAt: now, IsTaken: false, IsOmitted: func() *bool { v := false; return &v }()}

	note := "nausea"
	result, err := uc.OmitDrugPlansByResidentIDToday("resident-1", medicineModels.OmitDrugPlansByResidentRequest{
		StaffFirstName: "John",
		StaffLastName:  "Doe",
		OmittedReason:  "resident refused",
		Note:           &note,
	}, "user-1")

	assert.NoError(t, err)
	assert.Len(t, result, 1)
	if assert.NotNil(t, result[0]) {
		assert.False(t, result[0].IsTaken)
		if assert.NotNil(t, result[0].IsOmitted) {
			assert.True(t, *result[0].IsOmitted)
		}
		if assert.NotNil(t, result[0].OmittedReason) {
			assert.Equal(t, "resident refused", *result[0].OmittedReason)
		}
	}
	assert.Equal(t, 1, drugRepo.updateDrugPlanCalls)
	assert.Equal(t, 1, auditRepo.createAuditLogCalls)
}

func TestTakeDrugPlansByResidentIDToday_Error_WhenResidentNotFound(t *testing.T) {
	uc, drugRepo, _, auditRepo := newDrugPlanUsecase()
	drugRepo.residentExists = false

	_, err := uc.TakeDrugPlansByResidentIDToday("resident-x", medicineModels.TakeDrugPlansByResidentRequest{
		StaffFirstName: "John",
		StaffLastName:  "Doe",
	}, "user-1")

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "resident not found")
	assert.Equal(t, 1, drugRepo.residentExistsByIDCalls)
	assert.Equal(t, 0, drugRepo.getDrugPlansByResidentTodayCalls)
	assert.Equal(t, 0, drugRepo.updateDrugPlanCalls)
	assert.Equal(t, 0, auditRepo.createAuditLogCalls)
}
