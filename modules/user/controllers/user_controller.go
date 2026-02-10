package controllers

import (
	// "fmt"
	// "mime/multipart"
	// "strconv"
	"mime/multipart"

	"github.com/aikidoaikido115/New-Acis-BE/modules/entities"
	"github.com/aikidoaikido115/New-Acis-BE/modules/user/models"
	"github.com/aikidoaikido115/New-Acis-BE/modules/user/usecases"

	"github.com/gofiber/fiber/v2"
)

type UserController struct {
	userusecase usecases.UserUsecase
}

func NewUserController(userusecase usecases.UserUsecase) *UserController {
	return &UserController{
		userusecase: userusecase,
	}
}

// RegisterHandler godoc
// @Summary User Registration
// @Description Register a new user with username, email, password, and role
// @Tags Authentication
// @Accept multipart/form-data
// @Produce json
// @Param username formData string true "Username"
// @Param email formData string true "Email"
// @Param password formData string true "Password"
// @Param first_name formData string false "First Name"
// @Param last_name formData string false "Last Name"
// @Param nickname formData string false "Nickname"
// @Param role_name formData string true "Role Name"
// @Param profile_image formData file false "Profile Image"
// @Success 201 {object} object{status=string,status_code=int,message=string,result=object} "User created successfully"
// @Failure 400 {object} object{status=string,status_code=int,message=string,result=any} "Bad Request - Missing required fields"
// @Failure 500 {object} object{status=string,status_code=int,message=string,result=any} "Internal Server Error"
// @Router /api/auth/register [post]
func (c *UserController) RegisterHandler(ctx *fiber.Ctx) error {
	form, err := ctx.MultipartForm()
	if err != nil {
		return ctx.Status(fiber.ErrBadRequest.Code).JSON(fiber.Map{
			"status":      fiber.ErrBadRequest.Message,
			"status_code": fiber.ErrBadRequest.Code,
			"message":     "Invalid form data: " + err.Error(),
			"result":      nil,
		})
	}

	// Extract values from form
	var username, email, password, firstName, lastName, nickname, roleName string

	if usernames := form.Value["username"]; len(usernames) > 0 {
		username = usernames[0]
	}
	if emails := form.Value["email"]; len(emails) > 0 {
		email = emails[0]
	}
	if passwords := form.Value["password"]; len(passwords) > 0 {
		password = passwords[0]
	}
	if firstNames := form.Value["first_name"]; len(firstNames) > 0 {
		firstName = firstNames[0]
	}
	if lastNames := form.Value["last_name"]; len(lastNames) > 0 {
		lastName = lastNames[0]
	}
	if nicknames := form.Value["nickname"]; len(nicknames) > 0 {
		nickname = nicknames[0]
	}
	if roleNames := form.Value["role_name"]; len(roleNames) > 0 {
		roleName = roleNames[0]
	}

	// Validate required fields
	if username == "" {
		return ctx.Status(fiber.ErrBadRequest.Code).JSON(fiber.Map{
			"status":      fiber.ErrBadRequest.Message,
			"status_code": fiber.ErrBadRequest.Code,
			"message":     "Username is missing",
			"result":      nil,
		})
	}

	if email == "" {
		return ctx.Status(fiber.ErrBadRequest.Code).JSON(fiber.Map{
			"status":      fiber.ErrBadRequest.Message,
			"status_code": fiber.ErrBadRequest.Code,
			"message":     "Email is missing",
			"result":      nil,
		})
	}

	if password == "" {
		return ctx.Status(fiber.ErrBadRequest.Code).JSON(fiber.Map{
			"status":      fiber.ErrBadRequest.Message,
			"status_code": fiber.ErrBadRequest.Code,
			"message":     "Password is missing",
			"result":      nil,
		})
	}

	// Get profile image file (optional)
	files := form.File["profile_image"]
	var file multipart.File

	if len(files) > 0 {
		fileHeader := files[0]
		file, err = fileHeader.Open()
		if err != nil {
			return ctx.Status(fiber.ErrBadRequest.Code).JSON(fiber.Map{
				"status":      fiber.ErrBadRequest.Message,
				"status_code": fiber.ErrBadRequest.Code,
				"message":     "Failed to open uploaded file: " + err.Error(),
				"result":      nil,
			})
		}
		defer file.Close()
	}

	user := &entities.User{
		Username:  username,
		Email:     email,
		Password:  password,
		FirstName: firstName,
		LastName:  lastName,
		Nickname:  nickname,
	}

	data, err := c.userusecase.Register(user, roleName, file)
	if err != nil {
		return ctx.Status(fiber.ErrInternalServerError.Code).JSON(fiber.Map{
			"status":      fiber.ErrInternalServerError.Message,
			"status_code": fiber.ErrInternalServerError.Code,
			"message":     err.Error(),
			"result":      nil,
		})
	}

	return ctx.Status(fiber.StatusCreated).JSON(fiber.Map{
		"status":      "Success",
		"status_code": fiber.StatusCreated,
		"message":     "user created successfully",
		"result":      data,
	})
}

// LoginHandler godoc
// @Summary User Login
// @Description Login with username or email and password to get access token. You must provide EITHER username OR email, not both
// @Tags Authentication
// @Accept json
// @Produce json
// @Param request body object{username=string,email=string,password=string} true "Login credentials (provide EITHER username OR email with password, not both)"
// @Success 200 {object} object{status=string,status_code=int,message=string,result=object{token=string,user_id=string,username=string,email=string,profile_image=string}} "Login successful"
// @Failure 400 {object} object{status=string,status_code=int,message=string,result=any} "Bad Request - Missing credentials, sent both username and email, or missing password"
// @Failure 500 {object} object{status=string,status_code=int,message=string,result=any} "Internal Server Error"
// @Router /api/auth/login [post]
func (c *UserController) LoginHandler(ctx *fiber.Ctx) error {
	var req models.LoginRequest

	if err := ctx.BodyParser(&req); err != nil {
		return ctx.Status(fiber.ErrBadRequest.Code).JSON(fiber.Map{
			"status":      fiber.ErrBadRequest.Message,
			"status_code": fiber.ErrBadRequest.Code,
			"message":     err.Error(),
			"result":      nil,
		})
	}

	if req.Username == "" && req.Email == "" {
		return ctx.Status(fiber.ErrBadRequest.Code).JSON(fiber.Map{
			"status":      fiber.ErrBadRequest.Message,
			"status_code": fiber.ErrBadRequest.Code,
			"message":     "Username or Email is missing",
			"result":      nil,
		})
	}

	if req.Password == "" {
		return ctx.Status(fiber.ErrBadRequest.Code).JSON(fiber.Map{
			"status":      fiber.ErrBadRequest.Message,
			"status_code": fiber.ErrBadRequest.Code,
			"message":     "Password is missing",
			"result":      nil,
		})
	}

	token, user, err := c.userusecase.Login(req.Username, req.Email, req.Password)
	if err != nil {
		return ctx.Status(fiber.ErrInternalServerError.Code).JSON(fiber.Map{
			"status":      fiber.ErrInternalServerError.Message,
			"status_code": fiber.ErrInternalServerError.Code,
			"message":     err.Error(),
			"result":      nil,
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":      "Success",
		"status_code": fiber.StatusOK,
		"message":     "Login successful",
		"result": fiber.Map{
			"token":         token,
			"user_id":       user.ID,
			"username":      user.Username,
			"email":         user.Email,
			"profile_image": user.ProfileImage,
		},
	})
}

// ResetPasswordHandler godoc
// @Summary Reset Password (Authenticated)
// @Description Reset password for authenticated user with old password verification
// @Tags Authentication
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body object{old_password=string,new_password=string} true "Password reset information"
// @Success 200 {object} object{status=string,status_code=int,message=string,result=any} "Password reset successfully"
// @Failure 400 {object} object{status=string,status_code=int,message=string,result=any} "Bad Request - Old password invalid or passwords are the same"
// @Failure 401 {object} object{status=string,status_code=int,message=string,result=any} "Unauthorized - Missing user ID"
// @Failure 404 {object} object{status=string,status_code=int,message=string,result=any} "User not found"
// @Failure 500 {object} object{status=string,status_code=int,message=string,result=any} "Internal Server Error"
// @Router /api/auth/resetpassword [patch]
func (c *UserController) ResetPasswordHandler(ctx *fiber.Ctx) error {
	var req models.ResetPasswordRequest

	if err := ctx.BodyParser(&req); err != nil {
		return ctx.Status(fiber.ErrBadRequest.Code).JSON(fiber.Map{
			"status":      fiber.ErrBadRequest.Message,
			"status_code": fiber.ErrBadRequest.Code,
			"message":     err.Error(),
			"result":      nil,
		})
	}

	if req.OldPassword == "" || req.NewPassword == "" {
		return ctx.Status(fiber.ErrBadRequest.Code).JSON(fiber.Map{
			"status":      fiber.ErrBadRequest.Message,
			"status_code": fiber.ErrBadRequest.Code,
			"message":     "Old password and new password is missing",
			"result":      nil,
		})
	}

	userID, ok := ctx.Locals("user_id").(string)
	if !ok || userID == "" {
		return ctx.Status(fiber.ErrUnauthorized.Code).JSON(fiber.Map{
			"status":      fiber.ErrUnauthorized.Message,
			"status_code": fiber.ErrUnauthorized.Code,
			"message":     "Unauthorized: Missing user ID",
			"result":      nil,
		})
	}

	err := c.userusecase.ResetPassword(userID, req.OldPassword, req.NewPassword)
	if err != nil {
		// Return more specific status codes based on error type
		if err.Error() == "user invalid" {
			return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"status":      fiber.ErrNotFound.Message,
				"status_code": fiber.StatusNotFound,
				"message":     err.Error(),
				"result":      nil,
			})
		} else if err.Error() == "old password invalid" {
			return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"status":      fiber.ErrBadRequest.Message,
				"status_code": fiber.StatusBadRequest,
				"message":     err.Error(),
				"result":      nil,
			})
		} else {
			return ctx.Status(fiber.ErrInternalServerError.Code).JSON(fiber.Map{
				"status":      fiber.ErrInternalServerError.Message,
				"status_code": fiber.ErrInternalServerError.Code,
				"message":     err.Error(),
				"result":      nil,
			})
		}
	}
	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":      "Success",
		"status_code": fiber.StatusOK,
		"message":     "Password reset successfully",
		"result":      nil,
	})
}

// ForgotPasswordHandler godoc
// @Summary Forgot Password - Request OTP
// @Description Send OTP to user's email for password recovery
// @Tags Authentication
// @Accept json
// @Produce json
// @Param request body object{email=string} true "User email"
// @Success 200 {object} object{status=string,status_code=int,message=string} "Sent OTP successfully"
// @Failure 400 {object} object{status=string,status_code=int,message=string,result=any} "Bad Request - Email is missing"
// @Failure 500 {object} object{status=string,status_code=int,message=string,result=any} "Internal Server Error"
// @Router /api/auth/forgotpassword [post]
func (c *UserController) ForgotPasswordHandler(ctx *fiber.Ctx) error {
	var req models.ForgotPasswordRequest

	if err := ctx.BodyParser(&req); err != nil {
		return ctx.Status(fiber.ErrInternalServerError.Code).JSON(fiber.Map{
			"status":      fiber.ErrInternalServerError.Message,
			"status_code": fiber.ErrInternalServerError.Code,
			"message":     err.Error(),
			"result":      nil,
		})
	}

	if req.Email == "" {
		return ctx.Status(fiber.ErrBadRequest.Code).JSON(fiber.Map{
			"status":      fiber.ErrBadRequest.Message,
			"status_code": fiber.ErrBadRequest.Code,
			"message":     "Email is missing",
			"result":      nil,
		})
	}

	err := c.userusecase.ForgotPassword(req.Email)
	if err != nil {
		return ctx.Status(fiber.ErrInternalServerError.Code).JSON(fiber.Map{
			"status":      fiber.ErrInternalServerError.Message,
			"status_code": fiber.ErrInternalServerError.Code,
			"message":     err.Error(),
			"result":      nil,
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":      "Success",
		"status_code": fiber.StatusOK,
		"message":     "Sent OTP successfully",
	})
}

// VerifyOTPHandler godoc
// @Summary Verify OTP
// @Description Verify the OTP sent to user's email for password recovery
// @Tags Authentication
// @Accept json
// @Produce json
// @Param request body object{email=string,otp=string} true "Email and OTP"
// @Success 200 {object} object{status=string,status_code=int,message=string} "OTP is correct"
// @Failure 400 {object} object{status=string,status_code=int,message=string,result=any} "Bad Request - Email or OTP is missing"
// @Failure 500 {object} object{status=string,status_code=int,message=string,result=any} "Internal Server Error"
// @Router /api/auth/forgotpassword/otp [post]
func (c *UserController) VerifyOTPHandler(ctx *fiber.Ctx) error {
	var req models.VerifyOTPRequest

	if err := ctx.BodyParser(&req); err != nil {
		return ctx.Status(fiber.ErrInternalServerError.Code).JSON(fiber.Map{
			"status":      fiber.ErrInternalServerError.Message,
			"status_code": fiber.ErrInternalServerError.Code,
			"message":     err.Error(),
			"result":      nil,
		})
	}

	if req.Email == "" {
		return ctx.Status(fiber.ErrBadRequest.Code).JSON(fiber.Map{
			"status":      fiber.ErrBadRequest.Message,
			"status_code": fiber.ErrBadRequest.Code,
			"message":     "Email is missing",
			"result":      nil,
		})
	}

	if req.OTP == "" {
		return ctx.Status(fiber.ErrBadRequest.Code).JSON(fiber.Map{
			"status":      fiber.ErrBadRequest.Message,
			"status_code": fiber.ErrBadRequest.Code,
			"message":     "OTP is missing",
			"result":      nil,
		})
	}

	err := c.userusecase.VerifyOTP(req.Email, req.OTP) // เป็น _ เพราะไม่ใช้ token ตรงนี้มันจะถูกเทียบใน db
	if err != nil {
		return ctx.Status(fiber.ErrInternalServerError.Code).JSON(fiber.Map{
			"status":      fiber.ErrInternalServerError.Message,
			"status_code": fiber.ErrInternalServerError.Code,
			"message":     err.Error(),
			"result":      nil,
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":      "Success",
		"status_code": fiber.StatusOK,
		"message":     "OTP is correct",
	})
}

// ChangePasswordHandler godoc
// @Summary Change Password (Forgot Password Flow)
// @Description Change password after OTP verification in forgot password flow
// @Tags Authentication
// @Accept json
// @Produce json
// @Param request body object{email=string,new_password=string} true "Email and new password"
// @Success 200 {object} object{status=string,status_code=int,message=string} "Password changed successfully"
// @Failure 400 {object} object{status=string,status_code=int,message=string,result=any} "Bad Request - Email or new password is missing"
// @Failure 500 {object} object{status=string,status_code=int,message=string,result=any} "Internal Server Error"
// @Router /api/auth/forgotpassword/changepassword [patch]
func (c *UserController) ChangePasswordHandler(ctx *fiber.Ctx) error {
	var req models.ChangePasswordRequest

	if err := ctx.BodyParser(&req); err != nil {
		return ctx.Status(fiber.ErrInternalServerError.Code).JSON(fiber.Map{
			"status":      fiber.ErrInternalServerError.Message,
			"status_code": fiber.ErrInternalServerError.Code,
			"message":     err.Error(),
			"result":      nil,
		})
	}

	if req.Email == "" {
		return ctx.Status(fiber.ErrBadRequest.Code).JSON(fiber.Map{
			"status":      fiber.ErrBadRequest.Message,
			"status_code": fiber.ErrBadRequest.Code,
			"message":     "Email is missing",
			"result":      nil,
		})
	}

	if req.NewPassword == "" {
		return ctx.Status(fiber.ErrBadRequest.Code).JSON(fiber.Map{
			"status":      fiber.ErrBadRequest.Message,
			"status_code": fiber.ErrBadRequest.Code,
			"message":     "New password is missing",
			"result":      nil,
		})
	}

	err := c.userusecase.ChangePassword(req.Email, req.NewPassword)
	if err != nil {
		return ctx.Status(fiber.ErrInternalServerError.Code).JSON(fiber.Map{
			"status":      fiber.ErrInternalServerError.Message,
			"status_code": fiber.ErrInternalServerError.Code,
			"message":     err.Error(),
			"result":      nil,
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":      "Success",
		"status_code": fiber.StatusOK,
		"message":     "Password changed successfully",
	})
}

// GetUserByIDHandler godoc
// @Summary Get User Information
// @Description Get authenticated user's information
// @Tags User
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} object{status=string,status_code=int,message=string,result=object} "User Info retrieved successfully"
// @Failure 401 {object} object{status=string,status_code=int,message=string,result=any} "Unauthorized - Missing user ID"
// @Failure 404 {object} object{status=string,status_code=int,message=string,result=any} "User not found"
// @Router /api/user [get]
func (c *UserController) GetUserByIDHandler(ctx *fiber.Ctx) error {
	userID, ok := ctx.Locals("user_id").(string)
	if !ok || userID == "" {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"status":      fiber.ErrUnauthorized.Message,
			"status_code": fiber.StatusUnauthorized,
			"message":     "Unauthorized: Missing user ID",
			"result":      nil,
		})
	}

	data, err := c.userusecase.GetUserByID(userID)
	if err != nil {
		return ctx.Status(fiber.ErrNotFound.Code).JSON(fiber.Map{
			"status":      fiber.ErrNotFound.Message,
			"status_code": fiber.ErrNotFound.Code,
			"message":     err.Error(),
			"result":      nil,
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":      "Success",
		"status_code": fiber.StatusOK,
		"message":     "User Info retrieved successfully",
		"result":      data,
	})
}

// LogoutHandler godoc
// @Summary User Logout
// @Description Logout authenticated user
// @Tags Authentication
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} object{status=string,status_code=int,message=string,result=any} "Logout successful"
// @Router /api/auth/logout [post]
func (c *UserController) LogoutHandler(ctx *fiber.Ctx) error {
	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":      "Success",
		"status_code": fiber.StatusOK,
		"message":     "Logout successful",
		"result":      nil,
	})
}

// UpdateUserByIDHandler godoc
// @Summary Update User Information
// @Description Partially update authenticated user's profile. Only send fields that need to be updated.
// @Tags User
// @Accept multipart/form-data
// @Produce json
// @Security BearerAuth
// @Param username formData string false "Username"
// @Param first_name formData string false "First Name"
// @Param last_name formData string false "Last Name"
// @Param nickname formData string false "Nickname"
// @Param gender formData string false "Gender"
// @Param profile_image formData file false "Profile Image"
// @Success 200 {object} object{status=string,status_code=int,message=string,result=object} "User updated successfully"
// @Failure 400 {object} object{status=string,status_code=int,message=string,result=any} "Bad Request - Invalid form data"
// @Failure 401 {object} object{status=string,status_code=int,message=string,result=any} "Unauthorized - Missing user ID"
// @Failure 500 {object} object{status=string,status_code=int,message=string,result=any} "Internal Server Error"
// @Router /api/user [patch]
func (c *UserController) UpdateUserByIDHandler(ctx *fiber.Ctx) error {
	userID, ok := ctx.Locals("user_id").(string)
	if !ok || userID == "" {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"status":      fiber.ErrUnauthorized.Message,
			"status_code": fiber.StatusUnauthorized,
			"message":     "Unauthorized: Missing user ID",
			"result":      nil,
		})
	}

	form, err := ctx.MultipartForm()
	if err != nil {
		return ctx.Status(fiber.ErrBadRequest.Code).JSON(fiber.Map{
			"status":      fiber.ErrBadRequest.Message,
			"status_code": fiber.ErrBadRequest.Code,
			"message":     "Invalid form data: " + err.Error(),
			"result":      nil,
		})
	}

	user := &entities.User{}

	if usernames := form.Value["username"]; len(usernames) > 0 {
		user.Username = usernames[0]
	}
	if firstNames := form.Value["first_name"]; len(firstNames) > 0 {
		user.FirstName = firstNames[0]
	}
	if lastNames := form.Value["last_name"]; len(lastNames) > 0 {
		user.LastName = lastNames[0]
	}
	if nicknames := form.Value["nickname"]; len(nicknames) > 0 {
		user.Nickname = nicknames[0]
	}
	if genders := form.Value["gender"]; len(genders) > 0 {
		user.Gender = genders[0]
	}

	// Get the profile image files (optional)
	files := form.File["profile_image"]
	var file multipart.File

	if len(files) > 0 {
		// Use the first file if provided
		fileHeader := files[0]

		// Open the file
		file, err = fileHeader.Open()
		if err != nil {
			return ctx.Status(fiber.ErrBadRequest.Code).JSON(fiber.Map{
				"status":      fiber.ErrBadRequest.Message,
				"status_code": fiber.ErrBadRequest.Code,
				"message":     "Failed to open uploaded file: " + err.Error(),
				"result":      nil,
			})
		}
		defer file.Close()
	} // Call the usecase to update user with profile image
	updatedUser, err := c.userusecase.UpdateUserByID(userID, user, file)
	if err != nil {
		return ctx.Status(fiber.ErrInternalServerError.Code).JSON(fiber.Map{
			"status":      fiber.ErrInternalServerError.Message,
			"status_code": fiber.ErrInternalServerError.Code,
			"message":     err.Error(),
			"result":      nil,
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":      "Success",
		"status_code": fiber.StatusOK,
		"message":     "User updated successfully",
		"result":      updatedUser,
	})
}

// CreateStaffFileHandler godoc
// @Summary Upload Staff Files
// @Description Upload multiple files for staff member
// @Tags User
// @Accept multipart/form-data
// @Produce json
// @Security BearerAuth
// @Param file formData file true "Files to upload (can select multiple)"
// @Success 201 {object} object{status=string,status_code=int,message=string,result=array} "Staff files created successfully"
// @Failure 400 {object} object{status=string,status_code=int,message=string,result=any} "Bad Request - No files provided or invalid form data"
// @Failure 401 {object} object{status=string,status_code=int,message=string,result=any} "Unauthorized - Missing user ID"
// @Failure 500 {object} object{status=string,status_code=int,message=string,result=any} "Internal Server Error"
// @Router /api/user/staff/files [post]
func (c *UserController) CreateStaffFileHandler(ctx *fiber.Ctx) error {
	userID, ok := ctx.Locals("user_id").(string)
	if !ok || userID == "" {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"status":      fiber.ErrUnauthorized.Message,
			"status_code": fiber.StatusUnauthorized,
			"message":     "Unauthorized: Missing user ID",
			"result":      nil,
		})
	}

	form, err := ctx.MultipartForm()
	if err != nil {
		return ctx.Status(fiber.ErrBadRequest.Code).JSON(fiber.Map{
			"status":      fiber.ErrBadRequest.Message,
			"status_code": fiber.ErrBadRequest.Code,
			"message":     "Invalid form data: " + err.Error(),
			"result":      nil,
		})
	}

	// รับไฟล์หลายไฟล์จาก form
	files := form.File["file"]
	if len(files) == 0 {
		return ctx.Status(fiber.ErrBadRequest.Code).JSON(fiber.Map{
			"status":      fiber.ErrBadRequest.Message,
			"status_code": fiber.ErrBadRequest.Code,
			"message":     "No files provided",
			"result":      nil,
		})
	}

	// เรียก usecase เพื่อสร้างหลายไฟล์
	createdFiles, err := c.userusecase.CreateStaffFile(userID, files)
	if err != nil {
		return ctx.Status(fiber.ErrInternalServerError.Code).JSON(fiber.Map{
			"status":      fiber.ErrInternalServerError.Message,
			"status_code": fiber.ErrInternalServerError.Code,
			"message":     err.Error(),
			"result":      nil,
		})
	}

	return ctx.Status(fiber.StatusCreated).JSON(fiber.Map{
		"status":      "Success",
		"status_code": fiber.StatusCreated,
		"message":     "Staff files created successfully",
		"result":      createdFiles,
	})
}
