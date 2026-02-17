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
