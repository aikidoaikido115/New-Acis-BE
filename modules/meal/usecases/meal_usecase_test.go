package usecases_test

import (
	"context"
	"encoding/json"
	"errors"
	"testing"

	auditRepositories "github.com/aikidoaikido115/New-Acis-BE/modules/audit_logs/repositories"
	emrModels "github.com/aikidoaikido115/New-Acis-BE/modules/emr/models"
	emrRepositories "github.com/aikidoaikido115/New-Acis-BE/modules/emr/repositories"
	"github.com/aikidoaikido115/New-Acis-BE/modules/entities"
	mealConstants "github.com/aikidoaikido115/New-Acis-BE/modules/meal/constants"
	mealModels "github.com/aikidoaikido115/New-Acis-BE/modules/meal/models"
	mealRepositories "github.com/aikidoaikido115/New-Acis-BE/modules/meal/repositories"
	mealUsecases "github.com/aikidoaikido115/New-Acis-BE/modules/meal/usecases"
	userConstants "github.com/aikidoaikido115/New-Acis-BE/modules/user/constants"
	userRepositories "github.com/aikidoaikido115/New-Acis-BE/modules/user/repositories"
	aiinfra "github.com/aikidoaikido115/New-Acis-BE/pkg/ai"
	"github.com/stretchr/testify/assert"
)

type fakeMealRepo struct {
	*mealRepositories.GormMealRepository

	menusByID               map[string]*entities.Menu
	createdMealPlanResponse *entities.MealPlan
	deletedMealPlans        []*entities.MealPlan
	mealHistoryItems        []mealModels.MealHistoryItem
	mealHistoryTotal        int64

	getMenuByIDErr     error
	deleteMealPlansErr error
	createMealPlanErr  error
	mealHistoryErr     error

	capturedCreatedMealPlan *entities.MealPlan
	capturedHistoryDate     *string
	capturedHistoryMealType *string
	capturedHistorySearch   *string
	capturedHistoryPage     int
	capturedHistoryPageSize int

	getMenuByIDCalls     int
	deleteMealPlansCalls int
	createMealPlanCalls  int
	getMealHistoryCalls  int
}

func newFakeMealRepo() *fakeMealRepo {
	return &fakeMealRepo{
		GormMealRepository: mealRepositories.NewGormMealRepository(nil),
		menusByID:          map[string]*entities.Menu{},
	}
}

func (f *fakeMealRepo) GetMenuByID(id string) (*entities.Menu, error) {
	f.getMenuByIDCalls++
	if f.getMenuByIDErr != nil {
		return nil, f.getMenuByIDErr
	}
	menu, ok := f.menusByID[id]
	if !ok {
		return nil, errors.New("menu not found")
	}
	return menu, nil
}

func (f *fakeMealRepo) DeleteMealPlansToday() ([]*entities.MealPlan, error) {
	f.deleteMealPlansCalls++
	if f.deleteMealPlansErr != nil {
		return nil, f.deleteMealPlansErr
	}
	return f.deletedMealPlans, nil
}

func (f *fakeMealRepo) CreateMealPlan(mealPlan *entities.MealPlan) (*entities.MealPlan, error) {
	f.createMealPlanCalls++
	if f.createMealPlanErr != nil {
		return nil, f.createMealPlanErr
	}
	copied := *mealPlan
	f.capturedCreatedMealPlan = &copied

	if f.createdMealPlanResponse != nil {
		return f.createdMealPlanResponse, nil
	}
	return mealPlan, nil
}

func (f *fakeMealRepo) GetMealHistory(date *string, mealType *string, search *string, page int, pageSize int) ([]mealModels.MealHistoryItem, int64, error) {
	f.getMealHistoryCalls++
	f.capturedHistoryDate = date
	f.capturedHistoryMealType = mealType
	f.capturedHistorySearch = search
	f.capturedHistoryPage = page
	f.capturedHistoryPageSize = pageSize

	if f.mealHistoryErr != nil {
		return nil, 0, f.mealHistoryErr
	}
	return f.mealHistoryItems, f.mealHistoryTotal, nil
}

type fakeMealUserRepo struct {
	*userRepositories.GormUserRepository

	user  *entities.User
	role  *entities.Role
	staff *entities.Staff

	getUserErr  error
	getRoleErr  error
	getStaffErr error

	getUserCalls  int
	getRoleCalls  int
	getStaffCalls int
}

func newFakeMealUserRepo(roleName string) *fakeMealUserRepo {
	return &fakeMealUserRepo{
		GormUserRepository: userRepositories.NewGormUserRepository(nil),
		user: &entities.User{
			ID:       "user-1",
			RoleID:   "role-1",
			Username: "janedoe",
		},
		role: &entities.Role{
			ID:   "role-1",
			Name: roleName,
		},
		staff: &entities.Staff{
			ID:     "staff-1",
			UserID: "user-1",
			User: entities.User{
				FirstName: "Jane",
				LastName:  "Doe",
				Username:  "janedoe",
			},
		},
	}
}

func (f *fakeMealUserRepo) GetUserByID(id string) (*entities.User, error) {
	f.getUserCalls++
	if f.getUserErr != nil {
		return nil, f.getUserErr
	}
	return f.user, nil
}

func (f *fakeMealUserRepo) GetRoleByID(roleID string) (*entities.Role, error) {
	f.getRoleCalls++
	if f.getRoleErr != nil {
		return nil, f.getRoleErr
	}
	return f.role, nil
}

func (f *fakeMealUserRepo) GetStaffByUserID(userID string) (*entities.Staff, error) {
	f.getStaffCalls++
	if f.getStaffErr != nil {
		return nil, f.getStaffErr
	}
	return f.staff, nil
}

type fakeMealEmrRepo struct {
	*emrRepositories.GormEmrRepository

	allergyStatsResponse emrModels.ResidentAllergyStatsDashboardResponse
	allergyStatsErr      error

	getResidentAllergyStatsCalls int
}

func newFakeMealEmrRepo() *fakeMealEmrRepo {
	return &fakeMealEmrRepo{GormEmrRepository: emrRepositories.NewGormEmrRepository(nil)}
}

func (f *fakeMealEmrRepo) GetResidentAllergyStatsDashboard() (emrModels.ResidentAllergyStatsDashboardResponse, error) {
	f.getResidentAllergyStatsCalls++
	if f.allergyStatsErr != nil {
		return emrModels.ResidentAllergyStatsDashboardResponse{}, f.allergyStatsErr
	}
	return f.allergyStatsResponse, nil
}

type fakeMealAuditRepo struct {
	*auditRepositories.GormAuditLogRepository

	createAuditLogCalls int
}

func newFakeMealAuditRepo() *fakeMealAuditRepo {
	return &fakeMealAuditRepo{GormAuditLogRepository: auditRepositories.NewGormAuditLogRepository(nil)}
}

func (f *fakeMealAuditRepo) CreateAuditLog(auditLog *entities.AuditLogs) (*entities.AuditLogs, error) {
	f.createAuditLogCalls++
	return auditLog, nil
}

type fakeAllergyAIClient struct {
	responses []aiinfra.CheckAllergyResponse
	err       error

	checkAllergyCalls int
	payloads          []aiinfra.CheckAllergyRequest
}

func (f *fakeAllergyAIClient) CheckAllergy(ctx context.Context, payload aiinfra.CheckAllergyRequest) ([]byte, error) {
	f.checkAllergyCalls++
	f.payloads = append(f.payloads, payload)

	if f.err != nil {
		return nil, f.err
	}

	index := f.checkAllergyCalls - 1
	if index < 0 || index >= len(f.responses) {
		index = len(f.responses) - 1
	}

	if index < 0 {
		return nil, errors.New("no ai response configured")
	}

	return json.Marshal(f.responses[index])
}

func newMealUsecase(roleName string) (*mealUsecases.MealUseCaseImpl, *fakeMealRepo, *fakeMealUserRepo, *fakeMealEmrRepo, *fakeMealAuditRepo, *fakeAllergyAIClient) {
	mealRepo := newFakeMealRepo()
	userRepo := newFakeMealUserRepo(roleName)
	emrRepo := newFakeMealEmrRepo()
	auditRepo := newFakeMealAuditRepo()
	aiClient := &fakeAllergyAIClient{}

	uc := mealUsecases.NewMealUseCase(
		mealRepo,
		emrRepo,
		auditRepo,
		userRepo,
		aiClient,
	).(*mealUsecases.MealUseCaseImpl)

	return uc, mealRepo, userRepo, emrRepo, auditRepo, aiClient
}

func ptrString(v string) *string {
	return &v
}

func ptrInt16(v int16) *int16 {
	return &v
}

func TestCreateMealPlan_Success_WithAI(t *testing.T) {
	uc, mealRepo, userRepo, emrRepo, auditRepo, aiClient := newMealUsecase(userConstants.RoleKitchenStaff)

	mealRepo.menusByID["menu-main"] = &entities.Menu{ID: "menu-main", MenuName: "Main Menu", Description: "egg, rice"}
	mealRepo.menusByID["menu-backup"] = &entities.Menu{ID: "menu-backup", MenuName: "Backup Menu", Description: "fish, soup"}
	mealRepo.deletedMealPlans = []*entities.MealPlan{{
		ID:           "old-plan-1",
		MenuID:       "menu-main",
		BackUpMenuID: ptrString("menu-backup"),
		MainAmount:   5,
		BackUpAmount: ptrInt16(2),
		MealType:     "breakfast",
	}}

	emrRepo.allergyStatsResponse = emrModels.ResidentAllergyStatsDashboardResponse{
		AllergyDetails: []emrModels.AllergyStatisticDashboardResponse{
			{AllergyID: "a-1", AllergyName: "Egg", ResidentCount: 3},
		},
	}

	aiClient.responses = []aiinfra.CheckAllergyResponse{
		{Status: mealConstants.AllergyCheckStatusSafe, Reason: "safe"},
		{Status: mealConstants.AllergyCheckStatusSafe, Reason: "safe"},
	}

	mealRepo.createdMealPlanResponse = &entities.MealPlan{
		ID:               "new-plan-1",
		MenuID:           "menu-main",
		BackUpMenuID:     ptrString("menu-backup"),
		MainAmount:       10,
		BackUpAmount:     ptrInt16(3),
		MealType:         "breakfast",
		CreatedByStaffID: "staff-1",
		StaffName:        "Jane Doe",
	}

	mealPlan := &entities.MealPlan{
		MenuID:       "menu-main",
		BackUpMenuID: ptrString("menu-backup"),
		MainAmount:   10,
		BackUpAmount: ptrInt16(3),
		MealType:     "Breakfast",
	}

	result, warning, err := uc.CreateMealPlan(mealPlan, "user-1", false)

	assert.NoError(t, err)
	assert.Nil(t, warning)
	if assert.NotNil(t, result) {
		assert.Equal(t, "new-plan-1", result.ID)
		assert.Equal(t, "staff-1", result.CreatedByStaffID)
		assert.Equal(t, "Jane Doe", result.StaffName)
	}

	assert.Equal(t, 1, userRepo.getUserCalls)
	assert.Equal(t, 1, userRepo.getRoleCalls)
	assert.Equal(t, 1, userRepo.getStaffCalls)
	assert.Equal(t, 2, aiClient.checkAllergyCalls)
	assert.Equal(t, 1, emrRepo.getResidentAllergyStatsCalls)
	assert.Equal(t, 1, mealRepo.deleteMealPlansCalls)
	assert.Equal(t, 1, mealRepo.createMealPlanCalls)
	assert.Equal(t, 2, auditRepo.createAuditLogCalls)

	if assert.NotNil(t, mealRepo.capturedCreatedMealPlan) {
		assert.Equal(t, "breakfast", mealRepo.capturedCreatedMealPlan.MealType)
		assert.Equal(t, "staff-1", mealRepo.capturedCreatedMealPlan.CreatedByStaffID)
		assert.Equal(t, "Jane Doe", mealRepo.capturedCreatedMealPlan.StaffName)
	}
}

func TestCreateMealPlan_Success_WithAIWarningWhenMainMenuNotSafe(t *testing.T) {
	uc, mealRepo, _, emrRepo, auditRepo, aiClient := newMealUsecase(userConstants.RoleKitchenStaff)

	mealRepo.menusByID["menu-main"] = &entities.Menu{ID: "menu-main", MenuName: "Main Menu", Description: "shrimp, rice"}
	mealRepo.menusByID["menu-backup"] = &entities.Menu{ID: "menu-backup", MenuName: "Backup Menu", Description: "chicken, soup"}

	emrRepo.allergyStatsResponse = emrModels.ResidentAllergyStatsDashboardResponse{
		AllergyDetails: []emrModels.AllergyStatisticDashboardResponse{
			{AllergyID: "a-2", AllergyName: "Shrimp", ResidentCount: 2},
		},
	}

	aiClient.responses = []aiinfra.CheckAllergyResponse{
		{Status: mealConstants.AllergyCheckStatusAllergyWarn, Reason: "contains shrimp"},
		{Status: mealConstants.AllergyCheckStatusSafe, Reason: "safe"},
	}

	mealPlan := &entities.MealPlan{
		MenuID:       "menu-main",
		BackUpMenuID: ptrString("menu-backup"),
		MainAmount:   8,
		BackUpAmount: ptrInt16(4),
		MealType:     "lunch",
	}

	result, warning, err := uc.CreateMealPlan(mealPlan, "user-1", false)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	if assert.NotNil(t, warning) {
		assert.False(t, warning.MainMenuPassed)
		assert.True(t, warning.BackupMenuPassed)
		if assert.NotNil(t, warning.MainMenuResult) {
			assert.Equal(t, mealConstants.AllergyCheckStatusAllergyWarn, warning.MainMenuResult.Status)
		}
		if assert.NotNil(t, warning.BackupMenuResult) {
			assert.Equal(t, mealConstants.AllergyCheckStatusSafe, warning.BackupMenuResult.Status)
		}
	}

	assert.Equal(t, 2, aiClient.checkAllergyCalls)
	assert.Equal(t, 1, mealRepo.createMealPlanCalls)
	assert.Equal(t, 1, auditRepo.createAuditLogCalls)
}

func TestGetMealHistory_Success_NormalizeFiltersAndPagination(t *testing.T) {
	uc, mealRepo, _, _, _, _ := newMealUsecase(userConstants.RoleKitchenStaff)

	backupMenuName := "Steamed Veg"
	backupAmount := int16(3)
	mealRepo.mealHistoryItems = []mealModels.MealHistoryItem{
		{
			Date:           "2026-04-30",
			MealType:       "breakfast",
			MenuName:       "Congee",
			MainAmount:     10,
			BackupMenuName: &backupMenuName,
			BackupAmount:   &backupAmount,
			StaffName:      "Jane Doe",
		},
	}
	mealRepo.mealHistoryTotal = 7

	date := "2026-04-30"
	mealType := " เช้า "
	search := "jane"
	page := 2
	pageSize := 3

	result, err := uc.GetMealHistory("user-1", mealModels.MealHistoryQueryParams{
		Date:     &date,
		MealType: &mealType,
		Search:   &search,
		Page:     &page,
		PageSize: &pageSize,
	})

	assert.NoError(t, err)
	if assert.NotNil(t, result) {
		assert.Len(t, result.Items, 1)
		assert.Equal(t, 2, result.Pagination.Page)
		assert.Equal(t, 3, result.Pagination.PageSize)
		assert.Equal(t, 7, result.Pagination.TotalItems)
		assert.Equal(t, 3, result.Pagination.TotalPages)
	}

	assert.Equal(t, 1, mealRepo.getMealHistoryCalls)
	if assert.NotNil(t, mealRepo.capturedHistoryMealType) {
		assert.Equal(t, "breakfast", *mealRepo.capturedHistoryMealType)
	}
	if assert.NotNil(t, mealRepo.capturedHistoryDate) {
		assert.Equal(t, "2026-04-30", *mealRepo.capturedHistoryDate)
	}
	if assert.NotNil(t, mealRepo.capturedHistorySearch) {
		assert.Equal(t, "jane", *mealRepo.capturedHistorySearch)
	}
	assert.Equal(t, 2, mealRepo.capturedHistoryPage)
	assert.Equal(t, 3, mealRepo.capturedHistoryPageSize)
}

func TestGetMealHistory_Error_InvalidMealType(t *testing.T) {
	uc, mealRepo, _, _, _, _ := newMealUsecase(userConstants.RoleKitchenStaff)

	mealType := "snack"
	result, err := uc.GetMealHistory("user-1", mealModels.MealHistoryQueryParams{
		MealType: &mealType,
	})

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "meal_type must be one of")
	assert.Nil(t, result)
	assert.Equal(t, 0, mealRepo.getMealHistoryCalls)
}
