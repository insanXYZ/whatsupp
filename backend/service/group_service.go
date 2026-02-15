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
	groupRepository *repository.GroupRepository
	memberRepsitory *repository.MemberRepository
}

func NewGroupService(validator *validator.Validate, memberRepository *repository.MemberRepository, groupRepository *repository.GroupRepository) *GroupService {
	return &GroupService{
		validator:       validator,
		memberRepsitory: memberRepository,
		groupRepository: groupRepository,
	}
}

func (gs *GroupService) HandleFindGroups(ctx context.Context, claims *util.Claims, req *dto.SearchGroupRequest) ([]dto.SearchGroupResponse, error) {
	return gs.groupRepository.SearchGroupAndUserWithName(ctx, claims.Sub, req.Name)
}

func (gs *GroupService) HandleLoadRecentGroups(ctx context.Context, claims *util.Claims) ([]dto.LoadRecentGroup, error) {
	return gs.groupRepository.FindAllGroupWithMemberUserId(ctx, claims.Sub)
}
