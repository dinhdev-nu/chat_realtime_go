package service

import (
	"github.com/dinhdev-nu/realtime_auth_go/internal/repo"
	"github.com/gin-gonic/gin"
)

type AuthService struct{
	AuthRepo *repo.AuthRepo
}

func NewAuthServiceContructer(authRepo *repo.AuthRepo) *AuthService { 
	return &AuthService{AuthRepo: authRepo}
}

func DefaultAuthService() *AuthService {
	return NewAuthServiceContructer(repo.NewAuthRepo())
}

func (as *AuthService) Ping() *gin.H {
	return &gin.H{
		"message": as.AuthRepo.GetPing(),
	}
}
