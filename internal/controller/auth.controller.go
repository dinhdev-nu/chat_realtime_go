package controller

import (
	"fmt"

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

func (ac *AuthController) SignUp(c *gin.Context) {
	req, err := body.GetPayLoadFromRequestBody[input.SignUpInput](c)
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
	req, err := body.GetPayLoadFromRequestBody[input.LoginInput](c)
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
	req, err := body.GetPayLoadFromRequestBody[input.LogoutInput](c)
	if err != nil {
		res.BadRequestError(c, res.InvalidRequestPayloadCode, res.CodeMessage[res.InvalidRequestPayloadCode])
		return
	}
	req.UuidToken = c.GetString("uuidToken")
	fmt.Println("uuidToken", req.UuidToken)
	code := ac.AuthService.Logout(req.Email, req.UuidToken)
	if code != res.SuccessCode {
		res.BadRequestError(c, code, res.CodeMessage[code])
		return
	}
	res.SuccessResponse(c, nil)

}
