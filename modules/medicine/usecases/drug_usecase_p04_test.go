package usecases

import (
	"errors"
	"testing"
	"time"

	audit_repository "github.com/aikidoaikido115/New-Acis-BE/modules/audit_logs/repositories"
	"github.com/aikidoaikido115/New-Acis-BE/modules/entities"
	medicine_repository "github.com/aikidoaikido115/New-Acis-BE/modules/medicine/repositories"
	user_constants "github.com/aikidoaikido115/New-Acis-BE/modules/user/constants"
	user_repository "github.com/aikidoaikido115/New-Acis-BE/modules/user/repositories"
	"github.com/stretchr/testify/assert"
)

type fakeDrugPlanP04UserRepo struct {
	userByID   map[string]*entities.User
	roleByID   map[string]*entities.Role
	getUserErr error
	getRoleErr error
}

func (f *fakeDrugPlanP04UserRepo) CreateUser(user *entities.User) (*entities.User, error) {
	return nil, errors.New("not used")
}
func (f *fakeDrugPlanP04UserRepo) CreateStaff(user *entities.User, staff *entities.Staff) (*entities.Staff, error) {
	return nil, errors.New("not used")
}
func (f *fakeDrugPlanP04UserRepo) CreateStaffFile(staffFile *entities.StaffsFiles) (*entities.StaffsFiles, error) {
	return nil, errors.New("not used")
}
func (f *fakeDrugPlanP04UserRepo) GetUserByEmail(email string) (*entities.User, error) {
	return nil, errors.New("not used")
}
func (f *fakeDrugPlanP04UserRepo) GetUserByID(id string) (*entities.User, error) {
	if f.getUserErr != nil {
		return nil, f.getUserErr
	}
	if user, ok := f.userByID[id]; ok {
		return user, nil
	}
	return nil, errors.New("user not found")
}
func (f *fakeDrugPlanP04UserRepo) GetUsersByFirstAndLastName(firstName string, lastName string) ([]*entities.User, error) {
	return nil, errors.New("not used")
}
func (f *fakeDrugPlanP04UserRepo) GetStaffByID(id string) (*entities.Staff, error) {
	return nil, errors.New("not used")
}
func (f *fakeDrugPlanP04UserRepo) GetStaffByUserID(userID string) (*entities.Staff, error) {
	return nil, errors.New("not used")
}
func (f *fakeDrugPlanP04UserRepo) GetStaffFileByID(id string) (*entities.StaffsFiles, error) {
	return nil, errors.New("not used")
}
func (f *fakeDrugPlanP04UserRepo) GetUserByUsername(username string) (*entities.User, error) {
	return nil, errors.New("not used")
}
func (f *fakeDrugPlanP04UserRepo) GetRoleByName(roleName string) (*entities.Role, error) {
	return nil, errors.New("not used")
}
func (f *fakeDrugPlanP04UserRepo) GetRoleByID(roleID string) (*entities.Role, error) {
	if f.getRoleErr != nil {
		return nil, f.getRoleErr
	}
	if role, ok := f.roleByID[roleID]; ok {
		return role, nil
	}
	return nil, errors.New("role not found")
}
func (f *fakeDrugPlanP04UserRepo) UsernameExists(username string) (bool, error) {
	return false, errors.New("not used")
}
func (f *fakeDrugPlanP04UserRepo) EmailExists(email string) (bool, error) {
	return false, errors.New("not used")
}
func (f *fakeDrugPlanP04UserRepo) GetAllUsers() ([]*entities.User, error) {
	return nil, errors.New("not used")
}
func (f *fakeDrugPlanP04UserRepo) UpdateUserByID(user *entities.User) error {
	return errors.New("not used")
}
func (f *fakeDrugPlanP04UserRepo) UpdateUserApprovalByID(userID string, isApprove bool) error {
	return errors.New("not used")
}
func (f *fakeDrugPlanP04UserRepo) DeleteStaffAndUserByStaffID(staffID string) error {
	return errors.New("not used")
}
func (f *fakeDrugPlanP04UserRepo) DeleteRelativeAndUserByUserID(userID string) error {
	return errors.New("not used")
}
func (f *fakeDrugPlanP04UserRepo) DeleteUserByID(userID string) error {
	return errors.New("not used")
}
func (f *fakeDrugPlanP04UserRepo) GetRelativeUserByUserID(userID string) (*user_repository.AdminRelativeUser, error) {
	return nil, errors.New("not used")
}
func (f *fakeDrugPlanP04UserRepo) GetRelativeUsersWithResident() ([]user_repository.AdminRelativeUser, error) {
	return nil, errors.New("not used")
}
func (f *fakeDrugPlanP04UserRepo) GetStaffIDMapByUserIDs(userIDs []string) (map[string]string, error) {
	return map[string]string{}, nil
}
func (f *fakeDrugPlanP04UserRepo) CreateOTP(otp *entities.OTP) error { return errors.New("not used") }
func (f *fakeDrugPlanP04UserRepo) GetOTPByUserID(userID string) (*entities.OTP, error) {
	return nil, errors.New("not used")
}
func (f *fakeDrugPlanP04UserRepo) DeleteOTP(userID string) error { return errors.New("not used") }
func (f *fakeDrugPlanP04UserRepo) StoreResetToken(temptoken *entities.TempToken) error {
	return errors.New("not used")
}
func (f *fakeDrugPlanP04UserRepo) GetResetToken(userID string) (string, error) {
	return "", errors.New("not used")
}
func (f *fakeDrugPlanP04UserRepo) DeleteResetToken(userID string) error {
	return errors.New("not used")
}

type fakeDrugPlanP04AuditRepo struct{}

func (f *fakeDrugPlanP04AuditRepo) CreateAuditLog(auditLog *entities.AuditLogs) (*entities.AuditLogs, error) {
	return auditLog, nil
}
func (f *fakeDrugPlanP04AuditRepo) GetAllAuditLogs() ([]*entities.AuditLogs, error) { return nil, nil }
func (f *fakeDrugPlanP04AuditRepo) SearchAuditLogs(search string) ([]*entities.AuditLogs, error) {
	return nil, nil
}
func (f *fakeDrugPlanP04AuditRepo) GetAuditLogByID(id string) (*entities.AuditLogs, error) {
	return nil, nil
}

type fakeDrugPlanP04Repo struct {
	medicine_repository.DrugRepository
	expiredDrugs         []*entities.PersonalDrug
	activeDrugs          []*entities.PersonalDrug
	createResults        []bool
	deletedDrugPlanIDs   []string
	deletedPersonalDrugs []string
	createCalls          int
	getExpiredCalls      int
	getActiveCalls       int
	residentExists       bool
	residentExistsErr    error
	getExpiredErr        error
	getActiveErr         error
	createErr            error
}

func (f *fakeDrugPlanP04Repo) GetExpiredAsNeededPersonalDrugs(now time.Time, residentID *string) ([]*entities.PersonalDrug, error) {
	f.getExpiredCalls++
	if f.getExpiredErr != nil {
		return nil, f.getExpiredErr
	}
	return f.expiredDrugs, nil
}

func (f *fakeDrugPlanP04Repo) DeleteDrugPlansByPdID(pdID string) error {
	f.deletedDrugPlanIDs = append(f.deletedDrugPlanIDs, pdID)
	return nil
}

func (f *fakeDrugPlanP04Repo) DeletePersonalDrug(id string) error {
	f.deletedPersonalDrugs = append(f.deletedPersonalDrugs, id)
	return nil
}

func (f *fakeDrugPlanP04Repo) GetActivePersonalDrugsForDate(now time.Time, residentID *string) ([]*entities.PersonalDrug, error) {
	f.getActiveCalls++
	if f.getActiveErr != nil {
		return nil, f.getActiveErr
	}
	return f.activeDrugs, nil
}

func (f *fakeDrugPlanP04Repo) CreateDrugPlanIfNotExistsForDate(plan *entities.DrugPlan, date time.Time) (bool, error) {
	f.createCalls++
	if f.createErr != nil {
		return false, f.createErr
	}
	if len(f.createResults) == 0 {
		return true, nil
	}
	result := f.createResults[0]
	f.createResults = f.createResults[1:]
	return result, nil
}

func (f *fakeDrugPlanP04Repo) ResidentExistsByID(residentID string) (bool, error) {
	if f.residentExistsErr != nil {
		return false, f.residentExistsErr
	}
	return f.residentExists, nil
}

func TestForceGenerateTodayDrugPlans_Success(t *testing.T) {
	userRepo := &fakeDrugPlanP04UserRepo{
		userByID: map[string]*entities.User{
			"user-1": {ID: "user-1", RoleID: "role-med"},
		},
		roleByID: map[string]*entities.Role{
			"role-med": {ID: "role-med", Name: user_constants.RoleMedicalStaff},
		},
	}
	repo := &fakeDrugPlanP04Repo{
		expiredDrugs:  []*entities.PersonalDrug{{ID: "pd-expired"}},
		activeDrugs:   []*entities.PersonalDrug{{ID: "pd-1"}, {ID: "pd-2"}},
		createResults: []bool{true, false},
	}
	uc := NewDrugUseCase(repo, &fakeDrugPlanP04AuditRepo{}, userRepo)

	result, err := uc.ForceGenerateTodayDrugPlans("user-1")

	assert.NoError(t, err)
	assert.Equal(t, 1, repo.getExpiredCalls)
	assert.Equal(t, 1, repo.getActiveCalls)
	assert.Equal(t, 2, repo.createCalls)
	assert.Equal(t, []string{"pd-expired"}, repo.deletedDrugPlanIDs)
	assert.Equal(t, []string{"pd-expired"}, repo.deletedPersonalDrugs)
	assert.Equal(t, 1, result.GeneratedCount)
	assert.Equal(t, 1, result.SkippedExistingCount)
	assert.Equal(t, 1, result.ExpiredDeletedCount)
	assert.Equal(t, "all", result.Scope)
}

func TestForceGenerateTodayDrugPlansByResidentID_Success(t *testing.T) {
	userRepo := &fakeDrugPlanP04UserRepo{
		userByID: map[string]*entities.User{
			"user-1": {ID: "user-1", RoleID: "role-med"},
		},
		roleByID: map[string]*entities.Role{
			"role-med": {ID: "role-med", Name: user_constants.RoleMedicalStaff},
		},
	}
	repo := &fakeDrugPlanP04Repo{
		residentExists: true,
		activeDrugs:    []*entities.PersonalDrug{{ID: "pd-1"}},
		createResults:  []bool{true},
	}
	uc := NewDrugUseCase(repo, &fakeDrugPlanP04AuditRepo{}, userRepo)

	result, err := uc.ForceGenerateTodayDrugPlansByResidentID("resident-1", "user-1")

	assert.NoError(t, err)
	assert.Equal(t, "resident", result.Scope)
	assert.Equal(t, "resident-1", result.ResidentID)
	assert.Equal(t, 1, result.GeneratedCount)
	assert.Equal(t, 0, result.SkippedExistingCount)
	assert.Equal(t, 0, result.ExpiredDeletedCount)
}

func TestForceGenerateTodayDrugPlansByResidentID_ErrorWhenResidentMissing(t *testing.T) {
	userRepo := &fakeDrugPlanP04UserRepo{
		userByID: map[string]*entities.User{
			"user-1": {ID: "user-1", RoleID: "role-med"},
		},
		roleByID: map[string]*entities.Role{
			"role-med": {ID: "role-med", Name: user_constants.RoleMedicalStaff},
		},
	}
	repo := &fakeDrugPlanP04Repo{residentExists: false}
	uc := NewDrugUseCase(repo, &fakeDrugPlanP04AuditRepo{}, userRepo)

	result, err := uc.ForceGenerateTodayDrugPlansByResidentID("resident-1", "user-1")

	assert.Nil(t, result)
	assert.EqualError(t, err, "resident not found")
}

var _ user_repository.UserRepository = (*fakeDrugPlanP04UserRepo)(nil)
var _ audit_repository.AuditLogRepository = (*fakeDrugPlanP04AuditRepo)(nil)
