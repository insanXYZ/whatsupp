package service

import (
	"context"
	"errors"
	"fmt"
	"mime/multipart"
	"net/http"
	"path/filepath"
	"strconv"
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

type ChatService struct {
	validator                   *validator.Validate
	groupRepository             *repository.GroupRepository
	memberRepository            *repository.MemberRepository
	messageRepository           *repository.MessageRepository
	messageAttachmentRepository *repository.MessageAttachmentRepository
	userRepository              *repository.UserRepository
	hub                         *websocket.Hub
	storage                     *storage_go.Client
}

func NewChatService(
	validator *validator.Validate,
	groupRepository *repository.GroupRepository,
	memberRepository *repository.MemberRepository,
	messageRepository *repository.MessageRepository,
	messageAttachmentRepository *repository.MessageAttachmentRepository,
	userRepository *repository.UserRepository,
	hub *websocket.Hub,
	storage *storage_go.Client,
) *ChatService {
	return &ChatService{
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

func (cs *ChatService) HandleUpgradeWs(ctx context.Context, claims *util.Claims, w http.ResponseWriter, r *http.Request) error {
	ws, err := websocket.Upgrader.Upgrade(w, r, nil)
	if err != nil {
		return err
	}

	userId := claims.Sub

	client := &websocket.Client{
		Id:                    userId,
		Hub:                   cs.hub,
		Conn:                  ws,
		Send:                  make(chan *dto.BroadcastMessageWS, 250),
		HandlerSendMessage:    cs.handleSaveMessage,
		HandlerCreateNewGroup: cs.handleCreateNewGroup,
	}

	client.Hub.Register(client)

	go client.ReadPump()
	go client.WritePump()

	return nil
}

func (cs *ChatService) HandleUploadFileAttachments(ctx context.Context, messageID string, files []*multipart.FileHeader) error {

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

// return id of row insert message
func (cs *ChatService) handleSaveMessage(msg *dto.BroadcastMessageWS) (int, error) {
	ctx := context.Background()

	member := new(entity.Member)

	err := cs.memberRepository.TakeByUserIdAndGroupId(ctx, msg.ClientID, *msg.GroupID, member)
	if err != nil {
		return 0, err
	}

	newMessage := &entity.Message{
		MemberID: member.ID,
		Message:  msg.Message,
	}

	err = cs.messageRepository.Create(ctx, newMessage)

	return newMessage.ID, err

}

func (cs *ChatService) handleCreateNewGroup(senderId, receiverId int) (int, error) {
	ctx := context.Background()

	group := new(entity.Group)
	err := cs.groupRepository.TakePrivateGroupBySenderAndReceiverId(ctx, senderId, receiverId, group)

	if err != nil {

		if errors.Is(err, gorm.ErrRecordNotFound) {

			err := cs.groupRepository.Transaction(ctx, func(tx *gorm.DB) error {
				sender, receiver := new(entity.User), new(entity.User)

				userTx := cs.userRepository.WithTx(tx)
				groupTx := cs.groupRepository.WithTx(tx)
				memberTx := cs.memberRepository.WithTx(tx)

				err := userTx.TakeById(ctx, sender, senderId)
				if err != nil {
					return err
				}

				err = userTx.TakeById(ctx, receiver, receiverId)
				if err != nil {
					return err
				}

				newGroup := &entity.Group{
					Name:      fmt.Sprintf("%s-%s", "PRIVATE", strconv.Itoa(int(time.Now().Unix()))),
					Bio:       "~",
					GroupType: entity.GROUP,
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

				err = memberTx.Creates(ctx, members)
				if err != nil {
					return err
				}

				group = newGroup

				return nil

			})

			if err != nil {
				return 0, err
			}

		} else {
			return 0, err
		}

	}

	return group.ID, err
}
