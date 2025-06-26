package websocket

import (
	"encoding/json"

	"github.com/dinhdev-nu/realtime_auth_go/internal/input"
)

type AckMessage struct {
	Event      string `json:"event"`
	ReceiverID int64  `json:"receiver_id"`
	Status     string `json:"status"`
	Content    []byte `json:"content"`
	MessageID  uint64 `json:"message_id"`
}

func NewAckMessage(status string, content []byte, receiverId int64, msgId uint64) AckMessage {
	return AckMessage{
		Event:      "ack",
		Status:     status,
		Content:    content,
		ReceiverID: receiverId,
		MessageID:  msgId,
	}
}

func NewMessageResponse(sender_id int64, reseiver_id int64, event string, content string) []byte {
	data, _ := json.Marshal(
		input.Message{
			Event:      event,
			SendID:     sender_id,
			ReceiverID: reseiver_id,
			Content:    content,
		})
	return data
}
