package chat

import (
	c "github.com/dinhdev-nu/realtime_auth_go/internal/controller"
	"github.com/dinhdev-nu/realtime_auth_go/internal/utils/middleware/auth"
	"github.com/dinhdev-nu/realtime_auth_go/internal/websocket"
	"github.com/gin-gonic/gin"
)

type ChatRouter struct {
	cc *c.ChatController
}

func NewChatRouter(cc *c.ChatController) *ChatRouter {
	return &ChatRouter{
		cc: cc,
	}
}

func (cr *ChatRouter) InitRoutes(router *gin.RouterGroup) {

	hub := websocket.NewHub() // tạo 1 hub mới
	go hub.Run()              // chạy hub trong 1 goroutine

	chatRouter := router.Group("chat")
	{
		// websocket endpoint
		chatRouter.GET("/ws", websocket.HandleWebSocket(hub), auth.AuthMiddleware()) // upgrade http to websocket

		// middleware
		chatRouter.Use(auth.AuthMiddleware())
		// init chat
		chatRouter.GET("/init", cr.cc.InitChat) // Get info chat page

		// api
		chatRouter.GET("/get-messages/:room-id", cr.cc.GetMessages) // ChatController.GetMessages)

		// room
		chatRouter.GET("/get-room/:room_id", cr.cc.GetRoomChatById) // ChatController.GetRoomChatById)
		chatRouter.POST("/create-room", cr.cc.CreateNewRoom)

		// status
		chatRouter.POST("/set-status", cr.cc.UpdateStatusMessages) // ChatController.UpdateStatus
	}
}
