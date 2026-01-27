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
		Email: user.Email,
	}
}
