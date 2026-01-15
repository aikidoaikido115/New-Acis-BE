package repositories

import (
    "github.com/aikidoaikido115/New-Acis-BE/modules/entities"
    "github.com/stretchr/testify/mock"
)

type MockUserRepository struct {
    mock.Mock
}

func (m *MockUserRepository) CreateUser(user *entities.User) (*entities.User, error) {
    args := m.Called(user)
    if args.Get(0) == nil {
        return nil, args.Error(1)
    }
    return args.Get(0).(*entities.User), args.Error(1)
}

func (m *MockUserRepository) GetUserByUsername(username string) (*entities.User, error) {
    args := m.Called(username)
    if args.Get(0) == nil {
        return nil, args.Error(1)
    }
    return args.Get(0).(*entities.User), args.Error(1)
}

func (m *MockUserRepository) GetUserByEmail(email string) (*entities.User, error) {
    args := m.Called(email)
    if args.Get(0) == nil {
        return nil, args.Error(1)
    }
    return args.Get(0).(*entities.User), args.Error(1)
}

func (m *MockUserRepository) GetUserByID(id string) (*entities.User, error) {
    args := m.Called(id)
    if args.Get(0) == nil {
        return nil, args.Error(1)
    }
    return args.Get(0).(*entities.User), args.Error(1)
}

func (m *MockUserRepository) GetRoleByName(roleName string) (*entities.Role, error) {
    args := m.Called(roleName)
    if args.Get(0) == nil {
        return nil, args.Error(1)
    }
    return args.Get(0).(*entities.Role), args.Error(1)
}

func (m *MockUserRepository) UsernameExists(username string) (bool, error) {
    args := m.Called(username)
    return args.Bool(0), args.Error(1)
}

func (m *MockUserRepository) EmailExists(email string) (bool, error) {
    args := m.Called(email)
    return args.Bool(0), args.Error(1)
}

func (m *MockUserRepository) UpdateUserByID(user *entities.User) error {
    args := m.Called(user)
    return args.Error(0)
}

func (m *MockUserRepository) CreateStaff(user *entities.User, staff *entities.Staff) (*entities.Staff, error) {
    args := m.Called(user, staff)
    if args.Get(0) == nil {
        return nil, args.Error(1)
    }
    return args.Get(0).(*entities.Staff), args.Error(1)
}

func (m *MockUserRepository) GetStaffByID(id string) (*entities.Staff, error) {
    args := m.Called(id)
    if args.Get(0) == nil {
        return nil, args.Error(1)
    }
    return args.Get(0).(*entities.Staff), args.Error(1)
}

func (m *MockUserRepository) GetStaffByUserID(userID string) (*entities.Staff, error) {
    args := m.Called(userID)
    if args.Get(0) == nil {
        return nil, args.Error(1)
    }
    return args.Get(0).(*entities.Staff), args.Error(1)
}

func (m *MockUserRepository) CreateStaffFile(staffFile *entities.StaffsFiles) (*entities.StaffsFiles, error) {
    args := m.Called(staffFile)
    if args.Get(0) == nil {
        return nil, args.Error(1)
    }
    return args.Get(0).(*entities.StaffsFiles), args.Error(1)
}

func (m *MockUserRepository) GetStaffFileByID(id string) (*entities.StaffsFiles, error) {
    args := m.Called(id)
    if args.Get(0) == nil {
        return nil, args.Error(1)
    }
    return args.Get(0).(*entities.StaffsFiles), args.Error(1)
}

func (m *MockUserRepository) GetAllUsers() ([]*entities.User, error) {
    args := m.Called()
    if args.Get(0) == nil {
        return nil, args.Error(1)
    }
    return args.Get(0).([]*entities.User), args.Error(1)
}

func (m *MockUserRepository) CreateOTP(otp *entities.OTP) error {
    args := m.Called(otp)
    return args.Error(0)
}

func (m *MockUserRepository) GetOTPByUserID(userID string) (*entities.OTP, error) {
    args := m.Called(userID)
    if args.Get(0) == nil {
        return nil, args.Error(1)
    }
    return args.Get(0).(*entities.OTP), args.Error(1)
}

func (m *MockUserRepository) DeleteOTP(userID string) error {
    args := m.Called(userID)
    return args.Error(0)
}

func (m *MockUserRepository) StoreResetToken(temptoken *entities.TempToken) error {
    args := m.Called(temptoken)
    return args.Error(0)
}

func (m *MockUserRepository) GetResetToken(userID string) (string, error) {
    args := m.Called(userID)
    return args.String(0), args.Error(1)
}

func (m *MockUserRepository) DeleteResetToken(userID string) error {
    args := m.Called(userID)
    return args.Error(0)
}