package models

type LoginRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password" binding:"required"`
	Remember bool   `json:"remember"`
}

type ResetPasswordRequest struct {
	OldPassword string `json:"old_password" binding:"required"`
	NewPassword string `json:"new_password" binding:"required"`
}

type ForgotPasswordRequest struct {
	Email string `json:"email" binding:"required,email"`
}

type VerifyOTPRequest struct {
	Email string `json:"email" binding:"required,email"`
	OTP   string `json:"otp" binding:"required"`
}

type ChangePasswordRequest struct {
	Email       string `json:"email" binding:"required,email"`
	NewPassword string `json:"new_password" binding:"required"`
}

type GetUsersByFirstAndLastNameRequest struct {
	FirstName string `query:"first_name" binding:"required"`
	LastName  string `query:"last_name" binding:"required"`
}

type UpdateStaffRoleRequest struct {
	RoleName string `json:"role_name" binding:"required"`
}

type UpdateUserApprovalRequest struct {
	IsApprove *bool `json:"is_approve"`
}
