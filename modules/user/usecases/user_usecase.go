package usecases

import (
	"errors"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"time"

	"github.com/aikidoaikido115/New-Acis-BE/configs"
	"github.com/aikidoaikido115/New-Acis-BE/modules/entities"
	"github.com/aikidoaikido115/New-Acis-BE/modules/user/repositories"
	"github.com/aikidoaikido115/New-Acis-BE/pkg/utils"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"golang.org/x/text/unicode/norm"
)

type UserUsecase interface {
	Register(user *entities.User, roleName string) (*entities.User, error)
	Login(username, password string) (string, *entities.User, error)
	ResetPassword(userID, oldPassword, newPassword string) error

	GetUserByID(id string) (*entities.User, error)
	GetAllUsers() ([]*entities.User, error)
	UpdateUserByID(id string, user *entities.User, file multipart.File) (*entities.User, error)

	ForgotPassword(email string) error
	VerifyOTP(email, otpCode string) error
	ChangePassword(email, newPassword string) error

	CreateStaffFile(userID string, files []*multipart.FileHeader) ([]*entities.StaffsFiles, error)
}

type UserUseCaseImpl struct {
	userrepo  repositories.UserRepository
	jwtSecret string
	supa      configs.Supabase
	mail      configs.Mail
}

func NewUserUseCase(
	userrepo repositories.UserRepository,
	jwt configs.JWT,
	supa configs.Supabase,
	mail configs.Mail) UserUsecase {

	return &UserUseCaseImpl{
		userrepo:  userrepo,
		jwtSecret: jwt.Secret,
		supa:      supa,
		mail:      mail,
	}
}

func (u *UserUseCaseImpl) Register(user *entities.User, roleName string) (*entities.User, error) {
	normalizedEmail, err := utils.NormalizeEmail(user.Email)
	if err != nil {
		return nil, errors.New("invalid email format")
	}

	role, err := u.userrepo.GetRoleByName(roleName)
	if err != nil {
		return nil, errors.New("role not found")
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
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, errors.New("failed to hash password")
	}
	user.Password = string(hashedPassword)
	createdUser, err := u.userrepo.CreateUser(user)
	if err != nil {
		return nil, errors.New("failed to create user: " + err.Error())
	}

	if role.Name == "Medical Staff" || role.Name == "Kitchen Staff" {
		staff := &entities.Staff{
			ID:     uuid.New().String(),
			UserID: createdUser.ID}
		_, err := u.userrepo.CreateStaff(createdUser, staff)
		if err != nil {
			return nil, errors.New("failed to create staff: " + err.Error())
		}
	}

	return createdUser, nil
}

func (u *UserUseCaseImpl) Login(username, password string) (string, *entities.User, error) {
	user, err := u.userrepo.GetUserByUsername(username)
	if err != nil {
		return "", nil, errors.New("invalid username")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return "", nil, errors.New("invalid password")
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.ID,
		"exp":     time.Now().Add(time.Minute * 30).Unix(), // หมดอายุใน 30 นาที
		"iat":     time.Now().Unix(),                       // เวลาที่ออก
		"jti":     uuid.New().String(),                     // ให้ token นี้ unique
	})

	tokenString, err := token.SignedString([]byte(u.jwtSecret))
	if err != nil {
		return "", nil, errors.New("failed to generate token")
	}

	return tokenString, user, nil
}

func (u *UserUseCaseImpl) ResetPassword(userID, oldPassword, newPassword string) error {
	user, err := u.userrepo.GetUserByID(userID)
	if err != nil {
		return errors.New("user invalid")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(oldPassword)); err != nil {
		return errors.New("old password invalid")
	}

	hashedNewPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return errors.New("failed to hash new password")
	}
	user.Password = string(hashedNewPassword)

	if err := u.userrepo.UpdateUserByID(user); err != nil {
		return errors.New("failed to update password")
	}

	return nil
}

func (u *UserUseCaseImpl) GetUserByID(id string) (*entities.User, error) {
	user, err := u.userrepo.GetUserByID(id)
	if err != nil {
		return nil, errors.New("user not found")
	}
	return user, nil
}

func (u *UserUseCaseImpl) GetAllUsers() ([]*entities.User, error) {
	users, err := u.userrepo.GetAllUsers()
	if err != nil {
		return nil, errors.New("failed to retrieve all users")
	}
	return users, nil
}

func (u *UserUseCaseImpl) UpdateUserByID(id string, user *entities.User, file multipart.File) (*entities.User, error) {
	existingUser, err := u.userrepo.GetUserByID(id)
	if err != nil {
		return nil, errors.New("user not found")
	}

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

			if existingUser.NumberOfUsernames >= 1 {
				return nil, errors.New("username can change only once")
			}

			existingUser.NumberOfUsernames++ //นับจำนวนครั้งที่เปลี่ยน username
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

	if user.Gender != "" {
		existingUser.Gender = user.Gender
	}

	if file != nil {
		fileExtension, err := utils.DetectFileType(file)
		if err != nil {
			return nil, errors.New("invalid file: " + err.Error())
		}

		// Reset file pointer to beginning after DetectFileType
		file.Seek(0, io.SeekStart)

		fileName := uuid.New().String() + fileExtension

		profileURL, err := utils.UploadFile2Supa(file, fileName, "profiles/", u.supa)
		if err != nil {
			return nil, errors.New("failed to upload profile image: " + err.Error())
		}

		existingUser.ProfileImage = profileURL
	}

	if err := u.userrepo.UpdateUserByID(existingUser); err != nil {
		return nil, errors.New("failed to update user")
	}

	return existingUser, nil
}

func (u *UserUseCaseImpl) ForgotPassword(email string) error {
	user, err := u.userrepo.GetUserByEmail(email)
	if err != nil {
		return errors.New("user not found")
	}
	userID := user.ID
	otpCode, err := utils.GenerateRandomOTP(6)
	if err != nil {
		return errors.New("failed to generate OTP")
	}

	expiresAt := time.Now().Add(5 * time.Minute)
	otp, err := u.userrepo.GetOTPByUserID(userID)

	if err == nil && otp != nil {
		if err := u.userrepo.DeleteOTP(userID); err != nil {
			return errors.New("failed to delete existing OTP")
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

	if err := utils.SendMail(templatePath, user, otpCode, u.mail); err != nil {
		return errors.New("failed to send OTP email: " + err.Error())
	}

	return nil
}

func (u *UserUseCaseImpl) VerifyOTP(email, otpCode string) error {
	user, err := u.userrepo.GetUserByEmail(email)
	if err != nil {
		return errors.New("user not found")
	}

	userID := user.ID
	otp, err := u.userrepo.GetOTPByUserID(userID)
	if err != nil || otp == nil {
		return errors.New("OTP not found for user")
	}

	if time.Now().After(otp.ExpiresAt) {
		return errors.New("OTP has expired")
	}

	if otp.OTP != otpCode {
		return errors.New("invalid OTP code")
	}

	if err := u.userrepo.DeleteOTP(userID); err != nil {
		return errors.New("failed to delete existing OTP")
	}

	tempToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.ID,
		"exp":     time.Now().Add(time.Minute * 5).Unix(), // หมดอายุใน 5 นาที
		"iat":     time.Now().Unix(),                      // เวลาที่ออก
		"jti":     uuid.New().String(),                    // ให้ token นี้ unique
	})

	tempTokenString, err := tempToken.SignedString([]byte(u.jwtSecret))
	if err != nil {
		return errors.New("failed to generate token")
	}

	tempTokenTable := &entities.TempToken{
		UserID: userID,
		Token:  tempTokenString,
	}

	if err := u.userrepo.StoreResetToken(tempTokenTable); err != nil {
		return errors.New("failed to store resetToken")
	}

	return nil
}

func (u *UserUseCaseImpl) ChangePassword(email, newPassword string) error {
	user, err := u.userrepo.GetUserByEmail(email)
	if err != nil {
		return errors.New("user not found")
	}

	tokenString, err := u.userrepo.GetResetToken(user.ID)
	if err != nil {
		return errors.New("reset token not found")
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

	hashedNewPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return errors.New("failed to hash new password")
	}

	user.Password = string(hashedNewPassword)
	if err := u.userrepo.UpdateUserByID(user); err != nil {
		return errors.New("failed to update password")
	}

	if err := u.userrepo.DeleteResetToken(user.ID); err != nil {
		return errors.New("failed to delete reset token")
	}

	return nil
}

func (u *UserUseCaseImpl) CreateStaffFile(userID string, files []*multipart.FileHeader) ([]*entities.StaffsFiles, error) {
	existingUser, err := u.userrepo.GetUserByID(userID)
	if err != nil {
		return nil, errors.New("user not found")
	}
	if existingUser.Role.Name != "Medical Staff" && existingUser.Role.Name != "Kitchen Staff" {
		return nil, errors.New("user is not staff")
	}

	staff, err := u.userrepo.GetStaffByUserID(existingUser.ID)
	if err != nil {
		return nil, errors.New("staff not found")
	}

	if len(files) == 0 {
		return nil, errors.New("no files provided")
	}

	// สร้าง slice สำหรับเก็บผลลัพธ์
	createdFiles := make([]*entities.StaffsFiles, 0, len(files))

	// loop ผ่านแต่ละไฟล์และบันทึกแยกกัน
	for _, fileHeader := range files {
		// เปิดไฟล์
		file, err := fileHeader.Open()
		if err != nil {
			return nil, errors.New("failed to open file: " + err.Error())
		}
		defer file.Close()

		// ตรวจสอบประเภทไฟล์
		fileExtension, err := utils.DetectFileType(file)
		if err != nil {
			return nil, errors.New("invalid file: " + err.Error())
		}

		// Reset file pointer to beginning after DetectFileType
		file.Seek(0, io.SeekStart)

		// หาขนาดไฟล์โดย seek ไปท้ายไฟล์
		fileSize, err := file.Seek(0, io.SeekEnd)
		if err != nil {
			return nil, errors.New("failed to get file size: " + err.Error())
		}

		// Reset file pointer กลับไปจุดเริ่มต้นก่อนอัพโหลด
		file.Seek(0, io.SeekStart)

		fileName := uuid.New().String() + fileExtension

		// อัพโหลดไฟล์ไปยัง Supabase
		staffFileURL, err := utils.UploadFile2Supa(file, fileName, "staff_file/", u.supa)
		if err != nil {
			return nil, errors.New("failed to upload staff file: " + err.Error())
		}

		// สร้าง entity สำหรับแต่ละไฟล์
		staffFile := &entities.StaffsFiles{
			ID:       uuid.New().String(),
			StaffID:  staff.ID,
			File:      staffFileURL,
			FileName: fileName,
			FileType: fileExtension,
			FileSize: fileSize,
		}

		// เรียก repository เพื่อบันทึกแต่ละไฟล์
		createdStaffFile, err := u.userrepo.CreateStaffFile(staffFile)
		if err != nil {
			return nil, errors.New("failed to create staff file: " + err.Error())
		}

		createdFiles = append(createdFiles, createdStaffFile)
	}

	return createdFiles, nil
}
