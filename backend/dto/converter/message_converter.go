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
		ID:        message.ID,
		Message:   message.Message,
		Member:    MemberEntityToDto(message.Member),
		CreatedAt: message.CreatedAt,
	}
}

func MessageEntitytoGetMessageResponseDto(message *entity.Message, userId int) *dto.GetMessagesResponse {

	messageDto := MessageEntityToDto(message)

	if messageDto == nil {
		return nil
	}

	return &dto.GetMessagesResponse{
		IsMe:    messageDto.Member.User.ID == userId,
		Message: *messageDto,
	}

}
