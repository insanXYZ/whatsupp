package storage

import (
	"fmt"
	"os"
)

const (
	FILE_ATTACHMENT_BUCKET = "file-attachment"
	USER_PROFILE_BUCKET    = "user-profile"
	GROUP_PROFILE_BUCKET   = "group-profile"
)

var (
	DEFAULT_PROFILE_PICTURE_URL       string
	DEFAULT_GROUP_PROFILE_PICTURE_URL string
)

func InitStorage() {

	DEFAULT_PROFILE_PICTURE_URL = fmt.Sprintf("%s/object/public/%s/profile.png",
		os.Getenv("SUPABASE_STORAGE_RAW_URL"),
		USER_PROFILE_BUCKET,
	)

	DEFAULT_GROUP_PROFILE_PICTURE_URL = fmt.Sprintf("%s/object/public/%s/profile.png",
		os.Getenv("SUPABASE_STORAGE_RAW_URL"),
		GROUP_PROFILE_BUCKET,
	)
}
