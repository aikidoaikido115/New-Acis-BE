package usecases

import (
	"errors"
	"mime/multipart"
	"testing"
	"time"

	"github.com/aikidoaikido115/New-Acis-BE/configs"
	"github.com/aikidoaikido115/New-Acis-BE/modules/entities"
	"github.com/aikidoaikido115/New-Acis-BE/modules/user/repositories"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/google/uuid"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

// Mock Audit Log Repository
type MockAuditLogRepository struct {
	mock.Mock
}

func (m *MockAuditLogRepository) CreateAuditLog(log *entities.AuditLogs) (*entities.AuditLogs, error) {
	args := m.Called(log)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entities.AuditLogs), args.Error(1)
}

func (m *MockAuditLogRepository) GetAuditLogByID(id string) (*entities.AuditLogs, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entities.AuditLogs), args.Error(1)
}

// ========== GET ALL USERS TESTS ==========

func TestGetAllUsers_Success(t *testing.T) {
	usecase, mockRepo, _ := setupTestUseCase()

	users := []*entities.User{
		{ID: "1", Username: "user1", Email: "user1@example.com"},
		{ID: "2", Username: "user2", Email: "user2@example.com"},
	}

	mockRepo.On("GetAllUsers").Return(users, nil)

	result, err := usecase.GetAllUsers()

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Len(t, result, 2)
	mockRepo.AssertExpectations(t)
}

func TestGetAllUsers_Empty(t *testing.T) {
	usecase, mockRepo, _ := setupTestUseCase()

	users := []*entities.User{}

	mockRepo.On("GetAllUsers").Return(users, nil)

	result, err := usecase.GetAllUsers()

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Len(t, result, 0)
	mockRepo.AssertExpectations(t)
}

func TestGetAllUsers_Error(t *testing.T) {
	usecase, mockRepo, _ := setupTestUseCase()

	mockRepo.On("GetAllUsers").Return(nil, errors.New("database error"))

	result, err := usecase.GetAllUsers()

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, "failed to retrieve all users", err.Error())
	mockRepo.AssertExpectations(t)
}

// ========== UPDATE USER TESTS ==========

func TestUpdateUserByID_Success(t *testing.T) {
	usecase, mockRepo, mockAuditRepo := setupTestUseCase()

	existingUser := &entities.User{
		ID:       "user-id-1",
		Username: "oldusername",
		Email:    "test@example.com",
	}

	updatedUser := &entities.User{
		FirstName: "John",
		LastName:  "Doe",
		Nickname:  "JD",
		Gender:    "Male",
	}

	mockRepo.On("GetUserByID", "user-id-1").Return(existingUser, nil)
	mockRepo.On("UpdateUserByID", mock.AnythingOfType("*entities.User")).Return(nil)
	mockAuditRepo.On("CreateAuditLog", mock.AnythingOfType("*entities.AuditLogs")).Return(&entities.AuditLogs{}, nil)

	result, err := usecase.UpdateUserByID("user-id-1", updatedUser, nil)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "John", result.FirstName)
	assert.Equal(t, "Doe", result.LastName)
	mockRepo.AssertExpectations(t)
}

func TestUpdateUserByID_UserNotFound(t *testing.T) {
	usecase, mockRepo, _ := setupTestUseCase()

	updatedUser := &entities.User{FirstName: "John"}

	mockRepo.On("GetUserByID", "invalid-id").Return(nil, errors.New("not found"))

	result, err := usecase.UpdateUserByID("invalid-id", updatedUser, nil)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, "user not found", err.Error())
	mockRepo.AssertExpectations(t)
}

func TestUpdateUserByID_UsernameAlreadyTaken(t *testing.T) {
	usecase, mockRepo, _ := setupTestUseCase()

	existingUser := &entities.User{
		ID:       "user-id-1",
		Username: "oldusername",
	}

	updatedUser := &entities.User{Username: "takenusername"}

	mockRepo.On("GetUserByID", "user-id-1").Return(existingUser, nil)
	mockRepo.On("UsernameExists", "takenusername").Return(true, nil)

	result, err := usecase.UpdateUserByID("user-id-1", updatedUser, nil)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, "username already taken", err.Error())
	mockRepo.AssertExpectations(t)
}

func TestUpdateUserByID_UsernameChangeLimit(t *testing.T) {
	usecase, mockRepo, _ := setupTestUseCase()

	existingUser := &entities.User{
		ID:                "user-id-1",
		Username:          "oldusername",
		NumberOfUsernames: 1,
	}

	updatedUser := &entities.User{Username: "newusername"}

	mockRepo.On("GetUserByID", "user-id-1").Return(existingUser, nil)
	mockRepo.On("UsernameExists", "newusername").Return(false, nil)

	result, err := usecase.UpdateUserByID("user-id-1", updatedUser, nil)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, "username can change only once", err.Error())
	mockRepo.AssertExpectations(t)
}

// ========== CHANGE PASSWORD TESTS ==========

func TestChangePassword_Success(t *testing.T) {
    usecase, mockRepo, _ := setupTestUseCase()

    user := &entities.User{
        ID:    "user-id-1",
        Email: "test@example.com",
    }

    // Create a valid JWT token
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
        "user_id": user.ID,
        "exp":     time.Now().Add(time.Minute * 5).Unix(),
        "iat":     time.Now().Unix(),
        "jti":     uuid.New().String(),
    })
    tokenString, _ := token.SignedString([]byte("test-secret-key"))

    mockRepo.On("GetUserByEmail", "test@example.com").Return(user, nil)
    mockRepo.On("GetResetToken", "user-id-1").Return(tokenString, nil) // Return valid JWT
    mockRepo.On("UpdateUserByID", mock.AnythingOfType("*entities.User")).Return(nil)
    mockRepo.On("DeleteResetToken", "user-id-1").Return(nil)

    err := usecase.ChangePassword("test@example.com", "newpassword123")

    assert.NoError(t, err)
    mockRepo.AssertExpectations(t)
}

func TestChangePassword_UserNotFound(t *testing.T) {
	usecase, mockRepo, _ := setupTestUseCase()

	mockRepo.On("GetUserByEmail", "notfound@example.com").Return(nil, errors.New("not found"))

	err := usecase.ChangePassword("notfound@example.com", "newpassword")

	assert.Error(t, err)
	assert.Equal(t, "user not found", err.Error())
	mockRepo.AssertExpectations(t)
}

func TestChangePassword_TokenNotFound(t *testing.T) {
	usecase, mockRepo, _ := setupTestUseCase()

	user := &entities.User{
		ID:    "user-id-1",
		Email: "test@example.com",
	}

	mockRepo.On("GetUserByEmail", "test@example.com").Return(user, nil)
	mockRepo.On("GetResetToken", "user-id-1").Return("", errors.New("token not found"))

	err := usecase.ChangePassword("test@example.com", "newpassword")

	assert.Error(t, err)
	assert.Equal(t, "reset token not found", err.Error())
	mockRepo.AssertExpectations(t)
}

// ========== CREATE STAFF FILE TESTS ==========

func TestCreateStaffFile_UserNotFound(t *testing.T) {
	usecase, mockRepo, _ := setupTestUseCase()

	mockRepo.On("GetUserByID", "invalid-id").Return(nil, errors.New("not found"))

	result, err := usecase.CreateStaffFile("invalid-id", nil)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, "user not found", err.Error())
	mockRepo.AssertExpectations(t)
}

func TestCreateStaffFile_UserNotStaff(t *testing.T) {
	usecase, mockRepo, _ := setupTestUseCase()

	user := &entities.User{
		ID:   "user-id-1",
		Role: entities.Role{Name: "Relative"},
	}

	mockRepo.On("GetUserByID", "user-id-1").Return(user, nil)

	result, err := usecase.CreateStaffFile("user-id-1", nil)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, "user is not staff", err.Error())
	mockRepo.AssertExpectations(t)
}

func TestCreateStaffFile_NoFilesProvided(t *testing.T) {
	usecase, mockRepo, _ := setupTestUseCase()

	user := &entities.User{
		ID:   "user-id-1",
		Role: entities.Role{Name: "Medical Staff"},
	}

	staff := &entities.Staff{ID: "staff-1", UserID: "user-id-1"}

	mockRepo.On("GetUserByID", "user-id-1").Return(user, nil)
	mockRepo.On("GetStaffByUserID", "user-id-1").Return(staff, nil)

	result, err := usecase.CreateStaffFile("user-id-1", []*multipart.FileHeader{})

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, "no files provided", err.Error())
	mockRepo.AssertExpectations(t)
}

func TestCreateStaffFile_StaffNotFound(t *testing.T) {
	usecase, mockRepo, _ := setupTestUseCase()

	user := &entities.User{
		ID:   "user-id-1",
		Role: entities.Role{Name: "Medical Staff"},
	}

	mockRepo.On("GetUserByID", "user-id-1").Return(user, nil)
	mockRepo.On("GetStaffByUserID", "user-id-1").Return(nil, errors.New("staff not found"))

	result, err := usecase.CreateStaffFile("user-id-1", []*multipart.FileHeader{})

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, "staff not found", err.Error())
	mockRepo.AssertExpectations(t)
}

func TestRegister_InvalidEmail(t *testing.T) {
	usecase, mockRepo, _ := setupTestUseCase()

	user := &entities.User{
		Username: "testuser",
		Email:    "invalid-email",
		Password: "password123",
	}

	result, err := usecase.Register(user, "Medical Staff")

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, "invalid email format", err.Error())
	mockRepo.AssertNotCalled(t, "GetRoleByName")
}

func TestRegister_Success_KitchenStaff(t *testing.T) {
	usecase, mockRepo, _ := setupTestUseCase()

	role := &entities.Role{ID: "2", Name: "Kitchen Staff"}
	user := &entities.User{
		Username: "kitchenstaff",
		Email:    "kitchen@example.com",
		Password: "password123",
	}
	staff := &entities.Staff{ID: "staff-2"}

	mockRepo.On("GetRoleByName", "Kitchen Staff").Return(role, nil)
	mockRepo.On("UsernameExists", "kitchenstaff").Return(false, nil)
	mockRepo.On("EmailExists", "kitchen@example.com").Return(false, nil)
	mockRepo.On("CreateUser", mock.AnythingOfType("*entities.User")).Return(user, nil)
	mockRepo.On("CreateStaff", mock.AnythingOfType("*entities.User"), mock.AnythingOfType("*entities.Staff")).Return(staff, nil)

	result, err := usecase.Register(user, "Kitchen Staff")

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "kitchenstaff", result.Username)
	mockRepo.AssertExpectations(t)
}

func TestForgotPassword_ExistingOTP(t *testing.T) {
	usecase, mockRepo, _ := setupTestUseCase()

	user := &entities.User{
		ID:       "user-id-1",
		Email:    "test@example.com",
		Username: "testuser",
	}

	existingOTP := &entities.OTP{
		UserID:    "user-id-1",
		OTP:       "654321",
		ExpiresAt: time.Now().Add(3 * time.Minute),
	}

	mockRepo.On("GetUserByEmail", "test@example.com").Return(user, nil)
	mockRepo.On("GetOTPByUserID", "user-id-1").Return(existingOTP, nil)
	mockRepo.On("DeleteOTP", "user-id-1").Return(nil)
	mockRepo.On("CreateOTP", mock.AnythingOfType("*entities.OTP")).Return(nil)

	err := usecase.ForgotPassword("test@example.com")

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestVerifyOTP_UserNotFound(t *testing.T) {
	usecase, mockRepo, _ := setupTestUseCase()

	mockRepo.On("GetUserByEmail", "notfound@example.com").Return(nil, errors.New("not found"))

	err := usecase.VerifyOTP("notfound@example.com", "123456")

	assert.Error(t, err)
	assert.Equal(t, "user not found", err.Error())
	mockRepo.AssertExpectations(t)
}

func TestVerifyOTP_OTPNotFound(t *testing.T) {
	usecase, mockRepo, _ := setupTestUseCase()

	user := &entities.User{
		ID:    "user-id-1",
		Email: "test@example.com",
	}

	mockRepo.On("GetUserByEmail", "test@example.com").Return(user, nil)
	mockRepo.On("GetOTPByUserID", "user-id-1").Return(nil, errors.New("not found"))

	err := usecase.VerifyOTP("test@example.com", "123456")

	assert.Error(t, err)
	assert.Equal(t, "OTP not found for user", err.Error())
	mockRepo.AssertExpectations(t)
}

// Helper function to setup test dependencies
func setupTestUseCase() (UserUsecase, *repositories.MockUserRepository, *MockAuditLogRepository) {
	mockRepo := new(repositories.MockUserRepository)
	mockAuditRepo := new(MockAuditLogRepository)

	jwt := configs.JWT{Secret: "test-secret-key"}
	supa := configs.Supabase{}
	mail := configs.Mail{}

	usecase := NewUserUseCase(mockRepo, mockAuditRepo, jwt, supa, mail)
	return usecase, mockRepo, mockAuditRepo
}

// ========== REGISTER TESTS ==========

func TestRegister_Success_MedicalStaff(t *testing.T) {
	usecase, mockRepo, _ := setupTestUseCase()

	role := &entities.Role{ID: "1", Name: "Medical Staff"}
	user := &entities.User{
		Username: "testuser",
		Email:    "test@example.com",
		Password: "password123",
	}
	staff := &entities.Staff{ID: "staff-1"}

	mockRepo.On("GetRoleByName", "Medical Staff").Return(role, nil)
	mockRepo.On("UsernameExists", "testuser").Return(false, nil)
	mockRepo.On("EmailExists", "test@example.com").Return(false, nil)
	mockRepo.On("CreateUser", mock.AnythingOfType("*entities.User")).Return(user, nil)
	mockRepo.On("CreateStaff", mock.AnythingOfType("*entities.User"), mock.AnythingOfType("*entities.Staff")).Return(staff, nil)

	result, err := usecase.Register(user, "Medical Staff")

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "testuser", result.Username)
	mockRepo.AssertExpectations(t)
}

func TestRegister_Success_Relative(t *testing.T) {
	usecase, mockRepo, _ := setupTestUseCase()

	role := &entities.Role{ID: "3", Name: "Relative"}
	user := &entities.User{
		Username: "relative",
		Email:    "relative@example.com",
		Password: "password123",
	}

	mockRepo.On("GetRoleByName", "Relative").Return(role, nil)
	mockRepo.On("UsernameExists", "relative").Return(false, nil)
	mockRepo.On("EmailExists", "relative@example.com").Return(false, nil)
	mockRepo.On("CreateUser", mock.AnythingOfType("*entities.User")).Return(user, nil)

	result, err := usecase.Register(user, "Relative")

	assert.NoError(t, err)
	assert.NotNil(t, result)
	mockRepo.AssertExpectations(t)
}

func TestRegister_UsernameExists(t *testing.T) {
	usecase, mockRepo, _ := setupTestUseCase()

	role := &entities.Role{ID: "1", Name: "Medical Staff"}
	user := &entities.User{
		Username: "existinguser",
		Email:    "test@example.com",
		Password: "password123",
	}

	mockRepo.On("GetRoleByName", "Medical Staff").Return(role, nil)
	mockRepo.On("UsernameExists", "existinguser").Return(true, nil)

	result, err := usecase.Register(user, "Medical Staff")

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, "username already exists", err.Error())
	mockRepo.AssertExpectations(t)
}

func TestRegister_EmailExists(t *testing.T) {
	usecase, mockRepo, _ := setupTestUseCase()

	role := &entities.Role{ID: "1", Name: "Medical Staff"}
	user := &entities.User{
		Username: "newuser",
		Email:    "existing@example.com",
		Password: "password123",
	}

	mockRepo.On("GetRoleByName", "Medical Staff").Return(role, nil)
	mockRepo.On("UsernameExists", "newuser").Return(false, nil)
	mockRepo.On("EmailExists", "existing@example.com").Return(true, nil)

	result, err := usecase.Register(user, "Medical Staff")

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, "email already exists", err.Error())
	mockRepo.AssertExpectations(t)
}

func TestRegister_RoleNotFound(t *testing.T) {
	usecase, mockRepo, _ := setupTestUseCase()

	user := &entities.User{
		Username: "testuser",
		Email:    "test@example.com",
		Password: "password123",
	}

	mockRepo.On("GetRoleByName", "InvalidRole").Return(nil, errors.New("role not found"))

	result, err := usecase.Register(user, "InvalidRole")

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, "role not found", err.Error())
	mockRepo.AssertExpectations(t)
}

// ========== LOGIN TESTS ==========

func TestLogin_Success(t *testing.T) {
	usecase, mockRepo, _ := setupTestUseCase()

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)
	user := &entities.User{
		ID:       "user-id-1",
		Username: "testuser",
		Password: string(hashedPassword),
	}

	mockRepo.On("GetUserByUsername", "testuser").Return(user, nil)

	token, resultUser, err := usecase.Login("testuser", "password123")

	assert.NoError(t, err)
	assert.NotEmpty(t, token)
	assert.Equal(t, "testuser", resultUser.Username)
	mockRepo.AssertExpectations(t)
}

func TestLogin_InvalidUsername(t *testing.T) {
	usecase, mockRepo, _ := setupTestUseCase()

	mockRepo.On("GetUserByUsername", "wronguser").Return(nil, errors.New("user not found"))

	token, resultUser, err := usecase.Login("wronguser", "password123")

	assert.Error(t, err)
	assert.Empty(t, token)
	assert.Nil(t, resultUser)
	assert.Equal(t, "invalid username", err.Error())
	mockRepo.AssertExpectations(t)
}

func TestLogin_InvalidPassword(t *testing.T) {
	usecase, mockRepo, _ := setupTestUseCase()

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("correctpassword"), bcrypt.DefaultCost)
	user := &entities.User{
		ID:       "user-id-1",
		Username: "testuser",
		Password: string(hashedPassword),
	}

	mockRepo.On("GetUserByUsername", "testuser").Return(user, nil)

	token, resultUser, err := usecase.Login("testuser", "wrongpassword")

	assert.Error(t, err)
	assert.Empty(t, token)
	assert.Nil(t, resultUser)
	assert.Equal(t, "invalid password", err.Error())
	mockRepo.AssertExpectations(t)
}

// ========== RESET PASSWORD TESTS ==========

func TestResetPassword_Success(t *testing.T) {
	usecase, mockRepo, _ := setupTestUseCase()

	hashedOldPassword, _ := bcrypt.GenerateFromPassword([]byte("oldpassword"), bcrypt.DefaultCost)
	user := &entities.User{
		ID:       "user-id-1",
		Password: string(hashedOldPassword),
	}

	mockRepo.On("GetUserByID", "user-id-1").Return(user, nil)
	mockRepo.On("UpdateUserByID", mock.AnythingOfType("*entities.User")).Return(nil)

	err := usecase.ResetPassword("user-id-1", "oldpassword", "newpassword123")

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestResetPassword_UserNotFound(t *testing.T) {
	usecase, mockRepo, _ := setupTestUseCase()

	mockRepo.On("GetUserByID", "invalid-id").Return(nil, errors.New("user not found"))

	err := usecase.ResetPassword("invalid-id", "oldpassword", "newpassword")

	assert.Error(t, err)
	assert.Equal(t, "user invalid", err.Error())
	mockRepo.AssertExpectations(t)
}

func TestResetPassword_InvalidOldPassword(t *testing.T) {
	usecase, mockRepo, _ := setupTestUseCase()

	hashedOldPassword, _ := bcrypt.GenerateFromPassword([]byte("oldpassword"), bcrypt.DefaultCost)
	user := &entities.User{
		ID:       "user-id-1",
		Password: string(hashedOldPassword),
	}

	mockRepo.On("GetUserByID", "user-id-1").Return(user, nil)

	err := usecase.ResetPassword("user-id-1", "wrongoldpassword", "newpassword")

	assert.Error(t, err)
	assert.Equal(t, "old password invalid", err.Error())
	mockRepo.AssertExpectations(t)
}

// ========== FORGOT PASSWORD TESTS ==========

func TestForgotPassword_Success(t *testing.T) {
	usecase, mockRepo, _ := setupTestUseCase()

	user := &entities.User{
		ID:       "user-id-1",
		Email:    "test@example.com",
		Username: "testuser",
	}

	mockRepo.On("GetUserByEmail", "test@example.com").Return(user, nil)
	mockRepo.On("GetOTPByUserID", "user-id-1").Return(nil, errors.New("not found"))
	mockRepo.On("CreateOTP", mock.AnythingOfType("*entities.OTP")).Return(nil)

	err := usecase.ForgotPassword("test@example.com")

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestForgotPassword_UserNotFound(t *testing.T) {
	usecase, mockRepo, _ := setupTestUseCase()

	mockRepo.On("GetUserByEmail", "notfound@example.com").Return(nil, errors.New("user not found"))

	err := usecase.ForgotPassword("notfound@example.com")

	assert.Error(t, err)
	assert.Equal(t, "user not found", err.Error())
	mockRepo.AssertExpectations(t)
}

// ========== VERIFY OTP TESTS ==========

func TestVerifyOTP_Success(t *testing.T) {
	usecase, mockRepo, _ := setupTestUseCase()

	user := &entities.User{
		ID:    "user-id-1",
		Email: "test@example.com",
	}

	otp := &entities.OTP{
		UserID:    "user-id-1",
		OTP:       "123456",
		ExpiresAt: time.Now().Add(5 * time.Minute),
	}

	mockRepo.On("GetUserByEmail", "test@example.com").Return(user, nil)
	mockRepo.On("GetOTPByUserID", "user-id-1").Return(otp, nil)
	mockRepo.On("DeleteOTP", "user-id-1").Return(nil)
	mockRepo.On("StoreResetToken", mock.AnythingOfType("*entities.TempToken")).Return(nil)

	err := usecase.VerifyOTP("test@example.com", "123456")

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestVerifyOTP_InvalidCode(t *testing.T) {
	usecase, mockRepo, _ := setupTestUseCase()

	user := &entities.User{
		ID:    "user-id-1",
		Email: "test@example.com",
	}

	otp := &entities.OTP{
		UserID:    "user-id-1",
		OTP:       "123456",
		ExpiresAt: time.Now().Add(5 * time.Minute),
	}

	mockRepo.On("GetUserByEmail", "test@example.com").Return(user, nil)
	mockRepo.On("GetOTPByUserID", "user-id-1").Return(otp, nil)

	err := usecase.VerifyOTP("test@example.com", "wrong-otp")

	assert.Error(t, err)
	assert.Equal(t, "invalid OTP code", err.Error())
	mockRepo.AssertExpectations(t)
}

func TestVerifyOTP_Expired(t *testing.T) {
	usecase, mockRepo, _ := setupTestUseCase()

	user := &entities.User{
		ID:    "user-id-1",
		Email: "test@example.com",
	}

	otp := &entities.OTP{
		UserID:    "user-id-1",
		OTP:       "123456",
		ExpiresAt: time.Now().Add(-5 * time.Minute), // Expired
	}

	mockRepo.On("GetUserByEmail", "test@example.com").Return(user, nil)
	mockRepo.On("GetOTPByUserID", "user-id-1").Return(otp, nil)
	mockRepo.On("DeleteOTP", "user-id-1").Return(nil)

	err := usecase.VerifyOTP("test@example.com", "123456")

	assert.Error(t, err)
	assert.Equal(t, "OTP has expired", err.Error())
	mockRepo.AssertExpectations(t)
}

// ========== GET USER BY ID TESTS ==========

func TestGetUserByID_Success(t *testing.T) {
	usecase, mockRepo, _ := setupTestUseCase()

	user := &entities.User{
		ID:       "user-id-1",
		Username: "testuser",
		Email:    "test@example.com",
	}

	mockRepo.On("GetUserByID", "user-id-1").Return(user, nil)

	result, err := usecase.GetUserByID("user-id-1")

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "testuser", result.Username)
	mockRepo.AssertExpectations(t)
}

func TestGetUserByID_NotFound(t *testing.T) {
	usecase, mockRepo, _ := setupTestUseCase()

	mockRepo.On("GetUserByID", "invalid-id").Return(nil, errors.New("user not found"))

	result, err := usecase.GetUserByID("invalid-id")

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, "user not found", err.Error())
	mockRepo.AssertExpectations(t)
}
