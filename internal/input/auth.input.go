package input

type EmailInput struct {
	Email string `json:"email" binding:"required,email"`
}

type OtpInput struct {
	Email string `json:"email" binding:"required,email"`
	Otp   string `json:"otp" binding:"required"`
}
