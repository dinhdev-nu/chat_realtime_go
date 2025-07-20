package repo

import (
	"fmt"

	"github.com/dinhdev-nu/realtime_auth_go/global"
	"github.com/dinhdev-nu/realtime_auth_go/internal/dto"
	"github.com/dinhdev-nu/realtime_auth_go/internal/model"
	"github.com/gin-gonic/gin"
)

type IUserRepo interface {
	GetUserInfoByName(username string) (*model.GoDbUserInfo, error)
	GetUserInfoByIDs(userIDs []uint64) ([]*model.GoDbUserInfo, error)
	SearchUserByName(username string, userIdRes int64) ([]*dto.SearchUsersOutput, error)

	GetStatusByUserId(userId int64) (string, error)
}
type UserRepo struct{}

func NewUserRepo() IUserRepo {
	return &UserRepo{}
}

func (ur *UserRepo) GetUserInfoByIDs(userIDs []uint64) ([]*model.GoDbUserInfo, error) {
	var users []*model.GoDbUserInfo
	result := global.Mdb.Model(&model.GoDbUserInfo{}).
		Select("user_id, user_nickname, user_avatar").
		Where("user_id IN ?", userIDs).
		Find(&users)
	if result.Error != nil {
		if result.Error.Error() == "record not found" {
			return nil, nil // nếu không tìm thấy thì trả về nil
		}
		return nil, result.Error // nếu có lỗi khác thì trả về lỗi
	}
	return users, nil // nếu tìm thấy thì trả về danh sách người dùng
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

func (ur *UserRepo) GetUserInfoByName(username string) (*model.GoDbUserInfo, error) {

	var userInfo *model.GoDbUserInfo
	result := global.Mdb.Model(&model.GoDbUserInfo{}).
		Where("user_nickname = ?", username).
		First(&userInfo)
	if result.Error != nil {
		return userInfo, result.Error
	}
	return userInfo, nil
}

func (ur *UserRepo) SearchUserByName(username string, userIdRes int64) ([]*dto.SearchUsersOutput, error) {
	var users []*dto.SearchUsersOutput
	result := global.Mdb.Model(&model.GoDbUserInfo{}).
		Select("user_nickname, user_avatar, user_id, user_gender").
		Where("user_nickname LIKE ? AND user_id != ?", "%"+username+"%", userIdRes).
		Limit(LIMIT_USER).Scan(&users)
	if result.Error != nil {
		return nil, result.Error
	}
	return users, nil
}
