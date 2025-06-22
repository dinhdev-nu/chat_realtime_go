package websocket

import (
	"fmt"
	"net/http"
	"time"

	"github.com/dinhdev-nu/realtime_auth_go/global"
	"github.com/dinhdev-nu/realtime_auth_go/internal/utils"
	"github.com/dinhdev-nu/realtime_auth_go/pkg/response"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

// nhiệm vụ: xử lý các request từ client gửi đến server và upgrade kết nối từ http sang websocket

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true }, // cho phép tất cả các origin kết nối đến server
}

func HandleWebSocket(hub *Hub) gin.HandlerFunc {
	return func(c *gin.Context) {

		// check token
		token := c.Query("token")    // lấy token từ query string
		userID := c.Query("user_id") // lấy userID từ query string
		if userID == "" || token == "" {
			response.BadRequestError(c, response.UpgradeWebSocketErrorCode, "user_id is required") // nếu không có userID thì trả về lỗi
			return
		}
		c.Set("user_id", userID) // gán userID vào context để sử dụng sau này
		c.Set("token", token)    // gán token vào context để sử dụng sau này

		c.Next()

		conn, err := upgrader.Upgrade(c.Writer, c.Request, nil) // upgrade kết nối từ http sang websocket
		if err != nil {
			response.BadRequestError(c, response.UpgradeWebSocketErrorCode, "failed to upgrade connection") // nếu có lỗi thì trả về lỗi
			return
		}
		Client := &Client{
			UserID: utils.StringToInt64(userID), // gán userID cho client,
			Conn:   conn,                        // gán kết nối cho client
			Hub:    hub,                         // gán hub cho client
			Send:   make(chan []byte),           // khởi tạo kênh gửi tin nhắn đến client
		}
		hub.Register <- Client // thông báo hub có client mới kết nối

		go Client.ReadMessage()  // chạy goroutine đọc tin nhắn từ client
		go Client.WriteMessage() // chạy goroutine gửi tin nhắn đến client
	}
}

func (hub *Hub) handlePresence(userId int64, status string) {
	// save presence to cache
	if status == "online" {
		global.Rdb.Set(
			&gin.Context{},
			fmt.Sprintf("user:%d:presence", userId),
			status,
			time.Minute*2,
		)
	} else {
		global.Rdb.Del(
			&gin.Context{},
			fmt.Sprintf("user:%d:presence", userId),
		)
		global.Rdb.Set(
			&gin.Context{},
			fmt.Sprintf("user:%d:last_seen", userId),
			time.Now().Format("2006-01-02 15:04:05"),
			0,
		)
	}

}
