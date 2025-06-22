package repo

import (
	"github.com/dinhdev-nu/realtime_auth_go/global"
	"github.com/dinhdev-nu/realtime_auth_go/internal/model"
)

type IUserRepo interface {
	GetUserInfoByName(username string) (map[string]interface{}, error)
	SearchUserByName(username string, userIdRes int64) ([]map[string]interface{}, error)
}
type UserRepo struct{}

func NewUserRepo() IUserRepo {
	return &UserRepo{}
}

func (ur *UserRepo) GetUserInfoByName(username string) (map[string]interface{}, error) {

	var userInfo map[string]interface{}
	result := global.Mdb.Model(&model.GoDbUserInfo{}).
		Where("user_nickname = ?", username).
		First(&userInfo)
	if result.Error != nil {
		return userInfo, result.Error
	}
	return userInfo, nil
}

func (ur *UserRepo) SearchUserByName(username string, userIdRes int64) ([]map[string]interface{}, error) {
	var users []map[string]interface{}
	result := global.Mdb.Model(&model.GoDbUserInfo{}).
		Select("user_nickname, user_avatar, user_id, user_gender").
		Where("user_nickname LIKE ? AND user_id != ?", "%"+username+"%", userIdRes).
		Limit(LIMIT_USER).Find(&users)
	if result.Error != nil {
		return nil, result.Error
	}
	return users, nil
}
