package usecases

import (
    "mime/multipart"

    "github.com/aikidoaikido115/New-Acis-BE/modules/entities"
    "github.com/stretchr/testify/mock"
)

type MockUserUsecase struct {
    mock.Mock
}

func (m *MockUserUsecase) Register(user *entities.User, roleName string) (*entities.User, error) {
    args := m.Called(user, roleName)
    if args.Get(0) == nil {
        return nil, args.Error(1)
    }
    return args.Get(0).(*entities.User), args.Error(1)
}

func (m *MockUserUsecase) Login(username, password string) (string, *entities.User, error) {
    args := m.Called(username, password)
    return args.String(0), args.Get(1).(*entities.User), args.Error(2)
}

func (m *MockUserUsecase) ResetPassword(userID, oldPassword, newPassword string) error {
    args := m.Called(userID, oldPassword, newPassword)
    return args.Error(0)
}

func (m *MockUserUsecase) GetUserByID(id string) (*entities.User, error) {
    args := m.Called(id)
    if args.Get(0) == nil {
        return nil, args.Error(1)
    }
    return args.Get(0).(*entities.User), args.Error(1)
}

func (m *MockUserUsecase) GetAllUsers() ([]*entities.User, error) {
    args := m.Called()
    if args.Get(0) == nil {
        return nil, args.Error(1)
    }
    return args.Get(0).([]*entities.User), args.Error(1)
}

func (m *MockUserUsecase) UpdateUserByID(id string, user *entities.User, file multipart.File) (*entities.User, error) {
    args := m.Called(id, user, file)
    if args.Get(0) == nil {
        return nil, args.Error(1)
    }
    return args.Get(0).(*entities.User), args.Error(1)
}

func (m *MockUserUsecase) ForgotPassword(email string) error {
    args := m.Called(email)
    return args.Error(0)
}

func (m *MockUserUsecase) VerifyOTP(email, otpCode string) error {
    args := m.Called(email, otpCode)
    return args.Error(0)
}

func (m *MockUserUsecase) ChangePassword(email, newPassword string) error {
    args := m.Called(email, newPassword)
    return args.Error(0)
}

func (m *MockUserUsecase) CreateStaffFile(userID string, files []*multipart.FileHeader) ([]*entities.StaffsFiles, error) {
    args := m.Called(userID, files)
    if args.Get(0) == nil {
        return nil, args.Error(1)
    }
    return args.Get(0).([]*entities.StaffsFiles), args.Error(1)
}