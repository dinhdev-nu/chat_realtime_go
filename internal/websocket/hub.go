package websocket

import (
	"encoding/json"
	"sync"

	"github.com/dinhdev-nu/realtime_auth_go/internal/dto"
	"github.com/dinhdev-nu/realtime_auth_go/internal/utils"
)

type Hub struct {
	Clients   map[int64]*Client // danh sách client đang kết nối
	Following map[int64][]int64

	Register    chan *Client       // kênh thông báo client mới kết nối
	Unregister  chan *Client       // kênh thông báo client đã ngắt kết nối
	Broadcast   chan dto.OnMessage // kênh gửi tin nhắn đến tất cả client
	SubscribeTo chan dto.OnMessage // kênh để client đăng ký theo dõi người dùng khác
	Ack         chan AckMessage    // kênh gửi tin nhắn ack đến client
	mu          sync.Mutex         // Mutex để đồng bộ hóa truy cập đến danh sách client
}

// Tạo 1 Hub mới ( trung tâm điều phối )
func NewHub() *Hub {
	return &Hub{
		Clients:   make(map[int64]*Client),
		Following: make(map[int64][]int64),

		Register:    make(chan *Client),
		Unregister:  make(chan *Client),
		SubscribeTo: make(chan dto.OnMessage), // kênh để client đăng ký theo dõi người dùng khác
		Ack:         make(chan AckMessage),
		Broadcast:   make(chan dto.OnMessage),
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
			h.handlePresence(client.UserID, "online") // Gọi hàm handlePresence để lưu trạng thái online của client

		case client := <-h.Unregister: // Case có client ngắt kết nối
			if followers, ok := h.Following[client.UserID]; ok {
				for _, id := range followers { // Duyệt qua tất cả người dùng đang theo dõi client
					if follower, ok := h.Clients[id]; ok { // Kiểm tra client có trong danh sách không
						follower.Send <- NewMessageResponse(client.UserID, id, "status", "offline") // Gửi tin nhắn thông báo client đã offline
					}
				}
			}
			h.mu.Lock()
			if _, ok := h.Clients[client.UserID]; ok { // Kiểm tra client có trong danh sách không
				delete(h.Clients, client.UserID)   // Xóa client khỏi danh sách
				delete(h.Following, client.UserID) // Xóa client khỏi danh sách phòng chat
				close(client.Send)                 // Đóng kênh gửi tin nhắn của client
			}
			h.mu.Unlock()
			h.handlePresence(client.UserID, "offline") // Gọi hàm handlePresence để lưu trạng thái offline của client
		case followers := <-h.SubscribeTo:
			h.mu.Lock()
			if _, ok := h.Following[followers.SendID]; ok {
				followersID := h.Following[followers.SendID]
				for _, id := range followers.ReceiverIDs {
					if !utils.Contains(followersID, id) {
						h.Following[followers.SendID] = append(h.Following[followers.SendID], id)
					}
				}
			} else {
				h.Following[followers.SendID] = followers.ReceiverIDs
				for _, id := range followers.ReceiverIDs {
					if receiver, ok := h.Clients[id]; ok { // Kiểm tra client có trong danh sách không
						receiver.Send <- NewMessageResponse(followers.SendID, id, "status", "online") // Gửi tin nhắn thông báo đăng ký thành công
					}
				}
			}
			h.mu.Unlock()
		case ack := <-h.Ack: // Case có tin nhắn ack gửi đến client
			if receiver, ok := h.Clients[ack.ReceiverID]; ok {
				byteMsg, _ := json.Marshal(ack) // Chuyển đổi tin nhắn ack sang định dạng JSON
				receiver.Send <- byteMsg        // Gửi tin nhắn ack đến client nhận
			}
		case msg := <-h.Broadcast: // case có tin nhắn gửi đến tất cả client
			byteMsg, _ := json.Marshal(msg) // Chuyển đổi tin nhắn sang định dạng JSON
			switch msg.Type {               // Kiểm tra loại tin nhắn
			case "single": // Nếu là tin nhắn đơn (1 - 1 chat)
				if receiver, ok := h.Clients[msg.ReceiverID]; ok { // Kiểm tra client có trong danh sách không
					receiver.Send <- byteMsg // Gửi tin nhắn đến client nhận
				}
			case "group": // Nếu là tin nhắn nhóm (1 - n chat)
				for _, id := range msg.ReceiverIDs { // Duyệt qua tất cả client trong danh sách
					if receiver, ok := h.Clients[id]; ok { // Kiểm tra client có trong danh sách không
						receiver.Send <- byteMsg // Gửi tin nhắn đến client nhận
					}
				}
			default:
				for id, client := range h.Clients { // Duyệt qua tất cả client trong danh sách
					if id != msg.SendID {
						client.Send <- byteMsg // Gửi tin nhắn đến client nhận
					}
				}
			}
		}
	}
}

/*
	paload :
	1 - 1 chat : { "receiver_id": "user_id", "type": "single", "data": "message" , send_at, room_id}
	1 - n chat : { "receiver_ids": ["user_id1", "user_id2"], "type": "group", "data": "message" }
	1 - all chat : { "type": "broadcast", "data": "message" }
*/
