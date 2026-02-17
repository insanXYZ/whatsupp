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
	userId := claims.Sub

	user := new(entity.User)
	err := cs.userRepository.TakeById(ctx, user, userId)
	if err != nil {
		return err
	}

	ws, err := websocket.Upgrader.Upgrade(w, r, nil)
	if err != nil {
		return err
	}

	client := &websocket.Client{
		User:                   user,
		Hub:                    cs.hub,
		Conn:                   ws,
		Send:                   make(chan *dto.BroadcastMessageWs, 250),
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
	bc *dto.BroadcastMessageWs,
	hub *websocket.Hub,
) error {

	err := ms.messageRepository.Transaction(ctx, func(tx *gorm.DB) error {

		groupTx := ms.groupRepository.WithTx(tx)
		memberTx := ms.memberRepository.WithTx(tx)
		messageTx := ms.messageRepository.WithTx(tx)

		if bc.Request.GroupID == nil {

			group := new(entity.Group)
			err := groupTx.TakePersonalGroupBySenderAndReceiverId(ctx, bc.User.ID, *bc.Request.ReceiverID, group)
			if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
				return err
			}

			if err != nil {

				newGroup := &entity.Group{
					Bio:       "~",
					Name:      fmt.Sprintf("%s-%d", entity.PERSONAL, time.Now().Unix()),
					Image:     storage.DEFAULT_GROUP_PROFILE_PICTURE_URL,
					GroupType: entity.PERSONAL,
				}

				err = groupTx.Create(ctx, newGroup)
				if err != nil {
					return err
				}

				group = newGroup

				newMembers := []*entity.Member{
					{
						GroupID: newGroup.ID,
						UserID:  bc.User.ID,
						Role:    entity.MEMBER,
					},
					{
						GroupID: newGroup.ID,
						UserID:  *bc.Request.ReceiverID,
						Role:    entity.MEMBER,
					},
				}

				err := memberTx.Creates(ctx, newMembers)
				if err != nil {
					return err
				}

			}

			bc.Request.GroupID = &group.ID

		}

		member := new(entity.Member)
		if err := memberTx.TakeByUserIdAndGroupId(
			ctx,
			bc.User.ID,
			*bc.Request.GroupID,
			member,
		); err != nil {
			fmt.Println("errors TakeByUserIdAndGroupId")
			return err
		}

		newMessage := &entity.Message{
			MemberID: member.ID,
			Message:  bc.Request.Message,
		}

		if err := messageTx.Create(ctx, newMessage); err != nil {
			fmt.Println("errors create message")
			return err
		}

		newMessage.Member = &entity.Member{
			User: bc.User,
		}

		bc.Message = newMessage

		if !hub.IsExistGroup(*bc.Request.GroupID) {
			memberIds, err := memberTx.GetUserIdsWithGroupId(ctx, *bc.Request.GroupID)
			if err != nil {
				return err
			}

			hub.CreateGroup(*bc.Request.GroupID, memberIds)
		}

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
