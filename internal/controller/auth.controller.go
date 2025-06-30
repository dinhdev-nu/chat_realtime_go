package controller

import (
	"github.com/dinhdev-nu/realtime_auth_go/internal/dto"
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
func (ac *AuthController) Register(c *gin.Context) {
	req, err := body.GetPayLoadFromRequestBody[dto.RegisterDTO](c)
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
	req, err := body.GetPayLoadFromRequestBody[dto.SendOtpDTO](c)
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
	req, err := body.GetPayLoadFromRequestBody[dto.VerifyOtpDTO](c)
	if err != nil {
		res.BadRequestError(c, res.InvalidRequestPayloadCode, res.CodeMessage[res.InvalidRequestPayloadCode])
		return
	}
	code := ac.AuthService.VeryfyOtp(req.Email, req.Otp)
	if code != res.SuccessCode {
		res.BadRequestError(c, code, res.CodeMessage[code])
		return
	}
	res.SuccessResponse(c, nil)
}

func (ac *AuthController) SignUp(c *gin.Context) {
	req, err := body.GetPayLoadFromRequestBody[dto.SignUpInput](c)
	if err != nil {
		res.BadRequestError(c, res.InvalidRequestPayloadCode, res.CodeMessage[res.InvalidRequestPayloadCode])
		return
	}
	code := ac.AuthService.SignUp(req.Email, req.Password)
	if code != res.SuccessCode {
		res.BadRequestError(c, code, res.CodeMessage[code])
		return
	}
	res.SuccessResponse(c, nil)
}

func (ac *AuthController) Login(c *gin.Context) {
	req, err := body.GetPayLoadFromRequestBody[dto.LoginInput](c)
	if err != nil {
		res.BadRequestError(c, res.InvalidRequestPayloadCode, res.CodeMessage[res.InvalidRequestPayloadCode])
		return
	}
	req.LoginIp = c.ClientIP()
	data, code := ac.AuthService.Login(req.Email, req.Password, req.LoginIp)
	if code != res.SuccessCode {
		res.BadRequestError(c, code, res.CodeMessage[code])
		return
	}
	res.SuccessResponse(c, data)
}

func (ac *AuthController) Logout(c *gin.Context) {
	req, err := body.GetPayLoadFromRequestBody[dto.LogoutInput](c)
	if err != nil {
		res.BadRequestError(c, res.InvalidRequestPayloadCode, res.CodeMessage[res.InvalidRequestPayloadCode])
		return
	}
	req.UuidToken = c.GetString("uuidToken")
	code := ac.AuthService.Logout(req.Email, req.UuidToken)
	if code != res.SuccessCode {
		res.BadRequestError(c, code, res.CodeMessage[code])
		return
	}
	res.SuccessResponse(c, nil)

}

func (ac *AuthController) DelOtp(c *gin.Context) {
	req, err := body.GetPayLoadFromRequestBody[dto.EmailInput](c)
	if err != nil {
		res.BadRequestError(c, res.InvalidRequestPayloadCode, res.CodeMessage[res.InvalidRequestPayloadCode])
		return
	}

	code := ac.AuthService.DelOtp(req.Email)
	if code != res.SuccessCode {
		res.BadRequestError(c, code, res.CodeMessage[code])
		return
	}
	res.SuccessResponse(c, nil)
}
