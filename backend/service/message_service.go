package service

import (
	"context"
	"errors"
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

// returning messageId, groupId, error
func (cs *MessageService) handleIncomingMessage(ctx context.Context, msg *dto.BroadcastMessageWS, hub *websocket.Hub) error {
	err := cs.messageRepository.Transaction(ctx, func(tx *gorm.DB) error {

		groupTx := cs.groupRepository.WithTx(tx)
		userTx := cs.userRepository.WithTx(tx)
		memberTx := cs.memberRepository.WithTx(tx)
		messageTx := cs.messageRepository.WithTx(tx)

		isNewMessage := msg.GroupID == nil

		if isNewMessage {

			group := new(entity.Group)

			err := groupTx.TakePersonalGroupBySenderAndReceiverId(ctx, msg.ClientID, *msg.ReceiverID, group)
			if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
				return err
			}

			if errors.Is(err, gorm.ErrRecordNotFound) {

				sender, receiver := new(entity.User), new(entity.User)

				if err := userTx.TakeById(ctx, sender, msg.ClientID); err != nil {
					return err
				}

				if err := userTx.TakeById(ctx, receiver, msg.ReceiverID); err != nil {
					return err
				}

				newGroup := &entity.Group{
					Name:      fmt.Sprintf("PRIVATE-%d", time.Now().Unix()),
					Bio:       "~",
					GroupType: entity.PERSONAL,
					Image:     storage.DEFAULT_GROUP_PROFILE_PICTURE_URL,
				}

				err = groupTx.Create(ctx, newGroup)
				if err != nil {
					return err
				}

				members := []*entity.Member{
					{
						Role:    entity.MEMBER,
						UserID:  sender.ID,
						GroupID: newGroup.ID,
					},
					{
						Role:    entity.MEMBER,
						UserID:  receiver.ID,
						GroupID: newGroup.ID,
					},
				}

				err := memberTx.Creates(ctx, members)
				if err != nil {
					return err
				}

				group = newGroup

				msg.GroupID = &group.ID

				hub.CreateGroup(*msg.GroupID, []int{msg.ClientID, *msg.GroupID})
			}

		}

		member := new(entity.Member)
		err := memberTx.TakeByUserIdAndGroupId(ctx, msg.ClientID, *msg.GroupID, member)
		if err != nil {
			return err
		}

		newMessage := &entity.Message{
			MemberID: member.ID,
			Message:  msg.Message,
		}

		err = messageTx.Create(ctx, newMessage)
		if err != nil {
			return err
		}

		msg.MessageID = newMessage.ID

		return nil
	})

	return err
}

func (ms *MessageService) HandleGetMessages(ctx context.Context, groupId int, claims *util.Claims) ([]*entity.Message, error) {
	_, err := ms.groupRepository.TakeGroupWithGroupIdAndUserId(ctx, groupId, claims.Sub)
	if err != nil {
		return nil, err
	}

	return ms.messageRepository.GetMessages(ctx, groupId)
}
