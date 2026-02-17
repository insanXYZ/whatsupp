package storage

import (
	"fmt"
	"os"
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

func InitStorage() {

	DEFAULT_PROFILE_PICTURE_URL = fmt.Sprintf("%s/object/public/%s/default.png",
		os.Getenv("SUPABASE_STORAGE_RAW_URL"),
		USER_PROFILE_BUCKET,
	)

	DEFAULT_CONVERSATION_PROFILE_PICTURE = fmt.Sprintf("%s/object/public/%s/default.png",
		os.Getenv("SUPABASE_STORAGE_RAW_URL"),
		CONVERSATION_PROFILE_BUCKET,
	)
}
