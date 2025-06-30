package dto

import "github.com/dinhdev-nu/realtime_auth_go/internal/model"

type RegisterDTO struct {
	Email string `json:"email" binding:"required,email"`
}

type SendOtpDTO struct {
	Email string `json:"email" binding:"required,email"`
}

type VerifyOtpDTO struct {
	Email string `json:"email" binding:"required,email"`
	Otp   string `json:"otp" binding:"required"`
}

type SignUpInput struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type LoginInput struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
	LoginIp  string `json:"login_ip"`
}

type LoginOutput struct {
	Token string              `json:"token"`
	User  *model.GoDbUserInfo `json:"user"`
}

type LogoutInput struct {
	Email     string `json:"email" binding:"required,email"`
	UuidToken string `json:"uuid_token"`
}

type EmailInput struct {
	Email string `json:"email" binding:"required,email"`
}

type OtpValueRedisDTO struct {
	OTP       string `json:"otp"`
	FailCount int    `json:"fail_count"`
}
