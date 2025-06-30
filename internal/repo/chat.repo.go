package repo

import (
	"fmt"

	"github.com/dinhdev-nu/realtime_auth_go/global"
	"github.com/dinhdev-nu/realtime_auth_go/internal/database"
	"github.com/dinhdev-nu/realtime_auth_go/internal/dto"
	"github.com/gin-gonic/gin"
)

type IChatRepo interface {
	GetRoomsByUserId(userId int64) ([]database.GetPrivateRoomsByUserIdRow, error)
	GetRoomById(id uint64) (database.GoDbChatRoom, error)
	GetRoomByName(name string) (database.GetRoomByNameRow, error)
	CreateRoom(data *dto.CreateRoomDTO) error
	AddMembersToRoom(roomId uint64, users []uint64) error
	GetAnotherUserID(roomId uint64, userId int64) (int64, error)

	GetMessagesFromRoom(roomId uint64, page int64) ([]database.GetMessagesDirectByRoomIdRow, error)
	GetMessagesGruopFromRoom(roomId uint64, page int64) ([]database.GoDbChatMessagesGroup, error)
	SaveMessegeDirect(data dto.SaveMessageDTO) (int64, error)
	SaveMessegeGroup(data dto.SaveMessageDTO) (int64, error)
	SaveMessageStatus(msgId uint64, userId int64) error

	UpdateMessageStatus(data *dto.UpdateStatusInput) error
}

type chatRepo struct {
	sqlc *database.Queries
	ctx  *gin.Context
}

func NewChatRepo() IChatRepo {
	ctx := gin.Context{}
	return &chatRepo{
		sqlc: database.New(global.Mdbc),
		ctx:  &ctx,
	}
}

func (r *chatRepo) UpdateMessageStatus(data *dto.UpdateStatusInput) error {
	err := r.sqlc.UpdateMessageStatus(r.ctx, database.UpdateMessageStatusParams{
		RoomID:            data.RoomID,
		MessageReceiverID: data.UserId,
	})
	if err != nil {
		return err
	}
	return nil
}

func (r *chatRepo) SaveMessageStatus(msgId uint64, userId int64) error {
	err := r.sqlc.SaveMessageStatus(r.ctx, database.SaveMessageStatusParams{
		MessageID:     msgId,
		MessageUserID: uint64(userId),
	})
	if err != nil {
		return err
	}
	return nil
}

func (r *chatRepo) SaveMessegeDirect(data dto.SaveMessageDTO) (int64, error) {
	message := database.SaveMessageDirectParams{
		MessageRoomID:     data.MessageRoomID,
		MessageReceiverID: uint64(data.MessageReceiverID),
		MessageContent:    data.MessageContent,
		MessageType: database.NullGoDbChatMessagesDirectMessageType{
			Valid:                             true,
			GoDbChatMessagesDirectMessageType: database.GoDbChatMessagesDirectMessageTypeText,
		},
		MessageSentAt: data.MessageSentAt,
	}
	res, err := r.sqlc.SaveMessageDirect(r.ctx, message)
	if err != nil {
		return 0, err
	}
	id, _ := res.LastInsertId()
	return id, nil
}

func (r *chatRepo) SaveMessegeGroup(data dto.SaveMessageDTO) (int64, error) {
	message := database.SaveMessageGroupParams{
		MessageRoomID:   data.MessageRoomID,
		MessageSenderID: data.MessageSenderID,
		MessageContent:  data.MessageContent,
		MessageType: database.NullGoDbChatMessagesGroupMessageType{
			Valid:                            true,
			GoDbChatMessagesGroupMessageType: database.GoDbChatMessagesGroupMessageTypeText,
		},
		MessageSentAt: data.MessageSentAt,
	}
	res, err := r.sqlc.SaveMessageGroup(r.ctx, message)
	if err != nil {
		return 0, err
	}
	id, _ := res.LastInsertId()
	return id, nil
}

func (r *chatRepo) GetMessagesGruopFromRoom(roomId uint64, page int64) ([]database.GoDbChatMessagesGroup, error) {
	limit := 10
	messages, err := r.sqlc.GetMessagesGroupByRoomId(r.ctx, database.GetMessagesGroupByRoomIdParams{
		MessageRoomID: roomId,
		Limit:         int32(limit),
		Offset:        int32((page - 1) * int64(limit)),
	})
	if err != nil {
		return nil, err
	}
	return messages, nil
}

func (r *chatRepo) GetMessagesFromRoom(roomId uint64, page int64) ([]database.GetMessagesDirectByRoomIdRow, error) {
	limit := 10
	messages, err := r.sqlc.GetMessagesDirectByRoomId(r.ctx, database.GetMessagesDirectByRoomIdParams{
		MessageRoomID: roomId,
		Limit:         int32(limit),
		Offset:        int32((page - 1) * int64(limit)),
	})
	if err != nil {
		fmt.Println("Error fetching messages:", err)
		return nil, err
	}
	return messages, nil
}

func (r *chatRepo) GetAnotherUserID(roomId uint64, userId int64) (int64, error) {
	anotherUser, err := r.sqlc.GetAnotherPrivateMenberByRoomId(r.ctx, database.GetAnotherPrivateMenberByRoomIdParams{
		RoomID:       roomId,
		MemberUserID: uint64(userId),
	})
	if err != nil {
		return 0, err
	}
	return int64(anotherUser), nil
}

func (r *chatRepo) GetRoomsByUserId(userId int64) ([]database.GetPrivateRoomsByUserIdRow, error) {
	rooms, err := r.sqlc.GetPrivateRoomsByUserId(r.ctx, uint64(userId))
	if err != nil {
		return nil, err
	}
	return rooms, nil
}

func (r *chatRepo) GetRoomById(id uint64) (database.GoDbChatRoom, error) {
	room, err := r.sqlc.GetRoomById(r.ctx, id)
	if err != nil {
		return database.GoDbChatRoom{}, err
	}
	return room, nil
}

func (r *chatRepo) GetRoomByName(name string) (database.GetRoomByNameRow, error) {
	roomName := GetNameRoom(name)
	room, err := r.sqlc.GetRoomByName(r.ctx, database.GetRoomByNameParams{
		RoomName:   NullString(roomName[0]),
		RoomName_2: NullString(roomName[1]),
	})
	if err != nil {
		return database.GetRoomByNameRow{}, err
	}
	return room, nil
}

func (r *chatRepo) CreateRoom(data *dto.CreateRoomDTO) error {
	res, err := r.sqlc.CreateRoom(r.ctx, database.CreateRoomParams{
		RoomName:      NullString(data.RoomName),
		RoomCreatedBy: NullInt64(data.RoomCreateBy),
		RoomIsGroup:   data.RoomIsGroup,
	})
	if err != nil {
		fmt.Println("Error creating room:", err)
		return err
	}
	id, _ := res.LastInsertId()
	data.RoomID = uint64(id)
	return nil
}

func (r *chatRepo) AddMembersToRoom(roomId uint64, users []uint64) error {
	for _, user := range users {
		err := r.sqlc.InsetMemberToRoom(r.ctx, database.InsetMemberToRoomParams{
			RoomID:       roomId,
			MemberUserID: user,
		})
		if err != nil {
			return err
		}
	}
	return nil
}
