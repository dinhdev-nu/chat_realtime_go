package service

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/dinhdev-nu/realtime_auth_go/config"
	"github.com/dinhdev-nu/realtime_auth_go/internal/dto"
	"github.com/dinhdev-nu/realtime_auth_go/internal/repo"
	"github.com/dinhdev-nu/realtime_auth_go/internal/utils"
	"github.com/dinhdev-nu/realtime_auth_go/internal/utils/crypto"
	"github.com/dinhdev-nu/realtime_auth_go/internal/utils/jwt"
	"github.com/dinhdev-nu/realtime_auth_go/internal/utils/random"
	"github.com/dinhdev-nu/realtime_auth_go/pkg/response"
)

type IAuthService interface {
	Register(email string) (dto.RegisterDTO, int)
	SendOtp(email string) int // int là mã lỗi
	VeryfyOtp(email string, otp string) int
	DelOtp(email string) int
	SignUp(email string, password string) int
	Login(email string, password string, loginIp string) (dto.LoginOutput, int) // nên trả về 1 struct output
	UpdatePassword(assword string) int
	Logout(email string, uuidToken string) int //
}

type authService struct { // Viết thường là private và viết hoa là public
	repo repo.IAuthRepo
}

func NewAuthService(authRepo repo.IAuthRepo) IAuthService {
	return &authService{repo: authRepo}
}

func (as *authService) Register(email string) (dto.RegisterDTO, int) {
	// hash email tránh lộ email để otp
	emailHash := crypto.HashEmail(email)

	// Check email exist
	userInfo, err := as.repo.GetUserInfoByEmail(emailHash)
	if err != nil {
		fmt.Println("error server " + err.Error())
		return dto.RegisterDTO{}, response.ErrorCode
	}
	if userInfo != nil && userInfo.UserState != 2 {
		fmt.Println("email has exist")
		return dto.RegisterDTO{}, response.ErrorCodeEmailExist
	}

	if userInfo == nil {
		err := as.repo.CreateUserRegis(emailHash, email)
		if err != nil {
			return dto.RegisterDTO{}, response.ErrorCreateCode
		}
	}

	return dto.RegisterDTO{
		Email: email,
	}, response.SuccessCode
}

func (as *authService) SendOtp(email string) int {
	emailHash := crypto.HashEmail(email)

	// check otp exist
	if _, err := as.repo.GetOtp(emailHash); err == nil {
		fmt.Println("error otp exist", err)
		return response.ErrorCreateCode
	}
	// check email
	userBase, err := as.repo.GetUserBase(emailHash)
	if err != nil {
		fmt.Println("error get user " + err.Error())
		return response.ErrorUserNotExist
	}
	if userBase == nil {
		fmt.Println("error user not exist")
		return response.ErrorUserNotExist
	}

	// generaet opt
	opt := random.CreateOtp()

	// send otp
	err = config.SendOTPEmailByTemplate(email, opt) // implement for 5s
	if err != nil {
		fmt.Println("error send otp" + err.Error())
		return response.ErrorOtpFail
	}

	// save otp vào redis
	data := dto.OtpValueRedisDTO{
		OTP:       opt,
		FailCount: 0,
	}
	err = as.repo.AddOtp(emailHash, data, 180) // 3 phút
	if err != nil {
		fmt.Println("error save otp to redis" + err.Error())
		return response.ErrorOtpFail
	}
	// save otp to mysql
	go as.repo.SaveOtpTodb(userBase.UserID, email, emailHash, opt)

	return response.SuccessCode
}

func (as *authService) VeryfyOtp(email string, otp string) int {
	// get otp form redis
	hashEmail := crypto.HashEmail(email)
	data, err := as.repo.GetOtp(hashEmail)
	if err != nil {
		fmt.Println("error get otp from redis")
		return response.ErrorOtpNotExist
	}
	// check otp fail count
	if data.FailCount > 3 {
		err := as.repo.DelOtp(hashEmail)
		if err != nil {
			fmt.Println("error delete otp from redis")
			return response.ErrorOtpFail
		}
		fmt.Println("otp fail count > 3")
		return response.ErrorOtpFail
	}

	// compare otp
	if data.OTP != otp {
		// increment fail count
		err := as.repo.IncrementOtp(hashEmail, data)
		if err != nil {
			fmt.Println("error increment otp fail count" + err.Error())
			return response.ErrorOtpFail
		}

		fmt.Println("otp not match" + data.OTP + " != " + otp)
		return response.ErrorOtpFail
	}
	// delete otp
	err = as.repo.DelOtp(hashEmail)
	if err != nil {
		fmt.Println("error delete otp from redis")
		return response.ErrorOtpFail
	}
	// update otp in mysql
	go as.repo.UpdateOtpIndb(hashEmail)

	return response.SuccessCode
}

func (as *authService) SignUp(email string, password string) int {
	// hash email
	emailHash := crypto.HashEmail(email)

	// // recheck email exist
	// emailExist, err := as.repo.GetExistEmail(emailHash)
	// if err != nil {
	// 	fmt.Println("error check email exist" + err.Error())
	// 	return response.ErrorCodeEmailExist
	// }
	// if !emailExist {
	// 	fmt.Println("email not exist")
	// 	return response.ErrorUserNotExist
	// }

	// chech user state 1 is active
	user, err := as.repo.GetUserInfoByEmail(emailHash)
	if err != nil {
		fmt.Println("error get user" + err.Error())
		return response.ErrorUserNotExist
	}
	if user.UserState != 2 {
		return response.SuccessCode
	}

	// hash password
	salt, err := crypto.CreateSalt()
	if err != nil {
		return response.ErrorCreateCode
	}
	passwordHash := crypto.HashPassword(password, salt)
	// create user base account
	data := map[string]interface{}{
		"user_salt":     salt,
		"user_password": passwordHash,
	}
	err = as.repo.UpdateUserBase(emailHash, data)
	if err != nil {
		fmt.Println("error create user base" + err.Error())
		return response.ErrorCreateCode
	}

	// create user info default
	data = map[string]interface{}{
		"user_nickname": utils.GennarateUserName(email),
		"user_state":    1,
	}
	err = as.repo.UpdateUserInfo(emailHash, data)
	if err != nil {
		fmt.Println("error create user info" + err.Error())
		return response.ErrorCreateCode
	}

	return response.SuccessCode
}

func (as *authService) Login(email string, password string, loginIp string) (dto.LoginOutput, int) {
	// hash email
	emailHash := crypto.HashEmail(email)
	//Get user if exist
	user, err := as.repo.GetUserBase(emailHash)
	if err != nil {
		fmt.Println("error get user" + err.Error())
		return dto.LoginOutput{}, response.ErrorUserNotExist
	}
	userInfo, err := as.repo.GetUserInfoByID(user.UserID)
	if err != nil {
		fmt.Println("error get user info" + err.Error())
		return dto.LoginOutput{}, response.ErrorUserNotExist
	}
	switch userInfo.UserState {
	case 0: // user state 0 is banned
		fmt.Println("error user state is 0")
		return dto.LoginOutput{}, response.ErrorCode
	case 2: // user state 2 is not active
		fmt.Println("error user state is 2")
		return dto.LoginOutput{}, response.ErrorCode
	}

	// check password
	if !crypto.VerifyPassword(password, user.UserPassword, user.UserSalt) {
		fmt.Println("error password not match")
		return dto.LoginOutput{}, response.ErrorPasswordNotMatch
	}
	// check 2 factor auth

	// create token
	uuidToken := utils.GenerateUUIDToken(user.UserID)
	// get user info
	userData, err := json.Marshal(&userInfo)
	if err != nil {
		fmt.Println("error marshal user info" + err.Error())
		return dto.LoginOutput{}, response.ErrorUserNotExist
	}
	// save token to redis
	err = as.repo.AddUserKey(uuidToken, userData)

	if err != nil {
		fmt.Println("error save token to redis" + err.Error())
		return dto.LoginOutput{}, response.ErrorCreateCode
	}

	// create token
	token, err := jwt.CreateToken(uuidToken)
	if err != nil {
		fmt.Println("error create token" + err.Error())
		return dto.LoginOutput{}, response.ErrorCreateCode
	}

	// update user last login
	err = as.repo.UpdateLoginUser(user.UserID, loginIp)
	if err != nil {
		return dto.LoginOutput{}, response.ErrorUpdateLoginCode
	}
	return dto.LoginOutput{
		Token: token,
		User:  userInfo,
	}, response.SuccessCode

}

func (as *authService) UpdatePassword(password string) int {

	return response.SuccessCode
}

func (as *authService) Logout(email string, uuidToken string) int {
	hashEmail := crypto.HashEmail(email)
	// delete token in redis
	err := as.repo.DeleteFromRedis(uuidToken)
	if err != nil {
		fmt.Println("error delete token from redis" + err.Error())
		return response.ErrorCreateCode
	}
	// update user logout time
	data := map[string]interface{}{
		"user_logout_time": time.Now(),
	}
	if err := as.repo.UpdateUserBase(hashEmail, data); err != nil {
		fmt.Println("error update user logout time" + err.Error())
		return response.ErrorUpdateCode
	}

	return response.SuccessCode
}

// delete otp and send otp => resend otp
func (as *authService) DelOtp(email string) int {
	hashEmail := crypto.HashEmail(email)

	// delete otp in redis
	err := as.repo.DelOtp(hashEmail)
	if err != nil {
		return response.ErrorDeleteCode
	}

	// delete otp in mysql
	err = as.repo.DeleteOtpUser(hashEmail)
	if err != nil {
		return response.ErrorDeleteCode
	}

	return response.SuccessCode
}
