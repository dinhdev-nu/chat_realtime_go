package websocket

import (
	"encoding/json"
	"fmt"

	"github.com/dinhdev-nu/realtime_auth_go/internal/dto"
	"github.com/dinhdev-nu/realtime_auth_go/internal/repo"
	service "github.com/dinhdev-nu/realtime_auth_go/internal/service/chat"
	"github.com/gorilla/websocket"
)

type Client struct {
	UserID int64           // ID của người dùng
	Conn   *websocket.Conn // kết nối websocket
	Send   chan []byte     // kênh gửi tin nhắn đến client
	Hub    *Hub            // Hub kết nối của client
}

func (c *Client) ReadMessage() {
	defer func() {
		c.Hub.Unregister <- c // Thông báo Hub client đã ngắt kết nối
		c.Conn.Close()        // Đóng kết nối
	}()

	for {
		_, msg, err := c.Conn.ReadMessage() // Đọc tin nhắn từ client // _ là kiểu tin nhắn (text, binary, ping, pong)
		if err != nil {
			break // Nếu có lỗi thì ngắt kết nối
		}

		// Đưu tin nhắn đến broadcast phát cho tất cả client
		var message dto.OnMessage
		if err := json.Unmarshal(msg, &message); err != nil { // Giải mã tin nhắn từ client
			fmt.Println("Error unmarshalling message: ", err) // In ra lỗi nếu có
			newAck := NewAckMessage("error", msg, c.UserID, 0)
			c.Hub.Ack <- newAck
			continue // Tiếp tục vòng lặp nếu có lỗi
		}

		// Handle Message
		// Check Event
		switch message.Event {
		case "message":
			// Xử lý tin nhắn
			data, err := service.NewChatService(repo.NewChatRepo(), repo.NewAuthRepo(), repo.NewUserRepo()).HandleSendMesage(message)
			if err != nil {
				fmt.Println("Error handling message: ", err) // In ra lỗi nếu có
				newAck := NewAckMessage("error", msg, c.UserID, 0)
				c.Hub.Ack <- newAck // Gửi tin nhắn ack lỗi đến client

				continue // Tiếp tục vòng lặp nếu có lỗi
			}
			message.Message.ID = uint64(data.MessageID) // Lấy ID của tin nhắn từ dữ liệu trả về
			var ack = NewAckMessage("success", msg, message.SendID, message.Message.ID)
			c.Hub.Ack <- ack // Gửi tin nhắn ack thành công đến client
		case "status":
			message.ReceiverIDs = c.Hub.Following[message.SendID]               // Lấy danh sách người dùng theo dõi từ Hub
			fmt.Printf("User %d %s ...", message.SendID, message.Status.Status) // In ra thông báo trạng thái người dùng đã thay đổi
		case "subscribe":
			c.Hub.SubscribeTo <- message
			continue
		case "typing":
			fmt.Println("User is typing...") // In ra thông báo người dùng đang gõ
		case "read":
			fmt.Println("User has read the message") // In ra thông báo người dùng đã đọc tin nhắn
		}

		c.Hub.Broadcast <- message
	}
}

func (c *Client) WriteMessage() {
	defer c.Conn.Close() // Đóng kết nối khi kết thúc

	for msg := range c.Send { // Duyệt qua tất cả tin nhắn trong kênh gửi tin nhắn đến client
		if err := c.Conn.WriteMessage(websocket.TextMessage, msg); err != nil { // Gửi tin nhắn đến client
			fmt.Println("Error writing message: ", err) // In ra lỗi nếu có
			break                                       // Nếu có lỗi thì ngắt kết nối
		}
	}
}

// note payload from client to server need to convert data to []byte btoa(data)
