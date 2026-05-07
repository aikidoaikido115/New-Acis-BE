package usecases

import (
	"errors"
	"testing"
	"time"

	"github.com/aikidoaikido115/New-Acis-BE/configs"
	audit_repository "github.com/aikidoaikido115/New-Acis-BE/modules/audit_logs/repositories"
	"github.com/aikidoaikido115/New-Acis-BE/modules/entities"
	user_repository "github.com/aikidoaikido115/New-Acis-BE/modules/user/repositories"
	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
)

type fakeUserP04Repo struct {
	userByID       map[string]*entities.User
	userByUsername map[string]*entities.User
	userByEmail    map[string]*entities.User
	otpByUserID    map[string]*entities.OTP
	resetTokenByID map[string]string
	createdOTPs    []*entities.OTP
	deletedOTPs    []string
	createdTokens  []*entities.TempToken
	deletedTokens  []string
	updatedUsers   []*entities.User
}

func (f *fakeUserP04Repo) CreateUser(user *entities.User) (*entities.User, error) {
	return nil, errors.New("not used")
}
func (f *fakeUserP04Repo) CreateStaff(user *entities.User, staff *entities.Staff) (*entities.Staff, error) {
	return nil, errors.New("not used")
}
func (f *fakeUserP04Repo) CreateStaffFile(staffFile *entities.StaffsFiles) (*entities.StaffsFiles, error) {
	return nil, errors.New("not used")
}
func (f *fakeUserP04Repo) GetUserByEmail(email string) (*entities.User, error) {
	if user, ok := f.userByEmail[email]; ok {
		return user, nil
	}
	return nil, errors.New("user not found")
}
func (f *fakeUserP04Repo) GetUserByID(id string) (*entities.User, error) {
	if user, ok := f.userByID[id]; ok {
		return user, nil
	}
	return nil, errors.New("user not found")
}
func (f *fakeUserP04Repo) GetUsersByFirstAndLastName(firstName string, lastName string) ([]*entities.User, error) {
	return nil, errors.New("not used")
}
func (f *fakeUserP04Repo) GetStaffByID(id string) (*entities.Staff, error) {
	return nil, errors.New("not used")
}
func (f *fakeUserP04Repo) GetStaffByUserID(userID string) (*entities.Staff, error) {
	return nil, errors.New("not used")
}
func (f *fakeUserP04Repo) GetStaffFileByID(id string) (*entities.StaffsFiles, error) {
	return nil, errors.New("not used")
}
func (f *fakeUserP04Repo) GetUserByUsername(username string) (*entities.User, error) {
	if user, ok := f.userByUsername[username]; ok {
		return user, nil
	}
	return nil, errors.New("user not found")
}
func (f *fakeUserP04Repo) GetRoleByName(roleName string) (*entities.Role, error) {
	return nil, errors.New("not used")
}
func (f *fakeUserP04Repo) GetRoleByID(roleID string) (*entities.Role, error) {
	return nil, errors.New("not used")
}
func (f *fakeUserP04Repo) UsernameExists(username string) (bool, error) {
	return false, errors.New("not used")
}
func (f *fakeUserP04Repo) EmailExists(email string) (bool, error) {
	return false, errors.New("not used")
}
func (f *fakeUserP04Repo) GetAllUsers() ([]*entities.User, error) { return nil, errors.New("not used") }
func (f *fakeUserP04Repo) UpdateUserByID(user *entities.User) error {
	f.updatedUsers = append(f.updatedUsers, user)
	return nil
}
func (f *fakeUserP04Repo) UpdateUserApprovalByID(userID string, isApprove bool) error {
	return errors.New("not used")
}
func (f *fakeUserP04Repo) DeleteStaffAndUserByStaffID(staffID string) error {
	return errors.New("not used")
}
func (f *fakeUserP04Repo) DeleteRelativeAndUserByUserID(userID string) error {
	return errors.New("not used")
}
func (f *fakeUserP04Repo) GetRelativeUserByUserID(userID string) (*user_repository.AdminRelativeUser, error) {
	return nil, errors.New("not used")
}
func (f *fakeUserP04Repo) GetRelativeUsersWithResident() ([]user_repository.AdminRelativeUser, error) {
	return nil, errors.New("not used")
}
func (f *fakeUserP04Repo) GetStaffIDMapByUserIDs(userIDs []string) (map[string]string, error) {
	return map[string]string{}, nil
}
func (f *fakeUserP04Repo) CreateOTP(otp *entities.OTP) error {
	f.createdOTPs = append(f.createdOTPs, otp)
	return nil
}
func (f *fakeUserP04Repo) GetOTPByUserID(userID string) (*entities.OTP, error) {
	if otp, ok := f.otpByUserID[userID]; ok {
		return otp, nil
	}
	return nil, errors.New("otp not found")
}
func (f *fakeUserP04Repo) DeleteOTP(userID string) error {
	f.deletedOTPs = append(f.deletedOTPs, userID)
	return nil
}
func (f *fakeUserP04Repo) StoreResetToken(temptoken *entities.TempToken) error {
	f.createdTokens = append(f.createdTokens, temptoken)
	return nil
}
func (f *fakeUserP04Repo) GetResetToken(userID string) (string, error) {
	if token, ok := f.resetTokenByID[userID]; ok {
		return token, nil
	}
	return "", errors.New("reset token not found")
}
func (f *fakeUserP04Repo) DeleteResetToken(userID string) error {
	f.deletedTokens = append(f.deletedTokens, userID)
	return nil
}

type fakeUserP04AuditRepo struct{}

func (f *fakeUserP04AuditRepo) CreateAuditLog(auditLog *entities.AuditLogs) (*entities.AuditLogs, error) {
	return auditLog, nil
}
func (f *fakeUserP04AuditRepo) GetAllAuditLogs() ([]*entities.AuditLogs, error) { return nil, nil }
func (f *fakeUserP04AuditRepo) SearchAuditLogs(search string) ([]*entities.AuditLogs, error) {
	return nil, nil
}
func (f *fakeUserP04AuditRepo) GetAuditLogByID(id string) (*entities.AuditLogs, error) {
	return nil, nil
}

func makeHash(t *testing.T, password string) string {
	t.Helper()
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	assert.NoError(t, err)
	return string(hash)
}

func makeToken(t *testing.T, secret string, userID string, expiry time.Time) string {
	t.Helper()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": userID,
		"exp":     expiry.Unix(),
		"iat":     time.Now().Unix(),
		"jti":     "token-1",
	})
	value, err := token.SignedString([]byte(secret))
	assert.NoError(t, err)
	return value
}

func TestLogin_SuccessByUsername(t *testing.T) {
	repo := &fakeUserP04Repo{
		userByUsername: map[string]*entities.User{
			"staff1": {ID: "user-1", Username: "staff1", Password: makeHash(t, "secret"), IsApprove: true},
		},
	}
	uc := NewUserUseCase(repo, &fakeUserP04AuditRepo{}, configs.JWT{Secret: "secret"}, configs.Supabase{}, configs.Mail{})

	token, user, err := uc.Login("staff1", "", "secret", true)

	assert.NoError(t, err)
	assert.NotEmpty(t, token)
	assert.Equal(t, "user-1", user.ID)
}

func TestLogin_ErrorWhenPasswordMismatch(t *testing.T) {
	repo := &fakeUserP04Repo{
		userByUsername: map[string]*entities.User{
			"staff1": {ID: "user-1", Username: "staff1", Password: makeHash(t, "secret"), IsApprove: true},
		},
	}
	uc := NewUserUseCase(repo, &fakeUserP04AuditRepo{}, configs.JWT{Secret: "secret"}, configs.Supabase{}, configs.Mail{})

	token, user, err := uc.Login("staff1", "", "wrong", false)

	assert.Empty(t, token)
	assert.Nil(t, user)
	assert.ErrorIs(t, err, ErrInvalidCredentials)
}

func TestResetPassword_Success(t *testing.T) {
	repo := &fakeUserP04Repo{
		userByID: map[string]*entities.User{
			"user-1": {ID: "user-1", Password: makeHash(t, "old-pass")},
		},
	}
	uc := NewUserUseCase(repo, &fakeUserP04AuditRepo{}, configs.JWT{Secret: "secret"}, configs.Supabase{}, configs.Mail{})

	err := uc.ResetPassword("user-1", "old-pass", "new-pass")

	assert.NoError(t, err)
	assert.Len(t, repo.updatedUsers, 1)
	assert.NoError(t, bcrypt.CompareHashAndPassword([]byte(repo.updatedUsers[0].Password), []byte("new-pass")))
}

func TestForgotPassword_CreatesOtpAndDeletesExisting(t *testing.T) {
	repo := &fakeUserP04Repo{
		userByEmail: map[string]*entities.User{
			"user@example.com": {ID: "user-1", Email: "user@example.com"},
		},
		otpByUserID: map[string]*entities.OTP{
			"user-1": {UserID: "user-1", OTP: "123456", ExpiresAt: time.Now().Add(5 * time.Minute)},
		},
	}
	uc := NewUserUseCase(repo, &fakeUserP04AuditRepo{}, configs.JWT{Secret: "secret"}, configs.Supabase{}, configs.Mail{})

	err := uc.ForgotPassword("user@example.com")

	assert.NoError(t, err)
	assert.Equal(t, []string{"user-1"}, repo.deletedOTPs)
	if assert.Len(t, repo.createdOTPs, 1) {
		assert.Equal(t, "user-1", repo.createdOTPs[0].UserID)
		assert.NotEmpty(t, repo.createdOTPs[0].OTP)
	}
}

func TestVerifyOTP_Success(t *testing.T) {
	repo := &fakeUserP04Repo{
		userByEmail: map[string]*entities.User{
			"user@example.com": {ID: "user-1", Email: "user@example.com"},
		},
		otpByUserID: map[string]*entities.OTP{
			"user-1": {UserID: "user-1", OTP: "654321", ExpiresAt: time.Now().Add(5 * time.Minute)},
		},
	}
	uc := NewUserUseCase(repo, &fakeUserP04AuditRepo{}, configs.JWT{Secret: "secret"}, configs.Supabase{}, configs.Mail{})

	err := uc.VerifyOTP("user@example.com", "654321")

	assert.NoError(t, err)
	assert.Equal(t, []string{"user-1"}, repo.deletedOTPs)
	if assert.Len(t, repo.createdTokens, 1) {
		assert.Equal(t, "user-1", repo.createdTokens[0].UserID)
		assert.NotEmpty(t, repo.createdTokens[0].Token)
	}
}

func TestChangePassword_Success(t *testing.T) {
	repo := &fakeUserP04Repo{
		userByEmail: map[string]*entities.User{
			"user@example.com": {ID: "user-1", Email: "user@example.com", Password: makeHash(t, "old-pass")},
		},
		resetTokenByID: map[string]string{
			"user-1": makeToken(t, "secret", "user-1", time.Now().Add(5*time.Minute)),
		},
	}
	uc := NewUserUseCase(repo, &fakeUserP04AuditRepo{}, configs.JWT{Secret: "secret"}, configs.Supabase{}, configs.Mail{})

	err := uc.ChangePassword("user@example.com", "new-pass")

	assert.NoError(t, err)
	assert.Len(t, repo.updatedUsers, 1)
	assert.Equal(t, []string{"user-1"}, repo.deletedTokens)
	assert.NoError(t, bcrypt.CompareHashAndPassword([]byte(repo.updatedUsers[0].Password), []byte("new-pass")))
}

var _ audit_repository.AuditLogRepository = (*fakeUserP04AuditRepo)(nil)
