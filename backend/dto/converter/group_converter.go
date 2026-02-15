package converter

import (
	"whatsupp-backend/dto"
	"whatsupp-backend/entity"
)

func GroupEntityToDto(group *entity.Group) *dto.Group {
	if group == nil {
		return nil
	}

	return &dto.Group{
		ID:   group.ID,
		Name: group.Name,
		Bio:  group.Bio,
		Type: group.GroupType,
	}
}
