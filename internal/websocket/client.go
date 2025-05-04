package websocket

import (
	"encoding/json"
	"fmt"

	"github.com/gorilla/websocket"
)

type Client struct {
	UserID string          // ID của người dùng
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
		var message Message
		if err := json.Unmarshal(msg, &message); err != nil { // Giải mã tin nhắn từ client

			// case data là plain text // client go btoa data to byte
			var rawData RawData
			json.Unmarshal(msg, &rawData) // Giải mã tin nhắn từ client
			message.SendID = rawData.SendID
			message.ReceiverID = rawData.ReceiverID
			message.ReceiverIDs = rawData.ReceiverIDs
			message.Type = rawData.Type
			message.Data = []byte(rawData.Data) // Chuyển đổi nội dung tin nhắn từ string sang []byte
		}

		message.SendID = c.UserID // Gán ID của người gửi vào tin nhắn
		message.Data = msg        // Gán nội dung tin nhắn vào tin nhắn

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
