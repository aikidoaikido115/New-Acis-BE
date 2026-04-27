package usecases

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/aikidoaikido115/New-Acis-BE/configs"
	audit_constants "github.com/aikidoaikido115/New-Acis-BE/modules/audit_logs/constants"
	audit_repo "github.com/aikidoaikido115/New-Acis-BE/modules/audit_logs/repositories"
	"github.com/aikidoaikido115/New-Acis-BE/modules/entities"
	user_constants "github.com/aikidoaikido115/New-Acis-BE/modules/user/constants"
	"github.com/aikidoaikido115/New-Acis-BE/modules/user/repositories"
	"github.com/aikidoaikido115/New-Acis-BE/pkg/utils"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"golang.org/x/text/unicode/norm"
)

var (
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrAccountNotApproved = errors.New("account is pending approval")
	ErrAdminOnly          = errors.New("admin only")
	ErrTargetUserNotFound = errors.New("target user not found")
)

type UserUsecase interface {
	Register(user *entities.User, roleName string, file multipart.File) (*entities.User, error)
	Login(username, email, password string, remember bool) (string, *entities.User, error)
	ResetPassword(userID, oldPassword, newPassword string) error

	GetUserByID(id string) (*entities.User, error)
	GetAllUsers(userID string) ([]*entities.User, error)
	GetUsersByFirstAndLastName(firstName string, lastName string) ([]*entities.User, error)
	UpdateUserByID(id string, user *entities.User, file multipart.File) (*entities.User, error)
	UpdateUserApprovalByID(targetUserID string, isApprove bool, adminUserID string) (*entities.User, error)
	UpdateStaffRoleByID(staffID string, roleName string, userID string) (*entities.User, error)
	DeleteStaffByID(staffID string, userID string) error

	ForgotPassword(email string) error
	VerifyOTP(email, otpCode string) error
	ChangePassword(email, newPassword string) error

	CreateStaffFile(userID string, files []*multipart.FileHeader) ([]*entities.StaffsFiles, error)
}

type UserUseCaseImpl struct {
	userrepo     repositories.UserRepository
	auditlogrepo audit_repo.AuditLogRepository
	jwtSecret    string
	supa         configs.Supabase
	mail         configs.Mail
}

func NewUserUseCase(
	userrepo repositories.UserRepository,
	auditlogrepo audit_repo.AuditLogRepository,
	jwt configs.JWT,
	supa configs.Supabase,
	mail configs.Mail) UserUsecase {

	return &UserUseCaseImpl{
		userrepo:     userrepo,
		auditlogrepo: auditlogrepo,
		jwtSecret:    jwt.Secret,
		supa:         supa,
		mail:         mail,
	}
}

func (u *UserUseCaseImpl) Register(user *entities.User, roleName string, file multipart.File) (*entities.User, error) {
	normalizedEmail, err := utils.NormalizeEmail(user.Email)
	if err != nil {
		return nil, errors.New("invalid email format: " + err.Error())
	}

	role, err := u.userrepo.GetRoleByName(roleName)
	if err != nil {
		return nil, errors.New("role not found: " + err.Error())
	}

	user.Email = normalizedEmail
	user.Username = norm.NFC.String(user.Username)

	usernameExists, err := u.userrepo.UsernameExists(user.Username)
	if err != nil {
		return nil, errors.New("failed to check username availability: " + err.Error())
	}
	if usernameExists {
		return nil, errors.New("username already exists")
	}

	emailExists, err := u.userrepo.EmailExists(user.Email)
	if err != nil {
		return nil, errors.New("failed to check email availability: " + err.Error())
	}
	if emailExists {
		return nil, errors.New("email already exists")
	}

	user.ID = uuid.New().String()
	user.RoleID = role.ID
	user.IsApprove = false

	// Upload profile image if provided
	if file != nil {
		fileExtension, err := utils.DetectFileType(file)
		if err != nil {
			return nil, errors.New("invalid file: " + err.Error())
		}

		// Reset file pointer to beginning after DetectFileType
		if _, err = file.Seek(0, io.SeekStart); err != nil {
			return nil, errors.New("failed to reset file pointer: " + err.Error())
		}

		fileName := uuid.New().String() + fileExtension

		profileURL, err := utils.UploadFile2Supa(file, fileName, "profiles/", u.supa)
		if err != nil {
			return nil, errors.New("failed to upload profile image: " + err.Error())
		}

		user.ProfileImage = profileURL
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, errors.New("failed to hash password: " + err.Error())
	}
	user.Password = string(hashedPassword)
	createdUser, err := u.userrepo.CreateUser(user)
	if err != nil {
		return nil, errors.New("failed to create user: " + err.Error())
	}

	if role.Name == user_constants.RoleMedicalStaff || role.Name == user_constants.RoleKitchenStaff {
		staff := &entities.Staff{
			ID:     uuid.New().String(),
			UserID: createdUser.ID}
		_, err := u.userrepo.CreateStaff(createdUser, staff)
		if err != nil {
			return nil, errors.New("failed to create staff: " + err.Error())
		}
	}

	newUserData, _ := json.Marshal(map[string]interface{}{
		"username":   createdUser.Username,
		"email":      createdUser.Email,
		"role_id":    createdUser.RoleID,
		"first_name": createdUser.FirstName,
		"last_name":  createdUser.LastName,
	})

	auditLog := &entities.AuditLogs{
		ID:        uuid.New().String(),
		TableName: "users",
		RecordID:  createdUser.ID,
		UserID:    createdUser.ID,
		Action:    audit_constants.AuditActionInsert,
		OldValue:  "",
		NewValue:  string(newUserData),
	}

	_, err = u.auditlogrepo.CreateAuditLog(auditLog)
	if err != nil {
		log.Printf("[ERROR] Failed to create audit log for new user %s: %v", createdUser.ID, err)
	}

	return createdUser, nil
}

func (u *UserUseCaseImpl) Login(username, email, password string, remember bool) (string, *entities.User, error) {
	if username != "" && email != "" {
		return "", nil, errors.New("please provide either username or email, not both")
	}

	if username == "" && email == "" {
		return "", nil, errors.New("username or email is required")
	}

	var user *entities.User
	var err error

	if username != "" {
		user, err = u.userrepo.GetUserByUsername(username)
	} else {
		user, err = u.userrepo.GetUserByEmail(email)
	}

	if err != nil {
		return "", nil, ErrInvalidCredentials
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return "", nil, ErrInvalidCredentials
	}

	if !user.IsApprove {
		return "", nil, ErrAccountNotApproved
	}

	expiryDuration := time.Minute * 30
	if remember {
		expiryDuration = time.Hour * 24 * 2
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.ID,
		"exp":     time.Now().Add(expiryDuration).Unix(),
		"iat":     time.Now().Unix(),   // เวลาที่ออก
		"jti":     uuid.New().String(), // ให้ token นี้ unique
	})

	tokenString, err := token.SignedString([]byte(u.jwtSecret))
	if err != nil {
		return "", nil, errors.New("failed to generate token: " + err.Error())
	}

	return tokenString, user, nil
}

func (u *UserUseCaseImpl) UpdateUserApprovalByID(targetUserID string, isApprove bool, adminUserID string) (*entities.User, error) {
	if err := u.ensureAdmin(adminUserID); err != nil {
		return nil, err
	}

	targetUserID = strings.TrimSpace(targetUserID)
	if targetUserID == "" {
		return nil, errors.New("user id is required")
	}

	user, err := u.userrepo.GetUserByID(targetUserID)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrTargetUserNotFound, err)
	}

	oldUserData, _ := json.Marshal(map[string]interface{}{
		"user_id":    user.ID,
		"username":   user.Username,
		"role_id":    user.RoleID,
		"role_name":  user.Role.Name,
		"is_approve": user.IsApprove,
	})

	if err := u.userrepo.UpdateUserApprovalByID(user.ID, isApprove); err != nil {
		return nil, errors.New("failed to update user approval: " + err.Error())
	}

	updatedUser, err := u.userrepo.GetUserByID(user.ID)
	if err != nil {
		return nil, errors.New("failed to get updated user: " + err.Error())
	}

	newUserData, _ := json.Marshal(map[string]interface{}{
		"user_id":    updatedUser.ID,
		"username":   updatedUser.Username,
		"role_id":    updatedUser.RoleID,
		"role_name":  updatedUser.Role.Name,
		"is_approve": updatedUser.IsApprove,
	})

	auditLog := &entities.AuditLogs{
		ID:        uuid.New().String(),
		TableName: "users",
		RecordID:  updatedUser.ID,
		UserID:    adminUserID,
		Action:    audit_constants.AuditActionUpdate,
		OldValue:  string(oldUserData),
		NewValue:  string(newUserData),
	}

	if _, err := u.auditlogrepo.CreateAuditLog(auditLog); err != nil {
		log.Printf("[ERROR] Failed to create audit log for approval update %s: %v", updatedUser.ID, err)
	}

	return updatedUser, nil
}

func (u *UserUseCaseImpl) ResetPassword(userID, oldPassword, newPassword string) error {
	user, err := u.userrepo.GetUserByID(userID)
	if err != nil {
		return errors.New("user invalid: " + err.Error())
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(oldPassword)); err != nil {
		return errors.New("old password invalid")
	}

	// ตรวจสอบว่ารหัสผ่านใหม่ไม่ซ้ำกับรหัสผ่านเดิม
	if oldPassword == newPassword {
		return errors.New("new password cannot be the same as current password")
	}

	hashedNewPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return errors.New("failed to hash new password: " + err.Error())
	}
	user.Password = string(hashedNewPassword)

	if err := u.userrepo.UpdateUserByID(user); err != nil {
		return errors.New("failed to update password: " + err.Error())
	}

	auditLog := &entities.AuditLogs{
		ID:        uuid.New().String(),
		TableName: "users",
		RecordID:  user.ID,
		UserID:    userID,
		Action:    audit_constants.AuditActionUpdate,
		OldValue:  audit_constants.AuditOldNewValuePassword,
		NewValue:  audit_constants.AuditOldNewValuePassword,
	}

	_, err = u.auditlogrepo.CreateAuditLog(auditLog)
	if err != nil {
		log.Printf("[ERROR] Failed to create audit log for password reset %s: %v", userID, err)
	}

	return nil
}

func (u *UserUseCaseImpl) GetUserByID(id string) (*entities.User, error) {
	user, err := u.userrepo.GetUserByID(id)
	if err != nil {
		return nil, errors.New("user not found: " + err.Error())
	}
	return user, nil
}

func (u *UserUseCaseImpl) GetAllUsers(userID string) ([]*entities.User, error) {
	if err := u.ensureAdmin(userID); err != nil {
		return nil, err
	}

	users, err := u.userrepo.GetAllUsers()
	if err != nil {
		return nil, errors.New("failed to retrieve all users: " + err.Error())
	}
	return users, nil
}

func (u *UserUseCaseImpl) GetUsersByFirstAndLastName(firstName string, lastName string) ([]*entities.User, error) {
	users, err := u.userrepo.GetUsersByFirstAndLastName(firstName, lastName)
	if err != nil {
		return nil, errors.New("failed to get users by first and last name: " + err.Error())
	}
	return users, nil
}

func (u *UserUseCaseImpl) UpdateStaffRoleByID(staffID string, roleName string, userID string) (*entities.User, error) {
	if err := u.ensureAdmin(userID); err != nil {
		return nil, err
	}

	staffID = strings.TrimSpace(staffID)
	if staffID == "" {
		return nil, errors.New("staff id is required")
	}

	roleName = strings.TrimSpace(roleName)
	if roleName != user_constants.RoleMedicalStaff && roleName != user_constants.RoleKitchenStaff && roleName != user_constants.RoleSuperUser {
		return nil, errors.New("role_name must be Medical Staff, Kitchen Staff, or Super User")
	}

	staff, err := u.userrepo.GetStaffByID(staffID)
	if err != nil {
		return nil, errors.New("staff not found: " + err.Error())
	}

	currentRoleName := staff.User.Role.Name
	if currentRoleName == "" {
		currentRole, err := u.userrepo.GetRoleByID(staff.User.RoleID)
		if err != nil {
			return nil, errors.New("failed to get current role: " + err.Error())
		}
		currentRoleName = currentRole.Name
	}

	targetRole, err := u.userrepo.GetRoleByName(roleName)
	if err != nil {
		return nil, errors.New("role not found: " + err.Error())
	}

	oldUserData, _ := json.Marshal(map[string]interface{}{
		"user_id":    staff.User.ID,
		"staff_id":   staff.ID,
		"username":   staff.User.Username,
		"role_id":    staff.User.RoleID,
		"role_name":  currentRoleName,
		"first_name": staff.User.FirstName,
		"last_name":  staff.User.LastName,
	})

	staff.User.RoleID = targetRole.ID
	if err := u.userrepo.UpdateUserByID(&staff.User); err != nil {
		return nil, errors.New("failed to update staff role: " + err.Error())
	}

	updatedUser, err := u.userrepo.GetUserByID(staff.User.ID)
	if err != nil {
		return nil, errors.New("failed to get updated user: " + err.Error())
	}

	newUserData, _ := json.Marshal(map[string]interface{}{
		"user_id":    updatedUser.ID,
		"staff_id":   staff.ID,
		"username":   updatedUser.Username,
		"role_id":    updatedUser.RoleID,
		"role_name":  updatedUser.Role.Name,
		"first_name": updatedUser.FirstName,
		"last_name":  updatedUser.LastName,
	})

	auditLog := &entities.AuditLogs{
		ID:        uuid.New().String(),
		TableName: "users",
		RecordID:  updatedUser.ID,
		UserID:    userID,
		Action:    audit_constants.AuditActionUpdate,
		OldValue:  string(oldUserData),
		NewValue:  string(newUserData),
	}

	if _, err := u.auditlogrepo.CreateAuditLog(auditLog); err != nil {
		log.Printf("[ERROR] Failed to create audit log for staff role update %s: %v", updatedUser.ID, err)
	}

	return updatedUser, nil
}

func (u *UserUseCaseImpl) DeleteStaffByID(staffID string, userID string) error {
	if err := u.ensureAdmin(userID); err != nil {
		return err
	}

	staffID = strings.TrimSpace(staffID)
	if staffID == "" {
		return errors.New("staff id is required")
	}

	staff, err := u.userrepo.GetStaffByID(staffID)
	if err != nil {
		return errors.New("staff not found: " + err.Error())
	}

	oldStaffData, _ := json.Marshal(staff)

	if err := u.userrepo.DeleteStaffAndUserByStaffID(staffID); err != nil {
		return errors.New("failed to delete staff and user: " + err.Error())
	}

	auditLog := &entities.AuditLogs{
		ID:        uuid.New().String(),
		TableName: "users",
		RecordID:  staff.User.ID,
		UserID:    userID,
		Action:    audit_constants.AuditActionDelete,
		OldValue:  string(oldStaffData),
		NewValue:  "",
	}

	if _, err := u.auditlogrepo.CreateAuditLog(auditLog); err != nil {
		log.Printf("[ERROR] Failed to create audit log for staff deletion %s: %v", staffID, err)
	}

	return nil
}

func (u *UserUseCaseImpl) ensureAdmin(userID string) error {
	userID = strings.TrimSpace(userID)
	if userID == "" {
		return errors.New("user id is required")
	}

	user, err := u.userrepo.GetUserByID(userID)
	if err != nil {
		return errors.New("failed to get user: " + err.Error())
	}

	if user.Role.Name != user_constants.RoleAdmin {
		return ErrAdminOnly
	}

	return nil
}

func (u *UserUseCaseImpl) UpdateUserByID(id string, user *entities.User, file multipart.File) (*entities.User, error) {
	existingUser, err := u.userrepo.GetUserByID(id)
	if err != nil {
		return nil, errors.New("user not found: " + err.Error())
	}

	oldUserData, _ := json.Marshal(map[string]interface{}{
		"username":      existingUser.Username,
		"first_name":    existingUser.FirstName,
		"last_name":     existingUser.LastName,
		"nickname":      existingUser.Nickname,
		"gender":        existingUser.Gender,
		"phone":         existingUser.Phone,
		"profile_image": existingUser.ProfileImage,
	})

	// อัพเดต เท่าที่ อนุญาตให้อัพเดต
	// อัพเดต username เฉพาะเมื่อมีการส่งมา (ไม่ใช่ค่าว่าง)
	if user.Username != "" {
		if existingUser.Username != user.Username {

			usernameExists, err := u.userrepo.UsernameExists(user.Username)
			if err != nil {
				return nil, errors.New("failed to check username availability: " + err.Error())
			}
			if usernameExists {
				return nil, errors.New("username already taken")
			}
			existingUser.Username = user.Username
		}
	}

	if user.FirstName != "" {
		existingUser.FirstName = user.FirstName
	}

	if user.LastName != "" {
		existingUser.LastName = user.LastName
	}

	if user.Nickname != "" {
		existingUser.Nickname = user.Nickname
	}
	
    if user.Phone != "" {
		existingUser.Phone = user.Phone
	}

	if user.Gender != "" {
		existingUser.Gender = user.Gender
	}

	if file != nil {
		fileExtension, err := utils.DetectFileType(file)
		if err != nil {
			return nil, errors.New("invalid file: " + err.Error())
		}

		// Reset file pointer to beginning after DetectFileType
		if _, err = file.Seek(0, io.SeekStart); err != nil {
			return nil, errors.New("failed to reset file pointer: " + err.Error())
		}

		fileName := uuid.New().String() + fileExtension

		profileURL, err := utils.UploadFile2Supa(file, fileName, "profiles/", u.supa)
		if err != nil {
			return nil, errors.New("failed to upload profile image: " + err.Error())
		}

		existingUser.ProfileImage = profileURL
	}

	if err := u.userrepo.UpdateUserByID(existingUser); err != nil {
		return nil, errors.New("failed to update user: " + err.Error())
	}

	newUserData, _ := json.Marshal(map[string]interface{}{
		"username":      existingUser.Username,
		"first_name":    existingUser.FirstName,
		"last_name":     existingUser.LastName,
		"nickname":      existingUser.Nickname,
		"gender":        existingUser.Gender,
		"phone":         existingUser.Phone,
		"profile_image": existingUser.ProfileImage,
	})

	auditLog := &entities.AuditLogs{
		ID:        uuid.New().String(),
		TableName: "users",
		RecordID:  existingUser.ID,
		UserID:    id,
		Action:    audit_constants.AuditActionUpdate,
		OldValue:  string(oldUserData),
		NewValue:  string(newUserData),
	}

	// สร้าง audit log (ไม่ return error เพื่อไม่ให้กระทบกับการอัปเดต user)
	_, err = u.auditlogrepo.CreateAuditLog(auditLog)
	if err != nil {
		log.Printf("[ERROR] Failed to create audit log for user %s: %v", id, err)
	}

	return existingUser, nil
}

func (u *UserUseCaseImpl) ForgotPassword(email string) error {
	user, err := u.userrepo.GetUserByEmail(email)
	if err != nil {
		return errors.New("user not found: " + err.Error())
	}
	userID := user.ID
	otpCode, err := utils.GenerateRandomOTP(6)
	if err != nil {
		return errors.New("failed to generate OTP: " + err.Error())
	}

	expiresAt := time.Now().Add(5 * time.Minute)
	otp, err := u.userrepo.GetOTPByUserID(userID)

	if err == nil && otp != nil {
		if err := u.userrepo.DeleteOTP(userID); err != nil {
			return errors.New("failed to delete existing OTP: " + err.Error())
		}
	}

	newOTP := &entities.OTP{
		UserID:    userID,
		OTP:       otpCode,
		ExpiresAt: expiresAt,
	}

	if err := u.userrepo.CreateOTP(newOTP); err != nil {
		return errors.New("failed to create OTP " + err.Error())
	}

	workingDir, err := os.Getwd()
	if err != nil {
		return errors.New("failed to get working directory: " + err.Error())
	}
	templatePath := filepath.Join(workingDir, "assets", "OTPMail.html")

	// ส่งอีเมลแบบ asynchronous
	go func() {
		if err := utils.SendMail(templatePath, user, otpCode, u.mail); err != nil {
			log.Printf("[ERROR] Failed to send OTP email: %v", err)
		}
	}()

	return nil
}

func (u *UserUseCaseImpl) VerifyOTP(email, otpCode string) error {
	user, err := u.userrepo.GetUserByEmail(email)
	if err != nil {
		return errors.New("user not found: " + err.Error())
	}

	userID := user.ID
	otp, err := u.userrepo.GetOTPByUserID(userID)
	if err != nil || otp == nil {
		if err != nil {
			return errors.New("OTP not found for user: " + err.Error())
		}
		return errors.New("OTP not found for user")
	}

	// ตรวจสอบว่า OTP หมดอายุหรือไม่ และลบออกหากหมดอายุ
	if time.Now().After(otp.ExpiresAt) {
		// ลบ OTP ที่หมดอายุออกจากระบบ
		_ = u.userrepo.DeleteOTP(userID)
		return errors.New("OTP has expired")
	}

	if otp.OTP != otpCode {
		return errors.New("invalid OTP code")
	}

	if err := u.userrepo.DeleteOTP(userID); err != nil {
		return errors.New("failed to delete existing OTP: " + err.Error())
	}

	tempToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.ID,
		"exp":     time.Now().Add(time.Minute * 5).Unix(), // หมดอายุใน 5 นาที
		"iat":     time.Now().Unix(),                      // เวลาที่ออก
		"jti":     uuid.New().String(),                    // ให้ token นี้ unique
	})

	tempTokenString, err := tempToken.SignedString([]byte(u.jwtSecret))
	if err != nil {
		return errors.New("failed to generate token: " + err.Error())
	}

	tempTokenTable := &entities.TempToken{
		UserID: userID,
		Token:  tempTokenString,
	}

	if err := u.userrepo.StoreResetToken(tempTokenTable); err != nil {
		return errors.New("failed to store resetToken: " + err.Error())
	}

	return nil
}

func (u *UserUseCaseImpl) ChangePassword(email, newPassword string) error {
	user, err := u.userrepo.GetUserByEmail(email)
	if err != nil {
		return errors.New("user not found: " + err.Error())
	}

	tokenString, err := u.userrepo.GetResetToken(user.ID)
	if err != nil {
		return errors.New("reset token not found: " + err.Error())
	}

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("invalid token signing method")
		}
		return []byte(u.jwtSecret), nil
	})

	if err != nil || !token.Valid {
		// ลบ token ที่หมดอายุหรือไม่ถูกต้อง เพื่อทำความสะอาดฐานข้อมูล
		_ = u.userrepo.DeleteResetToken(user.ID)
		return errors.New("invalid or expired token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return errors.New("invalid token claims")
	}

	// ตรวจสอบว่า user_id ใน token ตรงกับ user ที่กำลังเปลี่ยนรหัสผ่าน
	if userId, ok := claims["user_id"].(string); !ok || userId != user.ID {
		return errors.New("token does not match user")
	}

	// ตรวจสอบว่ารหัสผ่านใหม่ไม่ซ้ำกับรหัสผ่านเดิม
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(newPassword)); err == nil {
		return errors.New("new password cannot be the same as current password")
	}

	hashedNewPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return errors.New("failed to hash new password")
	}

	user.Password = string(hashedNewPassword)
	if err := u.userrepo.UpdateUserByID(user); err != nil {
		return errors.New("failed to update password: " + err.Error())
	}

	if err := u.userrepo.DeleteResetToken(user.ID); err != nil {
		return errors.New("failed to delete reset token: " + err.Error())
	}

	auditLog := &entities.AuditLogs{
		ID:        uuid.New().String(),
		TableName: "users",
		RecordID:  user.ID,
		UserID:    user.ID,
		Action:    audit_constants.AuditActionUpdate,
		OldValue:  "password_reset_via_forgot_password",
		NewValue:  audit_constants.AuditOldNewValuePassword,
	}

	_, err = u.auditlogrepo.CreateAuditLog(auditLog)
	if err != nil {
		log.Printf("[ERROR] Failed to create audit log for password reset via forgot password %s: %v", user.ID, err)
	}

	return nil
}

func (u *UserUseCaseImpl) CreateStaffFile(userID string, files []*multipart.FileHeader) ([]*entities.StaffsFiles, error) {
	existingUser, err := u.userrepo.GetUserByID(userID)
	if err != nil {
		return nil, errors.New("user not found: " + err.Error())
	}
	if existingUser.Role.Name != user_constants.RoleMedicalStaff && existingUser.Role.Name != user_constants.RoleKitchenStaff && existingUser.Role.Name != user_constants.RoleSuperUser && existingUser.Role.Name != user_constants.RoleAdmin {
		return nil, errors.New("user is not staff")
	}

	staff, err := u.userrepo.GetStaffByUserID(existingUser.ID)
	if err != nil {
		return nil, errors.New("staff not found: " + err.Error())
	}

	if len(files) == 0 {
		return nil, errors.New("no files provided")
	}

	// สร้าง slice สำหรับเก็บผลลัพธ์และ error handling
	createdFiles := make([]*entities.StaffsFiles, len(files))
	var wg sync.WaitGroup
	var mu sync.Mutex
	var firstError error

	// อัปโหลดไฟล์แบบ parallel ด้วย goroutines
	for i, fileHeader := range files {
		wg.Add(1)
		go func(index int, fh *multipart.FileHeader) {
			defer wg.Done()

			// เปิดไฟล์
			file, err := fh.Open()
			if err != nil {
				mu.Lock()
				if firstError == nil {
					firstError = errors.New("failed to open file: " + err.Error())
				}
				mu.Unlock()
				return
			}
			defer file.Close()

			// ตรวจสอบประเภทไฟล์
			fileExtension, err := utils.DetectFileType(file)
			if err != nil {
				mu.Lock()
				if firstError == nil {
					firstError = errors.New("invalid file: " + err.Error())
				}
				mu.Unlock()
				return
			}

			// Reset file pointer to beginning after DetectFileType
			if _, err = file.Seek(0, io.SeekStart); err != nil {
				mu.Lock()
				if firstError == nil {
					firstError = errors.New("failed to reset file pointer: " + err.Error())
				}
				mu.Unlock()
				return
			}

			// หาขนาดไฟล์โดย seek ไปท้ายไฟล์
			fileSize, err := file.Seek(0, io.SeekEnd)
			if err != nil {
				mu.Lock()
				if firstError == nil {
					firstError = errors.New("failed to get file size: " + err.Error())
				}
				mu.Unlock()
				return
			}

			// Reset file pointer กลับไปจุดเริ่มต้นก่อนอัพโหลด
			if _, err = file.Seek(0, io.SeekStart); err != nil {
				mu.Lock()
				if firstError == nil {
					firstError = errors.New("failed to reset file pointer: " + err.Error())
				}
				mu.Unlock()
				return
			}

			fileName := uuid.New().String() + fileExtension

			// อัพโหลดไฟล์ไปยัง Supabase
			staffFileURL, err := utils.UploadFile2Supa(file, fileName, "staff_file/", u.supa)
			if err != nil {
				mu.Lock()
				if firstError == nil {
					firstError = errors.New("failed to upload staff file: " + err.Error())
				}
				mu.Unlock()
				return
			}

			// สร้าง entity สำหรับแต่ละไฟล์
			staffFile := &entities.StaffsFiles{
				ID:       uuid.New().String(),
				StaffID:  staff.ID,
				File:     staffFileURL,
				FileName: fileName,
				FileType: fileExtension,
				FileSize: fileSize,
			}

			// เรียก repository เพื่อบันทึกแต่ละไฟล์
			createdStaffFile, err := u.userrepo.CreateStaffFile(staffFile)
			if err != nil {
				mu.Lock()
				if firstError == nil {
					firstError = errors.New("failed to create staff file: " + err.Error())
				}
				mu.Unlock()
				return
			}

			// เก็บผลลัพธ์ในตำแหน่งที่ถูกต้อง (ไม่ต้อง lock เพราะแต่ละ goroutine เขียนคนละ index)
			createdFiles[index] = createdStaffFile
		}(i, fileHeader)
	}

	// รอให้ทุก goroutine ทำงานเสร็จ
	wg.Wait()

	// ตรวจสอบว่ามี error เกิดขึ้นหรือไม่
	if firstError != nil {
		return nil, firstError
	}

	// กรองเฉพาะไฟล์ที่อัปโหลดสำเร็จ (ไม่ใช่ nil)
	successFiles := make([]*entities.StaffsFiles, 0, len(files))
	for _, file := range createdFiles {
		if file != nil {
			successFiles = append(successFiles, file)
		}
	}

	return successFiles, nil
}
