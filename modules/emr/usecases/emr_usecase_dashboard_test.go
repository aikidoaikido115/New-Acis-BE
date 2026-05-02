package usecases_test

import (
	"errors"
	"testing"

	"github.com/aikidoaikido115/New-Acis-BE/configs"
	auditRepositories "github.com/aikidoaikido115/New-Acis-BE/modules/audit_logs/repositories"
	emrModels "github.com/aikidoaikido115/New-Acis-BE/modules/emr/models"
	emrRepositories "github.com/aikidoaikido115/New-Acis-BE/modules/emr/repositories"
	emrUsecases "github.com/aikidoaikido115/New-Acis-BE/modules/emr/usecases"
	"github.com/aikidoaikido115/New-Acis-BE/modules/entities"
	medicineModels "github.com/aikidoaikido115/New-Acis-BE/modules/medicine/models"
	medicineUsecases "github.com/aikidoaikido115/New-Acis-BE/modules/medicine/usecases"
	userConstants "github.com/aikidoaikido115/New-Acis-BE/modules/user/constants"
	userRepositories "github.com/aikidoaikido115/New-Acis-BE/modules/user/repositories"
	"github.com/stretchr/testify/assert"
)

type fakeDashboardUserRepo struct {
	*userRepositories.GormUserRepository
	user         *entities.User
	role         *entities.Role
	getUserCalls int
	getRoleCalls int
	getUserErr   error
	getRoleErr   error
}

func newFakeDashboardUserRepo(roleName string) *fakeDashboardUserRepo {
	return &fakeDashboardUserRepo{
		GormUserRepository: userRepositories.NewGormUserRepository(nil),
		user: &entities.User{
			ID:     "user-1",
			RoleID: "role-1",
		},
		role: &entities.Role{
			ID:   "role-1",
			Name: roleName,
		},
	}
}

func (f *fakeDashboardUserRepo) GetUserByID(id string) (*entities.User, error) {
	f.getUserCalls++
	if f.getUserErr != nil {
		return nil, f.getUserErr
	}
	return f.user, nil
}

func (f *fakeDashboardUserRepo) GetRoleByID(roleID string) (*entities.Role, error) {
	f.getRoleCalls++
	if f.getRoleErr != nil {
		return nil, f.getRoleErr
	}
	return f.role, nil
}

type fakeDashboardEmrRepo struct {
	*emrRepositories.GormEmrRepository
	numberOfResidentsResponse   emrModels.NumberOfResidentsDashboardResponse
	residentGenderResponse      emrModels.ResidentGenderStatsDashboardResponse
	residentAllergyResponse     emrModels.ResidentAllergyStatsDashboardResponse
	residentDrugAllergyResponse emrModels.ResidentDrugAllergyStatsDashboardResponse
	vitalSignsTodayLatest       []*entities.VitalSign
	vitalSignsTodayAll          []*entities.VitalSign
	numberOfResidentsErr        error
	residentGenderErr           error
	residentAllergyErr          error
	residentDrugAllergyErr      error
	vitalSignsTodayLatestErr    error
	vitalSignsTodayAllErr       error
	getNumberOfResidentsCalls   int
	getResidentGenderCalls      int
	getResidentAllergyCalls     int
	getResidentDrugAllergyCalls int
	getVitalSignsTodayCalls     []bool
}

func newFakeDashboardEmrRepo() *fakeDashboardEmrRepo {
	return &fakeDashboardEmrRepo{GormEmrRepository: emrRepositories.NewGormEmrRepository(nil)}
}

func (f *fakeDashboardEmrRepo) GetNumberOfResidentsDashboard() (emrModels.NumberOfResidentsDashboardResponse, error) {
	f.getNumberOfResidentsCalls++
	return f.numberOfResidentsResponse, f.numberOfResidentsErr
}

func (f *fakeDashboardEmrRepo) GetNumberOfResidentGender() (emrModels.ResidentGenderStatsDashboardResponse, error) {
	f.getResidentGenderCalls++
	return f.residentGenderResponse, f.residentGenderErr
}

func (f *fakeDashboardEmrRepo) GetResidentAllergyStatsDashboard() (emrModels.ResidentAllergyStatsDashboardResponse, error) {
	f.getResidentAllergyCalls++
	return f.residentAllergyResponse, f.residentAllergyErr
}

func (f *fakeDashboardEmrRepo) GetResidentDrugAllergyStatsDashboard() (emrModels.ResidentDrugAllergyStatsDashboardResponse, error) {
	f.getResidentDrugAllergyCalls++
	return f.residentDrugAllergyResponse, f.residentDrugAllergyErr
}

func (f *fakeDashboardEmrRepo) GetVitalSignsToday(isLatest bool) ([]*entities.VitalSign, error) {
	f.getVitalSignsTodayCalls = append(f.getVitalSignsTodayCalls, isLatest)
	if isLatest {
		return f.vitalSignsTodayLatest, f.vitalSignsTodayLatestErr
	}
	return f.vitalSignsTodayAll, f.vitalSignsTodayAllErr
}

type fakeDashboardAuditRepo struct {
	*auditRepositories.GormAuditLogRepository
}

func newFakeDashboardAuditRepo() *fakeDashboardAuditRepo {
	return &fakeDashboardAuditRepo{GormAuditLogRepository: auditRepositories.NewGormAuditLogRepository(nil)}
}

func (f *fakeDashboardAuditRepo) CreateAuditLog(auditLog *entities.AuditLogs) (*entities.AuditLogs, error) {
	return auditLog, nil
}

type fakeDrugDashboardUsecase struct {
	*medicineUsecases.DrugUseCaseImpl
	response []*medicineModels.DrugPlanTimeOfDaySummary
	err      error
	calls    int
}

func newFakeDrugDashboardUsecase() *fakeDrugDashboardUsecase {
	return &fakeDrugDashboardUsecase{DrugUseCaseImpl: &medicineUsecases.DrugUseCaseImpl{}}
}

func (f *fakeDrugDashboardUsecase) GetDrugPlansTodayTimeOfDayResidentSummary(userID string) ([]*medicineModels.DrugPlanTimeOfDaySummary, error) {
	f.calls++
	return f.response, f.err
}

func newDashboardUsecase(roleName string) (*emrUsecases.EmrUseCaseImpl, *fakeDashboardUserRepo, *fakeDashboardEmrRepo, *fakeDrugDashboardUsecase) {
	userRepo := newFakeDashboardUserRepo(roleName)
	auditRepo := newFakeDashboardAuditRepo()
	emrRepo := newFakeDashboardEmrRepo()
	drugUsecase := newFakeDrugDashboardUsecase()

	uc := emrUsecases.NewEmrUseCase(
		emrRepo,
		auditRepo,
		userRepo,
		drugUsecase,
		configs.Supabase{},
	).(*emrUsecases.EmrUseCaseImpl)

	return uc, userRepo, emrRepo, drugUsecase
}

func ptrFloat64(v float64) *float64 {
	return &v
}

func TestGetNumberOfResidentsDashboard_Success(t *testing.T) {
	uc, userRepo, emrRepo, _ := newDashboardUsecase(userConstants.RoleMedicalStaff)
	emrRepo.numberOfResidentsResponse = emrModels.NumberOfResidentsDashboardResponse{
		TotalResidents:         12,
		IndependentResidents:   5,
		PartialAssistResidents: 4,
		BedriddenResidents:     3,
	}

	result, err := uc.GetNumberOfResidentsDashboard("user-1")

	assert.NoError(t, err)
	assert.Equal(t, emrModels.NumberOfResidentsDashboardResponse{
		TotalResidents:         12,
		IndependentResidents:   5,
		PartialAssistResidents: 4,
		BedriddenResidents:     3,
	}, result)
	assert.Equal(t, 1, userRepo.getUserCalls)
	assert.Equal(t, 1, userRepo.getRoleCalls)
	assert.Equal(t, 1, emrRepo.getNumberOfResidentsCalls)
}

func TestGetNumberOfResidentsDashboard_KitchenStaff_Success(t *testing.T) {
    uc, userRepo, emrRepo, _ := newDashboardUsecase(userConstants.RoleKitchenStaff)

    result, err := uc.GetNumberOfResidentsDashboard("user-1")

    assert.NoError(t, err)
    assert.Equal(t, 1, userRepo.getUserCalls)
    assert.Equal(t, 1, userRepo.getRoleCalls)
    assert.Equal(t, 1, emrRepo.getNumberOfResidentsCalls)
    
    assert.NotNil(t, result)
}

func TestGetNumberOfResidentsDashboard_RepoError(t *testing.T) {
	uc, _, emrRepo, _ := newDashboardUsecase(userConstants.RoleMedicalStaff)
	emrRepo.numberOfResidentsErr = errors.New("db down")

	result, err := uc.GetNumberOfResidentsDashboard("user-1")

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to get dashboard data: db down")
	assert.Equal(t, emrModels.NumberOfResidentsDashboardResponse{}, result)
}

func TestGetResidentGenderStatsDashboard_Success(t *testing.T) {
	uc, _, emrRepo, _ := newDashboardUsecase(userConstants.RoleMedicalStaff)
	emrRepo.residentGenderResponse = emrModels.ResidentGenderStatsDashboardResponse{
		SumOfMale:      3,
		SumOfFemale:    1,
		TotalResidents: 4,
	}

	result, err := uc.GetResidentGenderStatsDashboard("user-1")

	assert.NoError(t, err)
	assert.Equal(t, int16(3), result.SumOfMale)
	assert.Equal(t, int16(1), result.SumOfFemale)
	assert.Equal(t, int16(4), result.TotalResidents)
	assert.InDelta(t, 75.0, result.MalePercentage, 0.001)
	assert.InDelta(t, 25.0, result.FemalePercentage, 0.001)
}

func TestGetResidentGenderStatsDashboard_ZeroTotalResidents(t *testing.T) {
	uc, _, emrRepo, _ := newDashboardUsecase(userConstants.RoleMedicalStaff)
	emrRepo.residentGenderResponse = emrModels.ResidentGenderStatsDashboardResponse{}

	result, err := uc.GetResidentGenderStatsDashboard("user-1")

	assert.NoError(t, err)
	assert.Equal(t, float32(0), result.MalePercentage)
	assert.Equal(t, float32(0), result.FemalePercentage)
}

func TestGetResidentGenderStatsDashboard_RepoError(t *testing.T) {
	uc, _, emrRepo, _ := newDashboardUsecase(userConstants.RoleMedicalStaff)
	emrRepo.residentGenderErr = errors.New("query failed")

	result, err := uc.GetResidentGenderStatsDashboard("user-1")

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to get resident gender stats: query failed")
	assert.Equal(t, emrModels.ResidentGenderStatsDashboardResponse{}, result)
}

func TestGetVitalSignStatsDashboard_Success(t *testing.T) {
	uc, _, emrRepo, _ := newDashboardUsecase(userConstants.RoleMedicalStaff)

	normalTemp := ptrFloat64(36.5)
	abnormalTemp := ptrFloat64(38.2)

	emrRepo.vitalSignsTodayLatest = []*entities.VitalSign{
		{ResidentID: "resident-1", Temperature: normalTemp},
		{ResidentID: "resident-2", Temperature: abnormalTemp},
		nil,
		{ResidentID: "", Temperature: normalTemp},
	}
	emrRepo.vitalSignsTodayAll = []*entities.VitalSign{
		{ResidentID: "resident-1", Temperature: normalTemp},
		{ResidentID: "resident-1", Temperature: abnormalTemp},
		{ResidentID: "resident-2", Temperature: normalTemp},
		{ResidentID: "resident-3", Temperature: abnormalTemp},
		{ResidentID: "resident-3", Temperature: normalTemp},
	}

	result, err := uc.GetVitalSignStatsDashboard("user-1")

	assert.NoError(t, err)
	assert.Equal(t, int64(1), result.CurrentNormalResidents)
	assert.Equal(t, int64(1), result.CurrentAbnormalResidents)
	assert.Equal(t, int64(2), result.CurrentTotalResidents)
	assert.Equal(t, int64(2), result.HadAbnormalTodayResidents)
	assert.Equal(t, int64(1), result.HadNormalOnlyTodayResidents)
	assert.Equal(t, int64(3), result.HadTotalResidents)
	assert.Equal(t, []bool{true, false}, emrRepo.getVitalSignsTodayCalls)
}

func TestGetVitalSignStatsDashboard_RepoError(t *testing.T) {
	uc, _, emrRepo, _ := newDashboardUsecase(userConstants.RoleMedicalStaff)
	emrRepo.vitalSignsTodayLatestErr = errors.New("latest fetch failed")

	result, err := uc.GetVitalSignStatsDashboard("user-1")

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to get vital signs summary: latest fetch failed")
	assert.Equal(t, emrModels.VitalSignDashboardSummary{}, result)
	assert.Equal(t, []bool{true}, emrRepo.getVitalSignsTodayCalls)
}

func TestGetDrugPlanTimeOfDayStatsDashboard_Success(t *testing.T) {
	uc, _, _, drugUsecase := newDashboardUsecase(userConstants.RoleMedicalStaff)
	drugUsecase.response = []*medicineModels.DrugPlanTimeOfDaySummary{
		nil,
		{TimeOfDay: "เช้า", TotalResidents: 4, GivenResidents: 3, WaitingResidents: 1, Status: "partial"},
		{TimeOfDay: "เย็น", TotalResidents: 2, GivenResidents: 2, WaitingResidents: 0, Status: "full"},
	}

	result, err := uc.GetDrugPlanTimeOfDayStatsDashboard("user-1")

	assert.NoError(t, err)
	assert.Equal(t, 1, drugUsecase.calls)
	assert.Equal(t, []emrModels.DrugPlanTimeOfDayDashboardSummary{
		{TimeOfDay: "เช้า", TotalResidents: 4, GivenResidents: 3, WaitingResidents: 1, Status: "partial"},
		{TimeOfDay: "เย็น", TotalResidents: 2, GivenResidents: 2, WaitingResidents: 0, Status: "full"},
	}, result)
}

func TestGetDrugPlanTimeOfDayStatsDashboard_RepoError(t *testing.T) {
	uc, _, _, drugUsecase := newDashboardUsecase(userConstants.RoleMedicalStaff)
	drugUsecase.err = errors.New("drug summary failed")

	result, err := uc.GetDrugPlanTimeOfDayStatsDashboard("user-1")

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "drug summary failed")
	assert.Nil(t, result)
	assert.Equal(t, 1, drugUsecase.calls)
}

func TestGetResidentAllergyStatsDashboard_Success(t *testing.T) {
	uc, _, emrRepo, _ := newDashboardUsecase(userConstants.RoleMedicalStaff)
	emrRepo.residentAllergyResponse = emrModels.ResidentAllergyStatsDashboardResponse{
		TotalNotAllergic: 7,
		TotalAllergic:    5,
		AllergyDetails: []emrModels.AllergyStatisticDashboardResponse{
			{AllergyID: "a-1", AllergyName: "Peanut", ResidentCount: 3},
		},
	}

	result, err := uc.GetResidentAllergyStatsDashboard("user-1")

	assert.NoError(t, err)
	assert.Equal(t, emrRepo.residentAllergyResponse, result)
	assert.Equal(t, 1, emrRepo.getResidentAllergyCalls)
}

func TestGetResidentAllergyStatsDashboard_RepoError(t *testing.T) {
	uc, _, emrRepo, _ := newDashboardUsecase(userConstants.RoleMedicalStaff)
	emrRepo.residentAllergyErr = errors.New("allergy query failed")

	result, err := uc.GetResidentAllergyStatsDashboard("user-1")

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to get resident allergy stats: allergy query failed")
	assert.Equal(t, emrModels.ResidentAllergyStatsDashboardResponse{}, result)
}

func TestGetResidentDrugAllergyStatsDashboard_Success(t *testing.T) {
	uc, _, emrRepo, _ := newDashboardUsecase(userConstants.RoleMedicalStaff)
	emrRepo.residentDrugAllergyResponse = emrModels.ResidentDrugAllergyStatsDashboardResponse{
		TotalNotDrugAllergic: 8,
		TotalDrugAllergic:    2,
		DrugAllergyDetails: []emrModels.DrugAllergyStatisticDashboardResponse{
			{DrugAllergyID: "d-1", AllergyName: "Penicillin", Count: 2},
		},
	}

	result, err := uc.GetResidentDrugAllergyStatsDashboard("user-1")

	assert.NoError(t, err)
	assert.Equal(t, emrRepo.residentDrugAllergyResponse, result)
	assert.Equal(t, 1, emrRepo.getResidentDrugAllergyCalls)
}

func TestGetResidentDrugAllergyStatsDashboard_RepoError(t *testing.T) {
	uc, _, emrRepo, _ := newDashboardUsecase(userConstants.RoleMedicalStaff)
	emrRepo.residentDrugAllergyErr = errors.New("drug allergy query failed")

	result, err := uc.GetResidentDrugAllergyStatsDashboard("user-1")

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to get resident drug allergy stats: drug allergy query failed")
	assert.Equal(t, emrModels.ResidentDrugAllergyStatsDashboardResponse{}, result)
}
