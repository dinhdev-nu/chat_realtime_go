package websocket

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
