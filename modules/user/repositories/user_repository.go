package repositories

import (
	"strings"
	"time"

	"github.com/aikidoaikido115/New-Acis-BE/modules/entities"
	user_constants "github.com/aikidoaikido115/New-Acis-BE/modules/user/constants"

	"gorm.io/gorm"
)

type AdminRelativeUser struct {
	UserID       string    `json:"user_id"`
	RelativeID   string    `json:"relative_id"`
	Username     string    `json:"username"`
	ResidentName string    `json:"resident_name"`
	CreatedAt    time.Time `json:"created_at"`
}

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
	GetUsersByFirstAndLastName(firstName string, lastName string) ([]*entities.User, error)
	GetStaffByID(id string) (*entities.Staff, error)
	GetStaffByUserID(userID string) (*entities.Staff, error)
	GetStaffFileByID(id string) (*entities.StaffsFiles, error)
	GetUserByUsername(username string) (*entities.User, error)
	GetRoleByName(roleName string) (*entities.Role, error)
	GetRoleByID(roleID string) (*entities.Role, error)
	UsernameExists(username string) (bool, error)
	EmailExists(email string) (bool, error)
	GetAllUsers() ([]*entities.User, error)
	GetRelativeUsersWithResident() ([]AdminRelativeUser, error)
	GetRelativeUserByUserID(userID string) (*AdminRelativeUser, error)
	GetStaffIDMapByUserIDs(userIDs []string) (map[string]string, error)
	UpdateUserByID(user *entities.User) error
	UpdateUserApprovalByID(userID string, isApprove bool) error
	DeleteStaffAndUserByStaffID(staffID string) error
	DeleteRelativeAndUserByUserID(userID string) error
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

func (r *GormUserRepository) GetUsersByFirstAndLastName(firstName string, lastName string) ([]*entities.User, error) {
	var users []*entities.User
	if err := r.db.
		Preload("Role").
		Where("LOWER(TRIM(first_name)) = LOWER(TRIM(?))", firstName).
		Where("LOWER(TRIM(last_name)) = LOWER(TRIM(?))", lastName).
		Find(&users).Error; err != nil {
		return nil, err
	}

	return users, nil
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

func (r *GormUserRepository) GetRoleByID(roleID string) (*entities.Role, error) {
	var role entities.Role
	if err := r.db.First(&role, "id = ?", roleID).Error; err != nil {
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
	if err := r.db.Preload("Role").Find(&users).Error; err != nil {
		return nil, err
	}
	return users, nil
}

func (r *GormUserRepository) GetRelativeUsersWithResident() ([]AdminRelativeUser, error) {
	var rows []AdminRelativeUser

	residentNameExpr := strings.Join([]string{
		"COALESCE(NULLIF(TRIM(CONCAT(COALESCE(residents.first_name, ''), ' ', COALESCE(residents.last_name, ''))), ''),",
		"NULLIF(TRIM(COALESCE(residents.nickname, '')), ''),",
		"users.username)",
	}, " ")

	if err := r.db.
		Table("relatives").
		Select("users.id AS user_id, relatives.id AS relative_id, users.username, "+residentNameExpr+" AS resident_name, users.created_at").
		Joins("JOIN users ON users.id = relatives.user_id").
		Joins("JOIN residents ON residents.id = relatives.resident_id").
		Joins("JOIN roles ON roles.id = users.role_id").
		Where("LOWER(TRIM(roles.name)) = LOWER(TRIM(?))", user_constants.RoleRelative).
		Order("users.created_at DESC").
		Scan(&rows).Error; err != nil {
		return nil, err
	}

	return rows, nil
}

func (r *GormUserRepository) GetRelativeUserByUserID(userID string) (*AdminRelativeUser, error) {
	residentNameExpr := strings.Join([]string{
		"COALESCE(NULLIF(TRIM(CONCAT(COALESCE(residents.first_name, ''), ' ', COALESCE(residents.last_name, ''))), ''),",
		"NULLIF(TRIM(COALESCE(residents.nickname, '')), ''),",
		"users.username)",
	}, " ")

	var row AdminRelativeUser
	err := r.db.
		Table("relatives").
		Select("users.id AS user_id, relatives.id AS relative_id, users.username, "+residentNameExpr+" AS resident_name, users.created_at").
		Joins("JOIN users ON users.id = relatives.user_id").
		Joins("JOIN residents ON residents.id = relatives.resident_id").
		Joins("JOIN roles ON roles.id = users.role_id").
		Where("users.id = ?", userID).
		Where("LOWER(TRIM(roles.name)) = LOWER(TRIM(?))", user_constants.RoleRelative).
		Take(&row).Error
	if err != nil {
		return nil, err
	}

	return &row, nil
}

func (r *GormUserRepository) GetStaffIDMapByUserIDs(userIDs []string) (map[string]string, error) {
	result := map[string]string{}
	if len(userIDs) == 0 {
		return result, nil
	}

	type staffRow struct {
		UserID  string
		StaffID string
	}

	var rows []staffRow
	if err := r.db.
		Table("staffs").
		Select("user_id, id as staff_id").
		Where("user_id IN ?", userIDs).
		Scan(&rows).Error; err != nil {
		return nil, err
	}

	for _, row := range rows {
		result[row.UserID] = row.StaffID
	}

	return result, nil
}

func (r *GormUserRepository) UpdateUserByID(user *entities.User) error {
	return r.db.Save(user).Error
}

func (r *GormUserRepository) UpdateUserApprovalByID(userID string, isApprove bool) error {
	return r.db.Model(&entities.User{}).Where("id = ?", userID).Update("is_approve", isApprove).Error
}

func (r *GormUserRepository) DeleteStaffAndUserByStaffID(staffID string) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		var staff entities.Staff
		if err := tx.First(&staff, "id = ?", staffID).Error; err != nil {
			return err
		}

		if err := tx.Exec("DELETE FROM support_tickets WHERE created_by_user_id = ?", staff.UserID).Error; err != nil {
			return err
		}
		if err := tx.Exec("DELETE FROM warehouse_transactions WHERE operator_user_id = ?", staff.UserID).Error; err != nil {
			return err
		}
		if err := tx.Exec("DELETE FROM otps WHERE user_id = ?", staff.UserID).Error; err != nil {
			return err
		}
		if err := tx.Exec("DELETE FROM temp_tokens WHERE user_id = ?", staff.UserID).Error; err != nil {
			return err
		}

		if err := tx.Exec("DELETE FROM staffs_files WHERE staff_id = ?", staffID).Error; err != nil {
			return err
		}
		if err := tx.Exec("DELETE FROM activities WHERE staff_id = ?", staffID).Error; err != nil {
			return err
		}
		if err := tx.Exec("UPDATE rooms SET staff_id = NULL WHERE staff_id = ?", staffID).Error; err != nil {
			return err
		}
		if err := tx.Exec("UPDATE vital_signs SET created_by_staff_id = NULL WHERE created_by_staff_id = ?", staffID).Error; err != nil {
			return err
		}
		if err := tx.Exec("UPDATE laboratory_values SET created_by_staff_id = NULL WHERE created_by_staff_id = ?", staffID).Error; err != nil {
			return err
		}
		if err := tx.Exec("DELETE FROM nurse_notes WHERE created_by_staff_id = ?", staffID).Error; err != nil {
			return err
		}
		if err := tx.Exec("DELETE FROM wound_care_notes WHERE created_by_staff_id = ?", staffID).Error; err != nil {
			return err
		}
		if err := tx.Exec("DELETE FROM relative_notes WHERE created_by_staff_id = ?", staffID).Error; err != nil {
			return err
		}
		if err := tx.Exec("DELETE FROM doctor_orders WHERE created_by_staff_id = ?", staffID).Error; err != nil {
			return err
		}
		if err := tx.Exec("DELETE FROM drug_plans WHERE given_by_staff_id = ?", staffID).Error; err != nil {
			return err
		}
		if err := tx.Delete(&entities.Staff{}, "id = ?", staffID).Error; err != nil {
			return err
		}
		if err := tx.Delete(&entities.User{}, "id = ?", staff.UserID).Error; err != nil {
			return err
		}

		return nil
	})
}

func (r *GormUserRepository) DeleteRelativeAndUserByUserID(userID string) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		var relative entities.Relative
		if err := tx.First(&relative, "user_id = ?", userID).Error; err != nil {
			return err
		}

		if err := tx.Exec("DELETE FROM relative_magic_link_tokens WHERE relative_id = ?", relative.ID).Error; err != nil {
			return err
		}
		if err := tx.Exec("DELETE FROM daily_updates WHERE relative_id = ?", relative.ID).Error; err != nil {
			return err
		}
		if err := tx.Exec("DELETE FROM support_tickets WHERE created_by_user_id = ?", userID).Error; err != nil {
			return err
		}
		if err := tx.Exec("DELETE FROM otps WHERE user_id = ?", userID).Error; err != nil {
			return err
		}
		if err := tx.Exec("DELETE FROM temp_tokens WHERE user_id = ?", userID).Error; err != nil {
			return err
		}

		if err := tx.Delete(&entities.Relative{}, "id = ?", relative.ID).Error; err != nil {
			return err
		}
		if err := tx.Delete(&entities.User{}, "id = ?", userID).Error; err != nil {
			return err
		}

		return nil
	})
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
		newTempToken := &entities.TempToken{
			UserID: tempToken.UserID,
			Token:  tempToken.Token,
		}
		return tx.Create(newTempToken).Error
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
