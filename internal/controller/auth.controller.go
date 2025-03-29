package controller

import (
	"github.com/dinhdev-nu/realtime_auth_go/internal/input"
	service "github.com/dinhdev-nu/realtime_auth_go/internal/service/auth"
	"github.com/dinhdev-nu/realtime_auth_go/internal/utils/body"
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

func (ac *AuthController) Ping(c *gin.Context) {
	// res.SuccessResponse(c, ac.AuthService.Ping())
	// res.BadRequestError(ctx, 4001, "Error")
}

func (ac *AuthController) Register(c *gin.Context) {
	req, err := body.GetPayLoadFromRequestBody[input.EmailInput](c)
	if err != nil {
		res.BadRequestError(c, res.InvalidRequestPayloadCode, res.CodeMessage[res.InvalidRequestPayloadCode])
		return
	}
	metadata, code := ac.AuthService.Register(req.Email)
	if code != res.SuccessCode {
		res.BadRequestError(c, code, res.CodeMessage[code])
		return

	}
	res.SuccessResponse(c, metadata)
}

func (ac *AuthController) SendOtp(c *gin.Context) {
	req, err := body.GetPayLoadFromRequestBody[input.EmailInput](c)
	if err != nil {
		res.BadRequestError(c, res.InvalidRequestPayloadCode, res.CodeMessage[res.InvalidRequestPayloadCode])
		return
	}
	code := ac.AuthService.SendOtp(req.Email)
	if code != res.SuccessCode {
		res.BadRequestError(c, code, res.CodeMessage[code])
		return
	}
	res.SuccessResponse(c, nil)
}

func (ac *AuthController) VerifyOtp(c *gin.Context) {
	req, err := body.GetPayLoadFromRequestBody[input.OtpInput](c)
	if err != nil {
		res.BadRequestError(c, res.InvalidRequestPayloadCode, res.CodeMessage[res.InvalidRequestPayloadCode])
		return
	}
	email, code := ac.AuthService.VeryfyOtp(req.Email, req.Otp)
	if code != res.SuccessCode {
		res.BadRequestError(c, code, res.CodeMessage[code])
		return
	}
	res.SuccessResponse(c, map[string]string{"keyPass": email})
}
