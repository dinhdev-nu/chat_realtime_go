package repo

import (
	"fmt"

	"github.com/dinhdev-nu/realtime_auth_go/global"
	"github.com/dinhdev-nu/realtime_auth_go/internal/model"
	"github.com/gin-gonic/gin"
)

type IUserRepo interface {
	GetUserInfoByName(username string) (map[string]interface{}, error)
	SearchUserByName(username string, userIdRes int64) ([]map[string]interface{}, error)

	GetStatusByUserId(userId int64) (string, error)
}
type UserRepo struct{}

func NewUserRepo() IUserRepo {
	return &UserRepo{}
}

func (ur *UserRepo) GetStatusByUserId(userId int64) (string, error) {
	var status string
	key := fmt.Sprintf("user:%d:presence", userId)
	result, err := global.Rdb.Get(&gin.Context{}, key).Result()
	if err != nil {
		status = "offline" // nếu không tìm thấy key thì trả về offline
		return status, nil
	}
	status = result // nếu tìm thấy key thì trả về giá trị của key
	return status, nil
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
