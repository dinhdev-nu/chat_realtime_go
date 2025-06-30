package websocket

import (
	"encoding/json"

	"github.com/dinhdev-nu/realtime_auth_go/internal/dto"
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
		dto.OnMessage{
			Event:       event,
			Type:        "single",
			SendID:      sender_id,
			ReceiverIDs: []int64{reseiver_id},
			Status: dto.StatusMessage{
				Status: content,
			},
		})
	return data
}
