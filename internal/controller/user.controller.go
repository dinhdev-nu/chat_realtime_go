package controller

import (
	"github.com/dinhdev-nu/realtime_auth_go/internal/model"
	service "github.com/dinhdev-nu/realtime_auth_go/internal/service/user"
	"github.com/dinhdev-nu/realtime_auth_go/pkg/response"
	"github.com/gin-gonic/gin"
)

type UserController struct {
	us service.IUserService
}

func NewUserController(s service.IUserService) *UserController {
	return &UserController{
		us: s,
	}
}

func (uc *UserController) GetUserInfoByName(c *gin.Context) {
	userName := c.Param("username")
	if userName == "" {
		response.BadRequestError(c, response.ErrorCodeInvalidRequest, "Username is required")
		return
	}
	res, err := uc.us.GetUserInfoByName(userName)
	if err != nil {
		response.BadRequestError(c, response.ErrorUserNotExist, "Failed to get user info")
		return
	}
	response.SuccessResponse(c, res)
}

func (uc *UserController) SearchUserByName(c *gin.Context) {
	userName := c.Query("username")
	value, exists := c.Get("user")
	if !exists {
		response.BadRequestError(c, response.InvalidRequestPayloadCode, "User ID is required")
		return
	}
	userInfo, ok := value.(*model.GoDbUserInfo)
	if !ok {
		response.BadRequestError(c, response.InvalidRequestPayloadCode, "Invalid user information")
		return
	}

	res, err := uc.us.SearchUserByName(userName, userInfo.UserID)
	if err != nil {
		response.BadRequestError(c, response.ErrorUserNotExist, "Failed to search user")
		return
	}
	response.SuccessResponse(c, res)
}
