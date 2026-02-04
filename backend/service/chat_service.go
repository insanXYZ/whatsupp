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
	"whatsupp-backend/entity"
	"whatsupp-backend/repository"
	"whatsupp-backend/storage"
	"whatsupp-backend/websocket"

	"github.com/go-playground/validator/v10"
	"github.com/golang-jwt/jwt/v5"
	storage_go "github.com/supabase-community/storage-go"
	"gorm.io/gorm"
)

type ChatService struct {
	validator                   *validator.Validate
	groupRepository             *repository.GroupRepository
	memberRepository            *repository.MemberRepository
	messageRepository           *repository.MessageRepository
	messageAttachmentRepository *repository.MessageAttachmentRepository
	hub                         *websocket.Hub
	storage                     *storage_go.Client
}

func NewChatService(
	groupRepository *repository.GroupRepository,
	memberRepository *repository.MemberRepository,
	messageRepository *repository.MessageRepository,
	messageAttachmentRepository *repository.MessageAttachmentRepository,
	validator *validator.Validate,
	hub *websocket.Hub,
	storage *storage_go.Client,
) *ChatService {
	return &ChatService{
		storage:                     storage,
		hub:                         hub,
		validator:                   validator,
		groupRepository:             groupRepository,
		memberRepository:            memberRepository,
		messageRepository:           messageRepository,
		messageAttachmentRepository: messageAttachmentRepository,
	}
}

func (cs *ChatService) HandleUpgradeWs(ctx context.Context, claims jwt.MapClaims, w http.ResponseWriter, r *http.Request) error {
	ws, err := websocket.Upgrader.Upgrade(w, r, nil)
	if err != nil {
		return err
	}

	userId, err := claims.GetSubject()
	if err != nil {
		return errors.New("error claims subject")
	}

	atoiId, err := strconv.Atoi(userId)
	if err != nil {

		return errors.New("error atoi")
	}

	client := &websocket.Client{
		Id:   atoiId,
		Hub:  cs.hub,
		Conn: ws,
		Send: make(chan []byte, 256),
	}

	client.Hub.Register(client)

	go client.ReadPump()
	go client.WritePump()

	return nil
}

func (cs *ChatService) HandleUploadFileAttachments(ctx context.Context, messageID string, files []*multipart.FileHeader) error {

	return cs.messageAttachmentRepository.DB.Transaction(func(tx *gorm.DB) error {

		messageRepoTx := cs.messageRepository.WithTx(tx)
		err := messageRepoTx.TakeById(ctx, &entity.Message{}, messageID)
		if err != nil {
			return err
		}

		messageAttTx := cs.messageAttachmentRepository.WithTx(tx)

		clientStorage := cs.storage

		for _, file := range files {

			ext := filepath.Ext(file.Filename)

			fileName := fmt.Sprintf("whatsupp-%s-%s%s",
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
