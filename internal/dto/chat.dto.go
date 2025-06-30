package dto

import (
	"time"

	"github.com/dinhdev-nu/realtime_auth_go/internal/database"
	"github.com/dinhdev-nu/realtime_auth_go/internal/model"
)

// API
// --- Init Chat struct ---
type InitChatOutPut struct {
	CurrentUser *model.GoDbUserInfo `json:"user"`
	Rooms       []RoomInitChat      `json:"rooms"`      // Danh sách phòng chat
	Followers   []int64             `json:"followers"`  // Danh sách người dùng theo dõi
	SocketUrl   string              `json:"socket_url"` // URL của WebSocket server
}
type RoomInitChat struct {
	RoomInfo database.GetPrivateRoomsByUserIdRow `json:"room"` // Thông tin phòng chat
	Users    InfoUserPrivateChat                 `json:"info"` // Danh sách người dùng trong phòng chat
}
type InfoUserPrivateChat struct {
	UserID     int64  `json:"user_id"`     // ID của người dùng
	UserName   string `json:"user_name"`   // Tên của người dùng
	UserAvatar string `json:"user_avatar"` // Avatar của người dùng
	UserStatus string `json:"user_status"` // Trạng thái của người dùng (online, offline, etc.)
}

type GetMessagesFromRoomOutput struct {
	RoomIsGroup    bool                                    `json:"room_is_group"`   // Kiểm tra phòng chat là nhóm hay không
	MessagesDriect []database.GetMessagesDirectByRoomIdRow `json:"messages_direct"` // Tin nhắn trong phòng chat
	MessagesGroup  []database.GoDbChatMessagesGroup        `json:"messages_group"`
}

type CreateRoomDTO struct {
	RoomID       uint64   `json:"room_id"` // Tên phòng chat
	RoomName     string   `json:"room_name" binding:"required"`
	RoomCreateBy int64    `json:"room_create_by" binding:"required"`
	RoomIsGroup  bool     `json:"room_is_group"`
	RoomMembers  []uint64 `json:"room_members" binding:"required"`
}

type UpdateStatusInput struct {
	RoomID uint64 `json:"room_id" binding:"required"` // ID của phòng chat
	UserId uint64 `json:"user_id" binding:"required"` // ID của người dùng
}

type SaveMessageDTO struct {
	MessageID         int64     `json:"message_id"`          // ID của tin nhắn
	MessageRoomID     uint64    `json:"message_room_id"`     // ID của phòng chat
	MessageSenderID   uint64    `json:"message_sender_id"`   // ID của người gửi
	MessageReceiverID uint64    `json:"message_receiver_id"` // ID của người nhận (nếu có)
	MessageContent    string    `json:"message_content"`     // Nội dung tin nhắn
	MessageSentAt     time.Time `json:"message_sent_at"`     // Thời gian gửi tin nhắn
}

// WebSocket
// --- Message struct ---
type Message struct {
	ID          uint64    `json:"id"`            // ID của tin nhắn
	RoomID      int64     `json:"room_id"`       // ID của phòng chat (nếu có)
	SendID      int64     `json:"sender_id"`     // ID của người gửi
	SendName    string    `json:"sender_name"`   // Tên của người gửi
	SendAvatar  string    `json:"sender_avatar"` // Avatar của người gửi
	ReceiverID  int64     `json:"receiver_id"`   // ID của người nhận (nếu có) 1 - 1 chat)
	ReceiverIDs []int64   `json:"receiver_ids"`  // ID của người nhận (nếu có) 1 - n chat
	SendAt      time.Time `json:"send_at"`       // Thời gian gửi tin nhắn
	Content     string    `json:"content"`       // Nội dung tin nhắn
	ContentType string    `json:"content_type"`  // Loại nội dung (text, image, video, file, etc.)
}

type StatusMessage struct {
	Status string `json:"status"` // Nội dung trạng thái (vd: "typing", "online", "offline")
}
type Typing struct {
	RoomID int64 `json:"room_id"` // ID của phòng chat (nếu có)
}
type Read struct {
	RoomID int64 `json:"room_id"` // ID của phòng chat (nếu có)
}

type OnMessage struct {
	Event       string        `json:"event"`                  // Loại tin nhắn (vd: "message", "status", "typing")
	Type        string        `json:"type,omitempty"`         // Loại tin nhắn (vd: "single", "multi", "broadcast")
	SendID      int64         `json:"sender_id,omitempty"`    // ID của người gửi
	ReceiverID  int64         `json:"receiver_id,omitempty"`  // ID của người nhận (nếu có) 1 - 1 chat
	ReceiverIDs []int64       `json:"receiver_ids,omitempty"` // ID của người nhận
	Message     Message       `json:"message,omitempty"`      // Tin nhắn
	Status      StatusMessage `json:"status,omitempty"`       // Tin nhắn trạng thái
	Typing      Typing        `json:"typing,omitempty"`       // Tin nhắn đang gõ
	Read        Read          `json:"read,omitempty"`         // Tin nhắn đã đọc
}
