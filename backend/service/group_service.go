package service

import (
	"context"
	"whatsupp-backend/dto"
	"whatsupp-backend/repository"

	"github.com/go-playground/validator/v10"
	"github.com/golang-jwt/jwt/v5"
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

func (g *GroupService) HandleLists(ctx context.Context, claims jwt.MapClaims, req *dto.SearchGroupRequest) ([]dto.SearchGroupResponse, error) {
	currUserId := claims["sub"].(int)
	return g.groupRepository.SearchGroupAndUserWithName(ctx, currUserId, req.Name)
}
