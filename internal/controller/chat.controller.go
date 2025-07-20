package controller

import (
	"github.com/dinhdev-nu/realtime_auth_go/internal/dto"
	"github.com/dinhdev-nu/realtime_auth_go/internal/model"
	service "github.com/dinhdev-nu/realtime_auth_go/internal/service/chat"
	"github.com/dinhdev-nu/realtime_auth_go/internal/utils/body"
	"github.com/dinhdev-nu/realtime_auth_go/pkg/response"
	"github.com/gin-gonic/gin"
)

type ChatController struct {
	ChatService service.IChatService
}

func NewChatController(ChatService service.IChatService) *ChatController {
	return &ChatController{
		ChatService: ChatService,
	}
}

func (cc *ChatController) InitChat(c *gin.Context) {
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

	res, err := cc.ChatService.InitChat(userInfo)
	if err != nil {
		response.BadRequestError(c, response.ErrorCode, err.Error())
		return
	}
	response.SuccessResponse(c, res)
}

func (cc *ChatController) GetMessages(c *gin.Context) {
	roomId := c.Param("room-id")
	page := c.Query("page")
	offset := c.Query("offset")
	if roomId == "" {
		response.BadRequestError(c, response.InvalidRequestPayloadCode, "Room ID is required")
	}
	messages, err := cc.ChatService.GetMessagesFromRoom(roomId, page, offset)
	if err != nil {
		response.BadRequestError(c, response.ErrorCode, err.Error())
		return
	}
	response.SuccessResponse(c, messages)
}

func (cc *ChatController) CreateNewRoom(c *gin.Context) {
	data, err := body.GetPayLoadFromRequestBody[dto.CreateRoomDTO](c)
	if err != nil {
		response.BadRequestError(c, response.InvalidRequestPayloadCode, err.Error())
		return
	}

	newRoom, err := cc.ChatService.CreateRoomChat(data)
	if err != nil {
		response.BadRequestError(c, response.ErrorCode, err.Error())
		return
	}
	response.SuccessResponse(c, newRoom)
}

func (cc *ChatController) UpdateStatusMessages(c *gin.Context) {
	data, err := body.GetPayLoadFromRequestBody[dto.UpdateStatusInput](c)
	if err != nil {
		response.BadRequestError(c, response.InvalidRequestPayloadCode, err.Error())
		return
	}

	err = cc.ChatService.UpdateStatusMessages(data)
	if err != nil {
		response.BadRequestError(c, response.ErrorCode, err.Error())
		return
	}
	response.SuccessResponse(c, "Status updated successfully")
}

func (cc *ChatController) GetRooms(c *gin.Context) {
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
	res, err := cc.ChatService.GetGroupRoomsByUserID2(userInfo)
	if err != nil {
		response.BadRequestError(c, response.ErrorCode, err.Error())
		return
	}
	response.SuccessResponse(c, res)
}
