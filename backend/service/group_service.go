package service

import (
	"context"
	"whatsupp-backend/dto"
	"whatsupp-backend/repository"
	"whatsupp-backend/util"

	"github.com/go-playground/validator/v10"
)

type GroupService struct {
	validator       *validator.Validate
	userRepository  *repository.UserRepository
	groupRepository *repository.GroupRepository
}

func NewGroupService(validator *validator.Validate, userRepository *repository.UserRepository, groupRepository *repository.GroupRepository) *GroupService {
	return &GroupService{
		validator:       validator,
		userRepository:  userRepository,
		groupRepository: groupRepository,
	}
}

func (g *GroupService) HandleLists(ctx context.Context, claims *util.Claims, req *dto.SearchGroupRequest) ([]dto.SearchGroupResponse, error) {
	return g.groupRepository.SearchGroupAndUserWithName(ctx, claims.Sub, req.Name)
}
