package controllers

import (
	// "fmt"
	// "mime/multipart"
	// "strconv"
	"mime/multipart"

	"github.com/aikidoaikido115/New-Acis-BE/modules/entities"
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

func (c *UserController) RegisterHandler(ctx *fiber.Ctx) error {
	var req struct {
		Username string `json:"username"`
		Email    string `json:"email"`
		Password string `json:"password"`
		RoleName string `json:"role_name"`
	}

	if err := ctx.BodyParser(&req); err != nil {
		return ctx.Status(fiber.ErrBadRequest.Code).JSON(fiber.Map{
			"status":      fiber.ErrBadRequest.Message,
			"status_code": fiber.ErrBadRequest.Code,
			"message":     err.Error(),
			"result":      nil,
		})
	}

	if req.Username == "" {
		return ctx.Status(fiber.ErrBadRequest.Code).JSON(fiber.Map{
			"status":      "Error",
			"status_code": fiber.ErrBadRequest.Code,
			"message":     "Username is missing",
			"result":      nil,
		})
	}

	if req.Email == "" {
		return ctx.Status(fiber.ErrBadRequest.Code).JSON(fiber.Map{
			"status":      "Error",
			"status_code": fiber.ErrBadRequest.Code,
			"message":     "Email is missing",
			"result":      nil,
		})
	}

	if req.Password == "" {
		return ctx.Status(fiber.ErrBadRequest.Code).JSON(fiber.Map{
			"status":      "Error",
			"status_code": fiber.ErrBadRequest.Code,
			"message":     "Password is missing",
			"result":      nil,
		})
	}

	user := &entities.User{
		Username: req.Username,
		Email:    req.Email,
		Password: req.Password,
	}

	data, err := c.userusecase.Register(user, req.RoleName)
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
		"status_code": fiber.StatusOK,
		"message":     "user created successfully",
		"result":      data,
	})
}

func (c *UserController) LoginHandler(ctx *fiber.Ctx) error {
	var req struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	if err := ctx.BodyParser(&req); err != nil {
		return ctx.Status(fiber.ErrBadRequest.Code).JSON(fiber.Map{
			"status":      fiber.ErrBadRequest.Message,
			"status_code": fiber.ErrBadRequest.Code,
			"message":     err.Error(),
			"result":      nil,
		})
	}

	if req.Username == "" {
		return ctx.Status(fiber.ErrBadRequest.Code).JSON(fiber.Map{
			"status":      "Error",
			"status_code": fiber.ErrBadRequest.Code,
			"message":     "Username is missing",
			"result":      nil,
		})
	}

	if req.Password == "" {
		return ctx.Status(fiber.ErrBadRequest.Code).JSON(fiber.Map{
			"status":      "Error",
			"status_code": fiber.ErrBadRequest.Code,
			"message":     "Password is missing",
			"result":      nil,
		})
	}

	token, user, err := c.userusecase.Login(req.Username, req.Password)
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

func (c *UserController) ResetPasswordHandler(ctx *fiber.Ctx) error {
	var req struct {
		Email       string `json:"email"`
		OldPassword string `json:"old_password"`
		NewPassword string `json:"new_password"`
	}

	if err := ctx.BodyParser(&req); err != nil {
		return ctx.Status(fiber.ErrBadRequest.Code).JSON(fiber.Map{
			"status":      "Error",
			"status_code": fiber.ErrBadRequest.Code,
			"message":     err.Error(),
			"result":      nil,
		})
	}

	if req.OldPassword == "" || req.NewPassword == "" {
		return ctx.Status(fiber.ErrBadRequest.Code).JSON(fiber.Map{
			"status":      "Error",
			"status_code": fiber.ErrBadRequest.Code,
			"message":     "Old password and new password is missing",
			"result":      nil,
		})
	}

	userID, ok := ctx.Locals("user_id").(string)
	if !ok || userID == "" {
		return ctx.Status(fiber.ErrUnauthorized.Code).JSON(fiber.Map{
			"status":      "Error",
			"status_code": fiber.ErrUnauthorized.Code,
			"message":     "Unauthorized Missing user ID",
			"result":      nil,
		})
	}

	err := c.userusecase.ResetPassword(userID, req.OldPassword, req.NewPassword)
	if err != nil {
		// Return more specific status codes based on error type
		if err.Error() == "user invalid" {
			return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"status":      "Error",
				"status_code": fiber.StatusNotFound,
				"message":     err.Error(),
				"result":      nil,
			})
		} else if err.Error() == "old password invalid" {
			return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"status":      "Error",
				"status_code": fiber.StatusBadRequest,
				"message":     err.Error(),
				"result":      nil,
			})
		} else {
			return ctx.Status(fiber.ErrInternalServerError.Code).JSON(fiber.Map{
				"status":      "Error",
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

func (c *UserController) ForgotPasswordHandler(ctx *fiber.Ctx) error {
	type ForgotPasswordRequest struct {
		Email string `json:"email" validate:"required,email"`
	}

	var req ForgotPasswordRequest
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
			"status":      "Error",
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

func (c *UserController) VerifyOTPHandler(ctx *fiber.Ctx) error {
	type OTPRequest struct {
		Email string `json:"email" validate:"required,email"`
		OTP   string `json:"otp"`
	}

	var req OTPRequest
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
			"status":      "Error",
			"status_code": fiber.ErrBadRequest.Code,
			"message":     "Email is missing",
			"result":      nil,
		})
	}

	if req.OTP == "" {
		return ctx.Status(fiber.ErrBadRequest.Code).JSON(fiber.Map{
			"status":      "Error",
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

func (c *UserController) ChangePasswordHandler(ctx *fiber.Ctx) error {
	type OTPRequest struct {
		Email       string `json:"email" validate:"required,email"`
		NewPassword string `json:"new_password"`
	}

	var req OTPRequest
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
			"status":      "Error",
			"status_code": fiber.ErrBadRequest.Code,
			"message":     "Email is missing",
			"result":      nil,
		})
	}

	if req.NewPassword == "" {
		return ctx.Status(fiber.ErrBadRequest.Code).JSON(fiber.Map{
			"status":      "Error",
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

func (c *UserController) GetUserByIDHandler(ctx *fiber.Ctx) error {
	userID, ok := ctx.Locals("user_id").(string)
	if !ok || userID == "" {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"status":      "Error",
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

func (c *UserController) LogoutHandler(ctx *fiber.Ctx) error {
	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":      "Success",
		"status_code": fiber.StatusOK,
		"message":     "Logout successful",
		"result":      nil,
	})
}

func (c *UserController) UpdateUserByIDHandler(ctx *fiber.Ctx) error {
	userID, ok := ctx.Locals("user_id").(string)
	if !ok || userID == "" {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"status":      "Error",
			"status_code": fiber.StatusUnauthorized,
			"message":     "Unauthorized: Missing user ID",
			"result":      nil,
		})
	}


	form, err := ctx.MultipartForm()
	if err != nil {
		return ctx.Status(fiber.ErrBadRequest.Code).JSON(fiber.Map{
			"status":      "Error",
			"status_code": fiber.ErrBadRequest.Code,
			"message":     "Invalid form data",
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
				"status":      "Error",
				"status_code": fiber.ErrBadRequest.Code,
				"message":     "Failed to open uploaded file",
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
