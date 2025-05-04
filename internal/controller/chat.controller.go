package controller

import "github.com/gin-gonic/gin"

type ChatController struct{}

func NewChatController() *ChatController {
	return &ChatController{}
}

func (cc *ChatController) SendMessage(c *gin.Context) {

}
