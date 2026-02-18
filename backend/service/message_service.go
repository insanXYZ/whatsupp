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
	conversationRepository      *repository.ConversationRepository
	memberRepository            *repository.MemberRepository
	messageRepository           *repository.MessageRepository
	messageAttachmentRepository *repository.MessageAttachmentRepository
	userRepository              *repository.UserRepository
	hub                         *websocket.Hub
	storage                     *storage_go.Client
}

func NewMessageService(
	validator *validator.Validate,
	conversationRepository *repository.ConversationRepository,
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
		conversationRepository:      conversationRepository,
		memberRepository:            memberRepository,
		messageRepository:           messageRepository,
		messageAttachmentRepository: messageAttachmentRepository,
	}
}

func (cs *MessageService) HandleUpgradeWs(ctx context.Context, claims *util.Claims, w http.ResponseWriter, r *http.Request) error {
	userId := claims.Sub

	user, err := cs.userRepository.TakeById(ctx, userId)
	if err != nil {
		fmt.Println("error 1:", err.Error())
		return err
	}

	ws, err := websocket.Upgrader.Upgrade(w, r, nil)
	if err != nil {

		fmt.Println("error 2:", err.Error())
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
		_, err := messageRepoTx.TakeById(ctx, messageID)
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

		conversationTx := ms.conversationRepository.WithTx(tx)
		memberTx := ms.memberRepository.WithTx(tx)
		messageTx := ms.messageRepository.WithTx(tx)

		var conversation *entity.Conversation
		var isNewConversation bool
		var err error

		if bc.Request.ConversationID == nil {

			switch bc.Request.Target.Type {
			case dto.TYPE_TARGET_GROUP:
				conversation, err = conversationTx.TakeGroupConversationByUserAndConversationId(ctx, bc.Sender.ID, bc.Request.Target.ID)
			case dto.TYPE_TARGET_USER:
				senderId := bc.Sender.ID
				receiverId := bc.Request.Target.ID
				conversation, err = conversationTx.TakePrivateConversationBySenderAndReceiverId(ctx, senderId, receiverId)
			default:
				err = errors.New("invalid target type")
			}

			isNotFound := errors.Is(err, gorm.ErrRecordNotFound)

			if err != nil && ((bc.Request.Target.Type == dto.TYPE_TARGET_GROUP) || (!isNotFound && bc.Request.Target.Type == dto.TYPE_TARGET_USER)) {
				return err
			}

			// create new conversation if target type is user
			if isNotFound {
				newConversation := &entity.Conversation{
					Bio:              "~",
					Name:             fmt.Sprintf("%s-%d", entity.CONV_TYPE_PRIVATE, time.Now().Unix()),
					ConversationType: entity.CONV_TYPE_PRIVATE,
					Image:            storage.DEFAULT_CONVERSATION_PROFILE_PICTURE,
				}

				err = conversationTx.Create(ctx, newConversation)
				if err != nil {
					return err
				}

				conversation = newConversation
				isNewConversation = true

				newMembers := []*entity.Member{
					{
						ConversationID: newConversation.ID,
						UserID:         bc.Sender.ID,
						Role:           entity.MEMBER_TYPE_MEMBER,
					},
					{
						ConversationID: newConversation.ID,
						UserID:         bc.Request.Target.ID,
						Role:           entity.MEMBER_TYPE_MEMBER,
					},
				}

				err = memberTx.Creates(ctx, newMembers)
				if err != nil {
					return err
				}
			}
		}

		newMessage := &entity.Message{
			UserID:         bc.Sender.ID,
			ConversationID: conversation.ID,
			Message:        bc.Request.Message,
		}

		err = messageTx.Create(ctx, newMessage)
		if err != nil {
			return err
		}

		bc.Message = newMessage
		bc.Message.User = bc.Sender

		bc.Request.ConversationID = &conversation.ID

		if !hub.IsExistConversation(conversation.ID) {
			memberIds, err := memberTx.GetUserIdsWithConversationId(ctx, conversation.ID)
			if err != nil {
				return err
			}

			hub.CreateConversation(conversation.ID, memberIds)
		}

		if isNewConversation {
			// send to sender
			err := hub.SendNewConversation(bc.Request.Target.ID, &dto.NewConversationResponse{
				ID:               bc.Sender.ID,
				Image:            bc.Sender.Image,
				Name:             bc.Sender.Name,
				Bio:              bc.Sender.Bio,
				ConversationType: entity.CONV_TYPE_PRIVATE,
				ConversationID:   &conversation.ID,
			})
			if err != nil {
				return err
			}

			receiver := hub.GetClient(bc.Request.Target.ID)

			if receiver == nil {
				receiver, err = ms.userRepository.TakeById(ctx, bc.Request.Target.ID)
			}
			if err != nil {
				return err
			}

			err = hub.SendNewConversation(bc.Sender.ID, &dto.NewConversationResponse{
				ID:               receiver.ID,
				Image:            receiver.Image,
				Name:             receiver.Name,
				Bio:              receiver.Bio,
				ConversationType: entity.CONV_TYPE_PRIVATE,
				ConversationID:   &conversation.ID,
			})

			if err != nil {
				return err
			}

		}
		return nil
	})

	return err
}

func (ms *MessageService) HandleGetMessages(ctx context.Context, conversationId int, claims *util.Claims) ([]*entity.Message, error) {
	_, err := ms.memberRepository.TakeByUserIdAndConversationId(ctx, claims.Sub, conversationId)
	if err != nil {
		return nil, err
	}

	return ms.messageRepository.GetMessages(ctx, conversationId)
}
