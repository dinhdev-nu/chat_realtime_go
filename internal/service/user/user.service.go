package service

import (
	"github.com/dinhdev-nu/realtime_auth_go/internal/dto"
	"github.com/dinhdev-nu/realtime_auth_go/internal/model"
	"github.com/dinhdev-nu/realtime_auth_go/internal/repo"
)

type IUserService interface {
	GetUserInfoByName(username string) (*model.GoDbUserInfo, error)
	SearchUserByName(username string, userIdReqes int64) ([]*dto.SearchUsersOutput, error)
}

type userService struct {
	repo repo.IUserRepo
}

func NewUserService(repo repo.IUserRepo) IUserService {
	return &userService{
		repo: repo,
	}
}

func (us *userService) GetUserInfoByName(username string) (*model.GoDbUserInfo, error) {
	userInfo, err := us.repo.GetUserInfoByName(username)
	if err != nil {
		return userInfo, err
	}
	return userInfo, nil
}

func (us *userService) SearchUserByName(username string, userIdRes int64) ([]*dto.SearchUsersOutput, error) {
	users, err := us.repo.SearchUserByName(username, userIdRes)
	if err != nil {
		return nil, err
	}
	if len(users) == 0 {
		return users, nil
	}
	// Lấy trạng thái online/offline của từng user
	for _, user := range users {
		status, _ := us.repo.GetStatusByUserId(user.UserID)
		user.UserStatus = status // Thêm trạng thái vào từng user
	}
	return users, nil
}
