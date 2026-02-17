package converter

import (
	"whatsupp-backend/dto"
	"whatsupp-backend/entity"
)

func MemberEntityToDto(member *entity.Member) *dto.Member {
	if member == nil {
		return nil
	}

	return &dto.Member{
		ID:           member.ID,
		Role:         string(member.Role),
		User:         UserEntityToDto(member.User),
		Conversation: ConversationEntityToDto(member.Conversation),
	}
}
