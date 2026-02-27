package converter

import (
	"whatsupp-backend/dto"
	"whatsupp-backend/entity"
)

func ConversationEntityToDto(ent *entity.Conversation) *dto.Conversation {
	if ent == nil {
		return nil
	}

	return &dto.Conversation{
		ID:        ent.ID,
		Name:      ent.Name,
		Bio:       ent.Bio,
		Image:     ent.Image,
		CreatedAt: ent.CreatedAt,
	}
}

func ConversationEntityToLoadRecentConversationDto(conversation *entity.Conversation, userId int) *dto.LoadRecentConversation {
	if conversation == nil {
		return nil
	}

	isPrivate := conversation.ConversationType == entity.CONV_TYPE_PRIVATE

	conversationSummary := &dto.ConversationSummary{
		ID:               conversation.ID,
		Name:             conversation.Name,
		Bio:              conversation.Bio,
		Image:            conversation.Image,
		ConversationType: conversation.ConversationType,
		HaveJoined:       true,
		ConversationID:   &conversation.ID,
	}

	if isPrivate {
		if len(conversation.Members) != 2 {
			return nil
		}

		receiver := conversation.Members[0]
		if receiver.User.ID == userId {
			receiver = conversation.Members[1]
		}

		userReceiver := receiver.User

		conversationSummary.Name = userReceiver.Name
		conversationSummary.Bio = userReceiver.Bio
		conversationSummary.Image = userReceiver.Image
	}

	recentConversation := &dto.LoadRecentConversation{
		ConversationSummary: conversationSummary,
		Members:             MemberEntitiesToDto(conversation.Members),
	}

	return recentConversation
}

func ConversationEntitiesToLoadRecentConversationsDto(conversations []*entity.Conversation, userId int) []*dto.LoadRecentConversation {
	if conversations == nil {
		return nil
	}

	recentConversations := make([]*dto.LoadRecentConversation, len(conversations))

	for i, conversation := range conversations {
		recentConversations[i] = ConversationEntityToLoadRecentConversationDto(conversation, userId)
	}

	return recentConversations

}
