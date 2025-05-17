package service

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/dinhdev-nu/realtime_auth_go/config"
	"github.com/dinhdev-nu/realtime_auth_go/internal/repo"
	"github.com/dinhdev-nu/realtime_auth_go/internal/utils"
	"github.com/dinhdev-nu/realtime_auth_go/internal/utils/crypto"
	"github.com/dinhdev-nu/realtime_auth_go/internal/utils/jwt"
	"github.com/dinhdev-nu/realtime_auth_go/internal/utils/random"
	"github.com/dinhdev-nu/realtime_auth_go/pkg/response"
)

type IAuthService interface {
	Register(email string) (map[string]string, int)
	SendOtp(email string) int // int là mã lỗi
	VeryfyOtp(email string, otp string) (string, int)
	SignUp(email string, password string) int
	Login(email string, password string, loginIp string) (map[string]interface{}, int) // nên trả về 1 struct output
	UpdatePassword(assword string) int
	Logout(email string, uuidToken string) int //
}

type authService struct { // Viết thường là private và viết hoa là public
	repo repo.IAuthRepo
}

func NewAuthService(authRepo repo.IAuthRepo) IAuthService {
	return &authService{repo: authRepo}
}

func (as *authService) Register(email string) (map[string]string, int) {

	// hash email tránh lộ email để otp
	emailHash := crypto.HashEmail(email)

	// Check email exist
	emailExist, err := as.repo.GetExistEmail(emailHash)
	if err != nil {
		fmt.Println("error check email exist" + err.Error())
		return nil, response.ErrorCodeEmailExist
	}
	if emailExist {
		return nil, response.ErrorCodeEmailExist
	}

	// create user info default
	go as.repo.CreateUserRegis(emailHash, email)

	return map[string]string{
		"email": email,
	}, response.SuccessCode
}

func (as *authService) SendOtp(email string) int {
	emailHash := crypto.HashEmail(email)
	// check otp exist
	if err := as.repo.GetOtp(emailHash); err != nil {
		fmt.Println("error otp exist")
		return response.ErrorCreateCode
	}
	// check email
	userBase, err := as.repo.GetUser(emailHash)
	if err != nil {
		fmt.Println("error get user" + err.Error())
		return response.ErrorUserNotExist
	}
	// generaet opt
	opt := random.CreateOtp()
	// send otp
	err = config.SendOTPEmailByTemplate(email, opt)
	if err != nil {
		fmt.Println("error send otp" + err.Error())
		return response.ErrorOtpFail
	}
	// save otp vào redis
	data := map[string]interface{}{
		"otp":        opt,
		"fail_count": 0,
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
		// increment fail count
		err := as.repo.IncrementOtp(hashEmail, data)
		if err != nil {
			fmt.Println("error increment otp fail count" + err.Error())
			return "", response.ErrorOtpFail
		}

		fmt.Println("otp not match" + data["otp"].(string) + " != " + otp)
		return "", response.ErrorOtpFail
	}
	// delete otp
	err := as.repo.DelOtp(hashEmail)
	if err != nil {
		fmt.Println("error delete otp from redis")
		return "", response.ErrorOtpFail
	}
	// update otp in mysql
	go as.repo.UpdateOtpIndb(hashEmail)

	return hashEmail, response.SuccessCode
}

func (as *authService) SignUp(email string, password string) int {
	// hash email
	emailHash := crypto.HashEmail(email)
	// recheck email exist
	emailExist, err := as.repo.GetExistEmail(emailHash)
	if err != nil {
		fmt.Println("error check email exist" + err.Error())
		return response.ErrorCodeEmailExist
	}
	if !emailExist {
		fmt.Println("email not exist")
		return response.ErrorUserNotExist
	}
	// chech user state 1 is active
	user, err := as.repo.GetUserInfoByEmail(emailHash)
	if err != nil {
		fmt.Println("error get user" + err.Error())
		return response.ErrorUserNotExist
	}
	if user.UserState == 1 {
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

func (as *authService) Login(email string, password string, loginIp string) (map[string]interface{}, int) {
	// hash email
	emailHash := crypto.HashEmail(email)
	//Get user if exist
	user, err := as.repo.GetUser(emailHash)
	if err != nil {
		fmt.Println("error get user" + err.Error())
		return nil, response.ErrorUserNotExist
	}
	userInfo, err := as.repo.GetUserInfoByID(user.UserID)
	if err != nil {
		fmt.Println("error get user info" + err.Error())
		return nil, response.ErrorUserNotExist
	}
	switch userInfo.UserState {
	case 0: // user state 0 is banned
		fmt.Println("error user state is 0")
		return nil, response.ErrorCode
	case 2: // user state 2 is not active
		fmt.Println("error user state is 2")
		return nil, response.ErrorCode
	}

	// check password
	if !crypto.VerifyPassword(password, user.UserPassword, user.UserSalt) {
		fmt.Println("error password not match")
		return nil, response.ErrorPasswordNotMatch
	}
	// check 2 factor auth

	// create token
	uuidToken := utils.GenerateUUIDToken(user.UserID)
	// get user info
	userData, err := json.Marshal(&userInfo)
	if err != nil {
		fmt.Println("error marshal user info" + err.Error())
		return nil, response.ErrorUserNotExist
	}
	// save token to redis
	err = as.repo.AddUserKey(uuidToken, userData)

	if err != nil {
		fmt.Println("error save token to redis" + err.Error())
		return nil, response.ErrorCreateCode
	}

	// create token
	token, err := jwt.CreateToken(uuidToken)
	if err != nil {
		fmt.Println("error create token" + err.Error())
		return nil, response.ErrorCreateCode
	}

	// update user last login
	err = as.repo.UpdateLoginUser(user.UserID, loginIp)
	if err != nil {
		return nil, response.ErrorUpdateLoginCode
	}
	return map[string]interface{}{
		"token": token,
		"id":    user.UserID,
		"email": userInfo.UserEmail,
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
