package main

import (
	"encoding/json"
	"fmt"
	"testing"
	"whatsupp-backend/config"
	"whatsupp-backend/entity"

	"github.com/joho/godotenv"
	"gorm.io/gorm"
)

var db *gorm.DB

func init() {
	err := godotenv.Load()
	if err != nil {
		panic(err)
	}

	db, err = config.NewGorm()
	if err != nil {
		panic(err)
	}

}

func TestMain(t *testing.T) {

	conversation := new(entity.Conversation)

	err := db.Debug().Joins("JOIN members ON members.conversation_id = conversations.id AND members.user_id = ?", 2).Preload("Members").Take(conversation, "conversations.id = ?", 57).Error
	if err != nil {
		panic(err)
	}

	b, err := json.MarshalIndent(conversation, "", " ")
	if err != nil {
		panic(err)
	}

	fmt.Println(string(b))
}
