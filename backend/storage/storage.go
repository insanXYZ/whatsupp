package storage

import (
	"fmt"
	"mime/multipart"
	"os"
	"path/filepath"
	"time"

	storage_go "github.com/supabase-community/storage-go"
)

const (
	FILE_ATTACHMENT_BUCKET      = "file-attachment"
	USER_PROFILE_BUCKET         = "user-profile"
	CONVERSATION_PROFILE_BUCKET = "conversation-profile"
)

var (
	DEFAULT_PROFILE_PICTURE_URL          string
	DEFAULT_CONVERSATION_PROFILE_PICTURE string
)

func init() {

	DEFAULT_PROFILE_PICTURE_URL = fmt.Sprintf("%s/object/public/%s/default.png",
		os.Getenv("SUPABASE_STORAGE_RAW_URL"),
		USER_PROFILE_BUCKET,
	)

	DEFAULT_CONVERSATION_PROFILE_PICTURE = fmt.Sprintf("%s/object/public/%s/default.png",
		os.Getenv("SUPABASE_STORAGE_RAW_URL"),
		CONVERSATION_PROFILE_BUCKET,
	)

}

type Storage struct {
	Client *storage_go.Client
}

func NewStorage(client *storage_go.Client) *Storage {
	return &Storage{
		Client: client,
	}
}

func (s *Storage) uploadFile(fileheader *multipart.FileHeader, filename, bucket string) (string, error) {
	file, err := fileheader.Open()
	if err != nil {
		return "", nil
	}

	defer file.Close()

	contentType := fileheader.Header.Get("Content-Type")

	isUpsert := true

	fileOption := storage_go.FileOptions{
		ContentType: &contentType,
		Upsert:      &isUpsert,
	}

	_, err = s.Client.UploadFile(bucket, filename, file, fileOption)
	if err != nil {
		return "", err
	}

	signed := s.Client.GetPublicUrl(bucket, filename)

	return signed.SignedURL, nil
}

func (s *Storage) UploadFileConversationProfile(image *multipart.FileHeader, conversationId int) (string, error) {

	filename := fmt.Sprintf("conversation-%d%s", conversationId, filepath.Ext(image.Filename))

	return s.uploadFile(image, filename, CONVERSATION_PROFILE_BUCKET)
}

func (s *Storage) UploadFileUserProfile(image *multipart.FileHeader, userId int) (string, error) {

	filename := fmt.Sprintf("user-%d%s", userId, filepath.Ext(image.Filename))

	return s.uploadFile(image, filename, CONVERSATION_PROFILE_BUCKET)
}

func (s *Storage) UploadFileAttachment(file *multipart.FileHeader, messageId int) (string, error) {
	filename := fmt.Sprintf("attachment-%d-%d%s", messageId, time.Now().Unix(), filepath.Ext(file.Filename))

	return s.uploadFile(file, filename, FILE_ATTACHMENT_BUCKET)
}
