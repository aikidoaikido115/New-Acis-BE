package models

type LoginRequest struct {
    Username string `json:"username"`
    Email    string `json:"email"`
    Password string `json:"password" binding:"required"`
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