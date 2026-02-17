package dto

import "whatsupp-backend/entity"

type BroadcastMessageWs struct {
	Request *SendMessageRequestWs
	Sender  *entity.User
	Message *entity.Message
}
