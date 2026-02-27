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
		ID:             member.ID,
		Role:           string(member.Role),
		ConversationId: member.ConversationID,
		User:           UserEntityToDto(member.User),
		Conversation:   ConversationEntityToDto(member.Conversation),
	}
}

func MemberEntitiesToDto(members []*entity.Member) []*dto.Member {
	if members == nil {
		return nil
	}

	membersDto := make([]*dto.Member, len(members))

	for i, member := range members {
		membersDto[i] = MemberEntityToDto(member)
	}

	return membersDto
}
