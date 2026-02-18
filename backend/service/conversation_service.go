package service

import (
	"context"
	"whatsupp-backend/dto"
	"whatsupp-backend/repository"
	"whatsupp-backend/util"

	"github.com/go-playground/validator/v10"
)

type ConversationService struct {
	validator              *validator.Validate
	conversationRepository *repository.ConversationRepository
	memberRepsitory        *repository.MemberRepository
}

func NewConversationService(validator *validator.Validate, memberRepository *repository.MemberRepository, conversationRepository *repository.ConversationRepository) *ConversationService {
	return &ConversationService{
		validator:              validator,
		memberRepsitory:        memberRepository,
		conversationRepository: conversationRepository,
	}
}

func (cs *ConversationService) HandleFindConversations(ctx context.Context, claims *util.Claims, req *dto.SearchConversationRequest) ([]dto.SearchConversationResponse, error) {
	return cs.conversationRepository.SearchConversationWithNameAndUserId(ctx, claims.Sub, req.Name)
}

func (cs *ConversationService) HandleLoadRecentConversations(ctx context.Context, claims *util.Claims) ([]dto.LoadRecentConversation, error) {
	return cs.conversationRepository.FindConversationsByUserId(ctx, claims.Sub)
}
