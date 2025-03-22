package service

import (
	"github.com/dinhdev-nu/realtime_auth_go/internal/repo"
	"github.com/gin-gonic/gin"
)

type IAuthService interface {
	Ping() *gin.H
}

type authService struct { // Viết thường là private và viết hoa là public
	repo repo.IAuthRepo
}

func NewAuthService(authRepo repo.IAuthRepo) IAuthService {
	return &authService{repo: authRepo}
}

func (as *authService) Ping() *gin.H {
	return &gin.H{
		"message": as.repo.GetPing(),
	}
}
