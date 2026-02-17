package converter

import (
	"whatsupp-backend/dto"
	"whatsupp-backend/entity"
)

func MessageEntityToDto(message *entity.Message) *dto.Message {
	if message == nil {
		return nil
	}

	return &dto.Message{
		ID:           message.ID,
		Message:      message.Message,
		CreatedAt:    message.CreatedAt,
		Conversation: ConversationEntityToDto(message.Conversation),
		User:         UserEntityToDto(message.User),
	}
}

func MessageEntitytoItemGetMessagesResponseDto(message *entity.Message, userId int) *dto.ItemGetMessagesResponse {

	messageDto := MessageEntityToDto(message)

	if messageDto == nil {
		return nil
	}

	return &dto.ItemGetMessagesResponse{
		IsMe:    messageDto.User.ID == userId,
		Message: messageDto,
	}

}
