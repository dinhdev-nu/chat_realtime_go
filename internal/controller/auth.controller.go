package controller

import (
	service "github.com/dinhdev-nu/realtime_auth_go/internal/service/auth"
	res "github.com/dinhdev-nu/realtime_auth_go/pkg/response"
	"github.com/gin-gonic/gin"
)

type AuthController struct {
	AuthService service.IAuthService
}

func NewAuthController(authService service.IAuthService) *AuthController {
	return &AuthController{
		AuthService: authService,
	}
}

func (ac *AuthController) Ping(ctx *gin.Context) {
	res.SuccessResponse(ctx, ac.AuthService.Ping())
	// res.BadRequestError(ctx, 4001, "Error")
}
