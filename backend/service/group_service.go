package service

import (
	"context"
	"fmt"
	"whatsupp-backend/dto"
	"whatsupp-backend/entity"
	"whatsupp-backend/repository"
)

type GroupService struct {
	userRepository  *repository.UserRepository
	groupRepository *repository.GroupRepository
}

func NewGroupRepository(userRepository *repository.UserRepository, groupRepository *repository.GroupRepository) *GroupService {
	return &GroupService{
		userRepository:  userRepository,
		groupRepository: groupRepository,
	}
}

func (u *GroupService) HandleLists(ctx context.Context, req *dto.ListGroupRequest) ([]entity.User, error) {
	nameFilter := fmt.Sprintf("%%%s%%", req.Name)
	var users []entity.User

	err := u.userRepository.FindByName(ctx, nameFilter, &users)
	return users, err
}
