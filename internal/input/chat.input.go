package input

import "time"

// --- Message struct ---
type Message struct {
	ID          uint64    `json:"id"`            // ID của tin nhắn
	RoomID      string    `json:"room_id"`       // ID của phòng chat (nếu có)
	SendID      int64     `json:"sender_id"`     // ID của người gửi
	SendName    string    `json:"sender_name"`   // Tên của người gửi
	SendAvatar  string    `json:"sender_avatar"` // Avatar của người gửi
	ReceiverID  int64     `json:"receiver_id"`   // ID của người nhận (nếu có) 1 - 1 chat
	ReceiverIDs []int64   `json:"receiver_ids"`  // ID của người nhận (nếu có) group chat 1 - n chat
	Type        string    `json:"type"`          // singgle, multi, broadcast
	Event       string    `json:"event"`         // Sự kiện (vd: "typing", "join", "leave")
	SendAt      time.Time `json:"send_at"`       // Thời gian gửi tin nhắn
	Content     string    `json:"content"`       // Nội dung tin nhắn
	ContentType string    `json:"content_type"`  // Loại nội dung (text, image, video, file, etc.)
}

type CreateRoomInput struct {
	RoomName     string  `json:"room_name" binding:"required"`
	RoomCreateBy int64   `json:"room_create_by" binding:"required"`
	RoomIsGroup  bool    `json:"room_is_group"`
	RoomMembers  []int64 `json:"room_members" binding:"required"`
}

type UpdateStatusInput struct {
	RoomID uint64 `json:"room_id" binding:"required"` // ID của phòng chat
	UserId int64  `json:"user_id" binding:"required"` // ID của người dùng
}
