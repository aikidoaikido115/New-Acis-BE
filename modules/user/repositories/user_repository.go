package repositories

import (
	"github.com/aikidoaikido115/New-Acis-BE/modules/entities"

	"gorm.io/gorm"
)

type GormUserRepository struct {
	db *gorm.DB
}

func NewGormUserRepository(db *gorm.DB) *GormUserRepository {
	return &GormUserRepository{
		db: db,
	}
}

type UserRepository interface {
	CreateUser(user *entities.User) (*entities.User, error)
	CreateStaff(user *entities.User, staff *entities.Staff) (*entities.Staff, error)
	CreateStaffFile(staffFile *entities.StaffsFiles) (*entities.StaffsFiles, error)
	// TODO CrateRelative()

	GetUserByEmail(email string) (*entities.User, error)
	GetUserByID(id string) (*entities.User, error)
	GetStaffByID(id string) (*entities.Staff, error)
	GetStaffByUserID(userID string) (*entities.Staff, error)
	GetStaffFileByID(id string) (*entities.StaffsFiles, error)
	GetUserByUsername(username string) (*entities.User, error)
	GetRoleByName(roleName string) (*entities.Role, error)
	UsernameExists(username string) (bool, error)
	EmailExists(email string) (bool, error)
	GetAllUsers() ([]*entities.User, error)
	UpdateUserByID(user *entities.User) error
	CreateOTP(otp *entities.OTP) error
	GetOTPByUserID(userID string) (*entities.OTP, error)
	DeleteOTP(userID string) error	
	StoreResetToken(temptoken *entities.TempToken) error
	GetResetToken(userID string) (string, error)
	DeleteResetToken(userID string) error
}

func (r *GormUserRepository) CreateUser(user *entities.User) (*entities.User, error) {
	err := r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&user).Error; err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	return r.GetUserByID(user.ID)
}

func (r *GormUserRepository) CreateStaff(user *entities.User, staff *entities.Staff) (*entities.Staff, error) {
	err := r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&staff).Error; err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	return r.GetStaffByID(staff.ID)
}

func (r *GormUserRepository) CreateStaffFile(staffFile *entities.StaffsFiles) (*entities.StaffsFiles, error) {
	err := r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&staffFile).Error; err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	return r.GetStaffFileByID(staffFile.ID)
}

func (r *GormUserRepository) GetStaffFileByID(id string) (*entities.StaffsFiles, error) {
	var staffFile entities.StaffsFiles
	err := r.db.First(&staffFile, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &staffFile, nil
}

func (r *GormUserRepository) GetUserByEmail(email string) (*entities.User, error) {
	var user entities.User
	err := r.db.Where("email = ?", email).First(&user).Error
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *GormUserRepository) GetUserByID(id string) (*entities.User, error) {
	var user entities.User
	if err := r.db.Preload("Role").First(&user, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *GormUserRepository) GetStaffByID(id string) (*entities.Staff, error) {
	var staff entities.Staff
	if err := r.db.Preload("User").First(&staff, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &staff, nil
}

func (r *GormUserRepository) GetStaffByUserID(userID string) (*entities.Staff, error) {
	var staff entities.Staff
	if err := r.db.Preload("User").First(&staff, "user_id = ?", userID).Error; err != nil {
		return nil, err
	}
	return &staff, nil
}

func (r *GormUserRepository) GetUserByUsername(username string) (*entities.User, error) {
	var user entities.User
	if err := r.db.First(&user, "username = ?", username).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *GormUserRepository) GetRoleByName(roleName string) (*entities.Role, error) {
	var role entities.Role
	if err := r.db.First(&role, "name = ?", roleName).Error; err != nil {
		return nil, err
	}
	return &role, nil
}

func (r *GormUserRepository) UsernameExists(username string) (bool, error) {
	var count int64
	err := r.db.Model(&entities.User{}).Where("username = ?", username).Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r *GormUserRepository) EmailExists(email string) (bool, error) {
	var count int64
	err := r.db.Model(&entities.User{}).Where("email = ?", email).Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r *GormUserRepository) GetAllUsers() ([]*entities.User, error) {
	var users []*entities.User
	if err := r.db.Find(&users).Error; err != nil {
		return nil, err
	}
	return users, nil
}

func (r *GormUserRepository) UpdateUserByID(user *entities.User) error {
	return r.db.Save(user).Error
}

// OTP methods implementation
func (r *GormUserRepository) CreateOTP(otp *entities.OTP) error {
	return r.db.Create(otp).Error
}

func (r *GormUserRepository) GetOTPByUserID(userID string) (*entities.OTP, error) {
	var otp entities.OTP
	if err := r.db.Where("user_id = ?", userID).First(&otp).Error; err != nil {
		return nil, err
	}
	return &otp, nil
}

func (r *GormUserRepository) DeleteOTP(userID string) error {
	if err := r.db.Delete(&entities.OTP{}, "user_id = ?", userID).Error; err != nil {
		return err
	}

	return nil
}

func (r *GormUserRepository) StoreResetToken(tempToken *entities.TempToken) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		// Delete existing token if it exists
		if err := tx.Where("user_id = ?", tempToken.UserID).Delete(&entities.TempToken{}).Error; err != nil {
			return err
		}

		// Create new token
		tempToken := &entities.TempToken{
			UserID: tempToken.UserID,
			Token:  tempToken.Token,
		}
		return tx.Create(tempToken).Error
	})
}

func (r *GormUserRepository) GetResetToken(userID string) (string, error) {
	var tempToken entities.TempToken
	if err := r.db.Where("user_id = ?", userID).First(&tempToken).Error; err != nil {
		return "", err
	}
	return tempToken.Token, nil
}

func (r *GormUserRepository) DeleteResetToken(userID string) error {
	if err := r.db.Where("user_id = ?", userID).Delete(&entities.TempToken{}).Error; err != nil {
		return err
	}
	return nil
}
