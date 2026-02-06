package config

import (
	"os"
	"strings"
	"whatsupp-backend/storage"

	storage_go "github.com/supabase-community/storage-go"
)

func createBucketIfNotExists(
	client *storage_go.Client,
	name string,
	opts storage_go.BucketOptions,
) error {

	_, err := client.CreateBucket(name, opts)
	if err == nil {
		return nil
	}

	msg := strings.ToLower(err.Error())

	if strings.Contains(msg, "already exists") ||
		strings.Contains(msg, "duplicate") ||
		strings.Contains(msg, "409") {
		return nil
	}

	return err
}

func NewSupabaseStorageClient() (*storage_go.Client, error) {
	client := storage_go.NewClient(
		os.Getenv("SUPABASE_STORAGE_RAW_URL"),
		os.Getenv("SUPABASE_STORAGE_TOKEN"),
		nil,
	)

	if err := createBucketIfNotExists(client, storage.FILE_ATTACHMENT_BUCKET, storage_go.BucketOptions{
		Public: true,
	}); err != nil {
		return nil, err
	}

	if err := createBucketIfNotExists(client, storage.USER_PROFILE_BUCKET, storage_go.BucketOptions{
		Public: true,
	}); err != nil {
		return nil, err
	}

	if err := createBucketIfNotExists(client, storage.GROUP_PROFILE_BUCKET, storage_go.BucketOptions{
		Public: true,
	}); err != nil {
		return nil, err
	}

	storage.InitStorage()

	return client, nil
}
