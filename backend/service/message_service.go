package service

import (
	"context"
	"fmt"
	"mime/multipart"
	"net/http"
	"path/filepath"
	"time"
	"whatsupp-backend/dto"
	"whatsupp-backend/entity"
	"whatsupp-backend/repository"
	"whatsupp-backend/storage"
	"whatsupp-backend/util"
	"whatsupp-backend/websocket"

	"github.com/go-playground/validator/v10"
	storage_go "github.com/supabase-community/storage-go"
	"gorm.io/gorm"
)

type MessageService struct {
	validator                   *validator.Validate
	groupRepository             *repository.GroupRepository
	memberRepository            *repository.MemberRepository
	messageRepository           *repository.MessageRepository
	messageAttachmentRepository *repository.MessageAttachmentRepository
	userRepository              *repository.UserRepository
	hub                         *websocket.Hub
	storage                     *storage_go.Client
}

func NewMessageService(
	validator *validator.Validate,
	groupRepository *repository.GroupRepository,
	memberRepository *repository.MemberRepository,
	messageRepository *repository.MessageRepository,
	messageAttachmentRepository *repository.MessageAttachmentRepository,
	userRepository *repository.UserRepository,
	hub *websocket.Hub,
	storage *storage_go.Client,
) *MessageService {
	return &MessageService{
		userRepository:              userRepository,
		storage:                     storage,
		hub:                         hub,
		validator:                   validator,
		groupRepository:             groupRepository,
		memberRepository:            memberRepository,
		messageRepository:           messageRepository,
		messageAttachmentRepository: messageAttachmentRepository,
	}
}

func (cs *MessageService) HandleUpgradeWs(ctx context.Context, claims *util.Claims, w http.ResponseWriter, r *http.Request) error {
	ws, err := websocket.Upgrader.Upgrade(w, r, nil)
	if err != nil {
		return err
	}

	userId := claims.Sub

	client := &websocket.Client{
		Id:                     userId,
		Hub:                    cs.hub,
		Conn:                   ws,
		Send:                   make(chan *dto.BroadcastMessageWS, 250),
		HandlerIncomingMessage: cs.handleIncomingMessage,
	}

	client.Hub.Register(client)

	go client.ReadPump()
	go client.WritePump()

	return nil
}

func (cs *MessageService) HandleUploadFileAttachments(ctx context.Context, messageID string, files []*multipart.FileHeader) error {

	return cs.messageAttachmentRepository.Transaction(ctx, func(tx *gorm.DB) error {

		messageRepoTx := cs.messageRepository.WithTx(tx)
		err := messageRepoTx.TakeById(ctx, &entity.Message{}, messageID)
		if err != nil {
			return err
		}

		messageAttTx := cs.messageAttachmentRepository.WithTx(tx)

		clientStorage := cs.storage

		for _, file := range files {

			ext := filepath.Ext(file.Filename)

			fileName := fmt.Sprintf("whatsupp-%s-%v%s",
				messageID,
				time.Now().Unix(),
				ext,
			)

			open, err := file.Open()
			if err != nil {
				return err
			}

			_, err = clientStorage.UploadFile(storage.FILE_ATTACHMENT_BUCKET, fileName, open)
			if err != nil {
				return err
			}

			publicUrl := clientStorage.GetPublicUrl(storage.FILE_ATTACHMENT_BUCKET, fileName)

			newMessageAtt := &entity.MessageAttachment{
				MessageID: messageID,
				FileExt:   ext,
				FileURL:   publicUrl.SignedURL,
			}

			err = messageAttTx.Create(ctx, newMessageAtt)
			if err != nil {
				return err
			}

		}

		return nil

	})
}

func (ms *MessageService) handleIncomingMessage(
	ctx context.Context,
	msg *dto.BroadcastMessageWS,
	hub *websocket.Hub,
) error {

	err := ms.messageRepository.Transaction(ctx, func(tx *gorm.DB) error {

		groupTx := ms.groupRepository.WithTx(tx)
		memberTx := ms.memberRepository.WithTx(tx)
		messageTx := ms.messageRepository.WithTx(tx)

		if msg.GroupID == nil {

			group, err := groupTx.TakeOrCreatePersonalGroup(
				ctx,
				msg.ClientID,
				*msg.ReceiverID,
			)
			if err != nil {
				return err
			}

			msg.GroupID = &group.ID
		}

		member := new(entity.Member)
		if err := memberTx.TakeByUserIdAndGroupId(
			ctx,
			msg.ClientID,
			*msg.GroupID,
			member,
		); err != nil {
			return err
		}

		newMessage := &entity.Message{
			MemberID: member.ID,
			Message:  msg.Message,
		}

		if err := messageTx.Create(ctx, newMessage); err != nil {
			return err
		}

		msg.MessageID = newMessage.ID

		return nil
	})

	if err != nil {
		return err
	}

	if !hub.IsExistGroup(*msg.GroupID) {

		memberIds, err := ms.memberRepository.GetUserIdsWithGroupId(ctx, *msg.GroupID)
		if err != nil {
			return err
		}

		hub.CreateGroup(*msg.GroupID, memberIds)
	}

	return nil
}

func (ms *MessageService) HandleGetMessages(ctx context.Context, groupId int, claims *util.Claims) ([]*entity.Message, error) {
	_, err := ms.groupRepository.TakeGroupWithGroupIdAndUserId(ctx, groupId, claims.Sub)
	if err != nil {
		return nil, err
	}

	return ms.messageRepository.GetMessages(ctx, groupId)
}
