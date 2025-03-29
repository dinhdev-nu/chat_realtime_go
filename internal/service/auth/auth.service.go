package service

import (
	"fmt"

	"github.com/dinhdev-nu/realtime_auth_go/internal/repo"
	"github.com/dinhdev-nu/realtime_auth_go/internal/utils/crypto"
	"github.com/dinhdev-nu/realtime_auth_go/internal/utils/random"
	"github.com/dinhdev-nu/realtime_auth_go/pkg/response"
	"github.com/gin-gonic/gin"
)

type IAuthService interface {
	Register(email string) (*gin.H, int)
	SendOtp(email string) int // int là mã lỗi
	VeryfyOtp(email string, otp string) (string, int)
}

type authService struct { // Viết thường là private và viết hoa là public
	repo repo.IAuthRepo
}

func NewAuthService(authRepo repo.IAuthRepo) IAuthService {
	return &authService{repo: authRepo}
}
func (as *authService) Register(email string) (*gin.H, int) {

	// hash email tránh lộ email để otp
	emailHash := crypto.HashEmail(email)
	fmt.Println("email hash ::: " + emailHash)

	// Check email exist
	if as.repo.GetExistEmail(emailHash) {
		return nil, response.ErrorCodeEmailExist
	}

	return &gin.H{
		"message": "ok",
		"pass":    emailHash,
	}, response.SuccessCode
}

func (as *authService) SendOtp(email string) int {
	emailHash := crypto.HashEmail(email)

	// generaet opt
	opt := random.CreateOtp()
	// send otp
	fmt.Printf("otp ::: %s\n", opt)
	// save otp vào redis
	data := map[string]interface{}{
		"otp":        opt,
		"fail_count": 0,
	}
	err := as.repo.AddOtp(emailHash, data, 180) // 3 phút
	if err != nil {
		fmt.Println("error save otp to redis" + err.Error())
		return response.ErrorOtpFail
	}

	return response.SuccessCode
}

func (as *authService) VeryfyOtp(email string, otp string) (string, int) {
	// get otp form redis
	hashEmail := crypto.HashEmail(email)
	data := as.repo.GetOtp(hashEmail)
	if data == nil {
		fmt.Println("error get otp from redis")
		return "", response.ErrorOtpNotExist
	}
	// check otp fail count
	if data["fail_count"].(float64) > 3 {
		err := as.repo.DelOtp(hashEmail)
		if err != nil {
			fmt.Println("error delete otp from redis")
			return "", response.ErrorOtpFail
		}
		fmt.Println("otp fail count > 3")
		return "", response.ErrorOtpFail
	}

	// compare otp
	if data["otp"].(string) != otp {
		fmt.Println("otp not match" + data["otp"].(string) + " != " + otp)
		return "", response.ErrorOtpFail
	}
	// delete otp
	err := as.repo.DelOtp(hashEmail)
	if err != nil {
		fmt.Println("error delete otp from redis")
		return "", response.ErrorOtpFail
	}
	return hashEmail, response.SuccessCode
}
