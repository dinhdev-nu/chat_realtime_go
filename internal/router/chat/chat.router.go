package chat

import (
	"github.com/dinhdev-nu/realtime_auth_go/internal/utils/middleware/auth"
	"github.com/dinhdev-nu/realtime_auth_go/internal/websocket"
	"github.com/gin-gonic/gin"
)

type ChatRouter struct {
	// ChatController *c.ChatController
}

func NewChatRouter() *ChatRouter {
	return &ChatRouter{
		// ChatController: cc,
	}
}

func (cr *ChatRouter) InitRoutes(router *gin.RouterGroup) {

	hub := websocket.NewHub() // tạo 1 hub mới
	go hub.Run()              // chạy hub trong 1 goroutine

	chatRouter := router.Group("chat")
	{
		// api
		chatRouter.POST("/send-message", nil)      // ChatController.SendMessage)
		chatRouter.POST("/get-messages", nil)      // ChatController.GetMessages)
		chatRouter.POST("/get-conversations", nil) // ChatController.GetConversations)

		// websocket endpoint
		chatRouter.GET("/ws", websocket.HandleWebSocket(hub), auth.AuthMiddleware()) // upgrade http to websocket
	}
}
