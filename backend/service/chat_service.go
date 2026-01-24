package service

import (
	"context"
	"net/http"
	"whatsupp-backend/repository"
	"whatsupp-backend/websocket"

	"github.com/go-playground/validator/v10"
)

type ChatService struct {
	validator                   *validator.Validate
	groupRepository             *repository.GroupRepository
	groupMemberRepository       *repository.GroupMemberRepository
	messageRepository           *repository.MessageRepository
	messageAttachmentRepository *repository.MessageAttachmentRepository
	hub                         *websocket.Hub
}

func NewChatService(
	groupRepository *repository.GroupRepository,
	groupMemberRepository *repository.GroupMemberRepository,
	messageRepository *repository.MessageRepository,
	messageAttachmentRepository *repository.MessageAttachmentRepository,
	validator *validator.Validate,
	hub *websocket.Hub,
) *ChatService {
	return &ChatService{
		hub:                         hub,
		validator:                   validator,
		groupRepository:             groupRepository,
		groupMemberRepository:       groupMemberRepository,
		messageRepository:           messageRepository,
		messageAttachmentRepository: messageAttachmentRepository,
	}
}

func (cs *ChatService) HandleUpgradeWs(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	ws, err := websocket.Upgrader.Upgrade(w, r, nil)
	if err != nil {
		return err
	}

	client := &websocket.Client{
		Hub:  cs.hub,
		Conn: ws,
		Send: make(chan []byte, 256),
	}

	client.Hub.Register(client)

	go client.WritePump()
	go client.WritePump()

	return nil
}
