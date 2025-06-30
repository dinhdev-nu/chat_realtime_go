package repo

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/dinhdev-nu/realtime_auth_go/global"
	"github.com/dinhdev-nu/realtime_auth_go/internal/dto"
	"github.com/dinhdev-nu/realtime_auth_go/internal/model"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type IAuthRepo interface {
	AddOtp(email string, data dto.OtpValueRedisDTO, ttl int64) error
	IncrementOtp(email string, value dto.OtpValueRedisDTO) error
	GetOtp(email string) (dto.OtpValueRedisDTO, error)
	DelOtp(email string) error        // delete otp from redis
	DeleteFromRedis(key string) error // delete key from redis
	AddUserKey(key string, userInfo []byte) error

	GetExistEmail(email string) (bool, error)
	GetUserBase(email string) (*model.GoDbUserBase, error)
	GetUserInfoByID(userId int64) (*model.GoDbUserInfo, error)
	GetUserInfoByEmail(email string) (*model.GoDbUserInfo, error)
	CreateUserRegis(emailHash string, email string) error

	SaveOtpTodb(userid int64, email string, keyHash string, otp string) error
	UpdateOtpIndb(email string) error
	UpdateUserInfo(email string, data map[string]interface{}) error
	UpdateUserBase(email string, data map[string]interface{}) error
	UpdatePassword(email string, password string) error
	UpdateLoginUser(userId int64, loginIp string) error
	DeleteOtpUser(email string) error
}

type authRepo struct {
	ctx context.Context
}

func NewAuthRepo() IAuthRepo {
	ctx := context.Background()
	return &authRepo{
		ctx: ctx,
	}
}

// mysql

func (ar *authRepo) GetUserInfoByEmail(email string) (*model.GoDbUserInfo, error) {
	userInfo := model.GoDbUserInfo{}
	result := global.Mdb.Model(&userInfo).Where("user_account = ?", email).First(&userInfo)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, result.Error
	}
	return &userInfo, nil
}

func (ar *authRepo) CreateUserRegis(emailHash string, email string) error {
	userInfo := model.GoDbUserInfo{
		UserAccount:          emailHash,
		UserState:            2,
		UserEmail:            email,
		UserIsAuthentication: 0,
	}
	newUser := global.Mdb.Create(&userInfo)
	if newUser.Error != nil {
		return newUser.Error
	}
	userBase := model.GoDbUserBase{
		UserID:       userInfo.UserID,
		UserAccount:  emailHash,
		UserPassword: "",
		UserSalt:     "",
	}
	newUserBase := global.Mdb.Create(&userBase)
	if newUserBase.Error != nil {
		return newUser.Error
	}
	return nil
}

func (ar *authRepo) GetUserInfoByID(userId int64) (*model.GoDbUserInfo, error) {
	userInfo := model.GoDbUserInfo{}
	result := global.Mdb.Model(&userInfo).Where("user_id = ?", userId).First(&userInfo)
	if result.Error != nil {
		return nil, result.Error
	}
	return &userInfo, nil
}

func (ar *authRepo) UpdateLoginUser(userId int64, loginIp string) error {
	result := global.Mdb.Model(&model.GoDbUserBase{}).Where("user_id = ?", userId).Updates(
		map[string]interface{}{
			"user_login_time": time.Now(),
			"user_login_ip":   loginIp,
		},
	)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("no rows affected for login time update")
	}
	return nil
}

func (ar *authRepo) GetUserBase(email string) (*model.GoDbUserBase, error) {
	user := model.GoDbUserBase{}
	result := global.Mdb.Model(&user).Where("user_account = ?", email).First(&user)

	if result.Error != nil {
		return nil, result.Error
	}

	return &user, nil
}

func (ar *authRepo) UpdatePassword(email string, password string) error {
	result := global.Mdb.Model(&model.GoDbUserBase{}).
		Where("user_account = ?", email).
		Update("user_password", password)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("no rows affected for password update")
	}
	return nil
}

func (ar *authRepo) UpdateOtpIndb(email string) error {
	// update otp
	otpModel := model.GoDbVerifyOtp{}
	result := global.Mdb.Model(&otpModel).Where("verify_key_hash = ?", email).Updates(
		map[string]interface{}{
			"is_verified":       1,
			"verify_updated_at": time.Now(),
		},
	)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("no rows affected for otp update")
	}

	return nil
}

func (ar *authRepo) UpdateUserBase(email string, data map[string]interface{}) error {
	err := global.Mdb.Model(&model.GoDbUserBase{}).
		Where("user_account = ?", email).
		Updates(data)
	if err.Error != nil {
		return err.Error
	}
	if err.RowsAffected == 0 {
		return errors.New("no rows affected for user base update")
	}
	return nil
}

func (ar *authRepo) UpdateUserInfo(email string, data map[string]interface{}) error {

	result := global.Mdb.Model(&model.GoDbUserInfo{}).
		Where("user_account = ?", email).
		Updates(data)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("no rows affected for user info update")
	}
	return nil
}

func (ar *authRepo) SaveOtpTodb(userid int64, email string, keyHash string, otp string) error {
	otpModel := model.GoDbVerifyOtp{
		UserID:        userid,
		VerifyOtp:     otp,
		VerifyKey:     email,
		VerifyKeyHash: keyHash,
		VerifyType:    1,
		IsVerified:    0,
	}
	err := global.Mdb.Clauses(clause.OnConflict{
		Columns: []clause.Column{{Name: "verify_key"}},
		DoUpdates: clause.AssignmentColumns([]string{
			"verify_otp",
			"verify_key_hash",
			"verify_updated_at",
			"verify_type",
			"is_verified",
		}),
	}).Create(&otpModel).Error
	if err != nil {
		return err
	}
	// if global.Mdb.Model(&otpModel).Where("verify_key = ?", email).Updates(map[string]interface{}{}).RowsAffected == 0 {
	// 	global.Mdb.Create(&otpModel)
	// }
	return nil
}

func (ar *authRepo) GetExistEmail(email string) (bool, error) {
	// check email exist in db
	var count int64
	err := global.Mdb.Model(&model.GoDbUserInfo{}).
		Where("user_account = ? AND user_state !=?", email, 2).
		Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func (ar *authRepo) DeleteOtpUser(email string) error {
	// delete otp
	otpModel := model.GoDbVerifyOtp{}
	result := global.Mdb.Model(&otpModel).Where("verify_key_hash = ?", email).Delete(&otpModel)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("no rows affected for otp delete")
	}
	return nil
}

// redis

func (ar *authRepo) AddUserKey(key string, userInfo []byte) error {
	ttl := 60 * 60 // 60 minutes
	err := global.Rdb.SetEx(ar.ctx, key, userInfo, time.Duration(ttl)*time.Second).Err()
	if err != nil {
		return err
	}
	return nil
}

func (ar *authRepo) IncrementOtp(email string, value dto.OtpValueRedisDTO) error {
	key := "otp:" + email + ":usr"
	// get ttl
	ttl, err := global.Rdb.TTL(ar.ctx, key).Result()
	if err != nil {
		return err
	}
	value.FailCount++
	jsonData, err := json.Marshal(value)
	if err != nil {
		return err
	}
	return global.Rdb.SetEx(ar.ctx, key, jsonData, ttl).Err()
}

func (ar *authRepo) AddOtp(email string, data dto.OtpValueRedisDTO, ttl int64) error {
	key := "otp:" + email + ":usr"
	jsonData, err := json.Marshal(data)
	if err != nil {
		return err
	}

	return global.Rdb.SetEx(ar.ctx, key, jsonData, time.Duration(ttl)*time.Second).Err()
}

func (ar *authRepo) GetOtp(email string) (dto.OtpValueRedisDTO, error) {
	key := "otp:" + email + ":usr"
	jsonData, err := global.Rdb.Get(ar.ctx, key).Result()
	if err != nil || jsonData == "" {
		return dto.OtpValueRedisDTO{}, errors.New("otp not found or empty")
	}
	var data dto.OtpValueRedisDTO
	err = json.Unmarshal([]byte(jsonData), &data)
	if err != nil {
		return dto.OtpValueRedisDTO{}, err
	}
	return data, nil
}

// DelOtp delete otp from redis
func (ar *authRepo) DelOtp(email string) error {
	key := "otp:" + email + ":usr"
	return global.Rdb.Del(ar.ctx, key).Err()
}

func (ar *authRepo) DeleteFromRedis(key string) error {
	return global.Rdb.Del(ar.ctx, key).Err()
}
