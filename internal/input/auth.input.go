package input

type EmailInput struct {
	Email string `json:"email" binding:"required,email"`
}

type OtpInput struct {
	Email string `json:"email" binding:"required,email"`
	Otp   string `json:"otp" binding:"required"`
}

type SendOtpInput struct {
	Email   string `json:"email" binding:"required,email"`
	Purpose int    `json:"purpose" binding:"required"`
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

type LogoutInput struct {
	Email     string `json:"email" binding:"required,email"`
	UuidToken string `json:"uuid_token"`
}
