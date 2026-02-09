package converter

import (
	"whatsupp-backend/dto"
	"whatsupp-backend/entity"
)

func UserEntityToDto(user *entity.User) *dto.User {
	if user == nil {
		return nil
	}

	return &dto.User{
		ID:    user.ID,
		Name:  user.Name,
		Image: user.Image,
		Bio:   user.Bio,
		Email: user.Email,
	}
}
