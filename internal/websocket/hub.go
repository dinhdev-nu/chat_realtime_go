package websocket

import (
	"sync"
)

// --- Message struct ---
type Message struct {
	SendID      string   `json:"sender_id"`    // ID của người gửi
	ReceiverID  string   `json:"receiver_id"`  // ID của người nhận (nếu có) 1 - 1 chat
	ReceiverIDs []string `json:"receiver_ids"` // ID của người nhận (nếu có) group chat 1 - n chat
	Type        string   `json:"type"`         // singgle, multi, broadcast
	Data        []byte   `json:"data"`         // Nội dung tin nhắn
}
type RawData struct {
	SendID      string   `json:"sender_id"`    // ID của người gửi
	ReceiverID  string   `json:"receiver_id"`  // ID của người nhận (nếu có) 1 - 1 chat
	ReceiverIDs []string `json:"receiver_ids"` // ID của người nhận (nếu có) group
	Type        string   `json:"type"`         // singgle, multi, broadcast
	Data        string   `json:"data"`         // Nội dung tin nhắn
}

type Hub struct {
	Clients    map[string]*Client // danh sách client đang kết nối
	Register   chan *Client       // kênh thông báo client mới kết nối
	Unregister chan *Client       // kênh thông báo client đã ngắt kết nối
	Broadcast  chan Message       // kênh gửi tin nhắn đến tất cả client
	mu         sync.Mutex         // Mutex để đồng bộ hóa truy cập đến danh sách client
}

// Tạo 1 Hub mới ( trung tâm điều phối )
func NewHub() *Hub {
	return &Hub{
		Clients:    make(map[string]*Client),
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
		Broadcast:  make(chan Message),
	}
}

// Bắt đầu Hub ( chạy trong 1 goroutine )
func (h *Hub) Run() {
	for {
		select {
		case client := <-h.Register: // Case có client mới kết nối
			h.mu.Lock()
			h.Clients[client.UserID] = client // Thêm client vào danh sách client connecting
			h.mu.Unlock()

		case client := <-h.Unregister: // Case có client ngắt kết nối
			h.mu.Lock()
			if _, ok := h.Clients[client.UserID]; ok { // Kiểm tra client có trong danh sách không
				delete(h.Clients, client.UserID) // Xóa client khỏi danh sách
				close(client.Send)               // Đóng kênh gửi tin nhắn của client
			}
			h.mu.Unlock()
		case msg := <-h.Broadcast: // case có tin nhắn gửi đến tất cả client
			switch msg.Type { // Kiểm tra loại tin nhắn
			case "single": // Nếu là tin nhắn đơn (1 - 1 chat)
				if receiver, ok := h.Clients[msg.ReceiverID]; ok { // Kiểm tra client có trong danh sách không
					receiver.Send <- msg.Data // Gửi tin nhắn đến client nhận
				}
			case "multi": // Nếu là tin nhắn nhóm (1 - n chat)
				for _, id := range msg.ReceiverIDs { // Duyệt qua tất cả client trong danh sách
					if receiver, ok := h.Clients[id]; ok { // Kiểm tra client có trong danh sách không
						receiver.Send <- msg.Data // Gửi tin nhắn đến client nhận
					}
				}
			default:
				for id, client := range h.Clients { // Duyệt qua tất cả client trong danh sách
					if id != msg.SendID {
						client.Send <- msg.Data // Gửi tin nhắn đến client nhận
					}
				}
			}

		}
	}
}

/*
	paload :
	1 - 1 chat : { "receiver_id": "user_id", "type": "single", "data": "message" }
	1 - n chat : { "receiver_ids": ["user_id1", "user_id2"], "type": "group", "data": "message" }
	1 - all chat : { "type": "broadcast", "data": "message" }
*/
