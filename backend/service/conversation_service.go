package service

import (
	"context"
	"errors"
	"whatsupp-backend/dto"
	"whatsupp-backend/dto/converter"
	"whatsupp-backend/dto/message"
	"whatsupp-backend/entity"
	"whatsupp-backend/repository"
	"whatsupp-backend/storage"
	"whatsupp-backend/util"
	"whatsupp-backend/websocket"

	"github.com/go-playground/validator/v10"
	"github.com/insanXYZ/sage"
	"gorm.io/gorm"
)

type ConversationService struct {
	validator              *validator.Validate
	conversationRepository *repository.ConversationRepository
	memberRepository       *repository.MemberRepository
	storage                *storage.Storage
	hub                    *websocket.Hub
}

func NewConversationService(validator *validator.Validate, memberRepository *repository.MemberRepository, conversationRepository *repository.ConversationRepository, storage *storage.Storage, hub *websocket.Hub) *ConversationService {

	conversationService := &ConversationService{
		storage:                storage,
		validator:              validator,
		memberRepository:       memberRepository,
		conversationRepository: conversationRepository,
		hub:                    hub,
	}

	hub.SyncConversation = conversationService.syncHubConversation

	return conversationService
}

func (cs *ConversationService) syncHubConversation(conversationId int) error {
	ctx := context.Background()
	hub := cs.hub

	if hub.IsExistConversation(conversationId) {
		return nil
	}

	memberIds, err := cs.memberRepository.GetUserIdsWithConversationId(ctx, conversationId)
	if err != nil {
		return err
	}

	hub.CreateConversation(conversationId, memberIds)
	return nil
}

func (cs *ConversationService) HandleListMembersConversation(ctx context.Context, req *dto.ListMembersConversationRequest, claims *util.Claims) ([]*entity.Member, error) {
	conversation, err := cs.conversationRepository.TakeConversationByConversationAndUserId(ctx, req.ConversationID, claims.Sub)
	if err != nil {
		return nil, err
	}

	return cs.memberRepository.FindByConversationId(ctx, conversation.ID)

}

func (cs *ConversationService) HandleFindConversations(ctx context.Context, claims *util.Claims, req *dto.SearchConversationRequest) ([]dto.SearchConversationResponse, error) {
	return cs.conversationRepository.SearchConversationWithNameAndUserId(ctx, claims.Sub, req.Name)
}

func (cs *ConversationService) HandleLoadRecentConversations(ctx context.Context, claims *util.Claims) ([]*entity.Conversation, error) {
	return cs.conversationRepository.FindConversationsByUserId(ctx, claims.Sub)
}

func (cs *ConversationService) HandleCreateGroupConversation(ctx context.Context, req *dto.CreateGroupConversationRequest, claims *util.Claims) error {
	err := cs.validator.Struct(req)
	if err != nil {
		return err
	}

	if req.Image != nil {
		err = sage.Validate(req.Image)
		if err != nil {
			return err
		}

	}

	err = cs.conversationRepository.Transaction(ctx, func(tx *gorm.DB) error {

		conversationTx := cs.conversationRepository.WithTx(tx)
		memberTx := cs.memberRepository.WithTx(tx)

		imageConversation := storage.DEFAULT_CONVERSATION_PROFILE_PICTURE

		newConversation := &entity.Conversation{
			Name:             req.Name,
			Bio:              req.Bio,
			ConversationType: entity.CONV_TYPE_GROUP,
			Image:            imageConversation,
			Members: []*entity.Member{
				{
					UserID: claims.Sub,
					Role:   entity.MEMBER_TYPE_ADMIN,
				},
			},
		}

		err = conversationTx.Create(ctx, newConversation)
		if err != nil {
			return err
		}

		if req.Image != nil {
			imageUrl, err := cs.storage.UploadFileConversationProfile(req.Image, newConversation.ID)

			if err != nil {
				return err
			}

			newConversation.Image = imageUrl

			err = conversationTx.Update(ctx, newConversation)
			if err != nil {
				return err
			}
		}

		conversationSummary := &dto.ConversationSummary{
			ID:               newConversation.ID,
			Name:             newConversation.Name,
			Image:            newConversation.Image,
			Bio:              newConversation.Bio,
			ConversationType: newConversation.ConversationType,
			ConversationID:   &newConversation.ID,
			HaveJoined:       true,
		}

		members, err := memberTx.FindByConversationId(ctx, newConversation.ID)
		if err != nil {
			return err
		}

		cs.hub.CreateConversation(newConversation.ID, []int{claims.Sub})

		newConversationResponse := &dto.NewConversationResponse{
			ConversationSummary: conversationSummary,
			Members:             converter.MemberEntitiesToDto(members),
		}

		err = cs.hub.SendNewConversation(claims.Sub, newConversationResponse)
		if err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return err
	}

	return nil

}

func (cs *ConversationService) HandleJoinGroupConversation(ctx context.Context, req *dto.JoinGroupConversationRequest, claims *util.Claims) (bool, error) {

	var newMember *entity.Member
	var conversation *entity.Conversation
	isJoin := true

	err := cs.conversationRepository.Transaction(ctx, func(tx *gorm.DB) error {

		conversationTx := cs.conversationRepository.WithTx(tx)
		memberTx := cs.memberRepository.WithTx(tx)

		conversationWithMember, err := conversationTx.TakeGroupConversationLeftJoinMemberByUserAndConversationId(ctx, claims.Sub, req.ConversationID)
		if err != nil {
			return err
		}

		// members field empty
		// its mean userId isn't this member conversation
		isMemberConversation := len(conversationWithMember.Members) != 0

		if isMemberConversation {
			isJoin = false
			err = memberTx.DeleteById(ctx, conversationWithMember.Members[0].ID)

		} else {
			isJoin = true

			newMember = &entity.Member{
				UserID:         claims.Sub,
				Role:           entity.MEMBER_TYPE_MEMBER,
				ConversationID: conversationWithMember.ID,
			}

			err = memberTx.Create(ctx, newMember)
		}

		if err != nil {
			return err
		}

		conversation = conversationWithMember

		return nil
	})

	if err != nil {
		return isJoin, err
	}

	err = cs.hub.SyncConversation(conversation.ID)
	if err != nil {
		return isJoin, err
	}

	if isJoin {

		conversationSummary := &dto.ConversationSummary{
			ID:               conversation.ID,
			Name:             conversation.Name,
			Image:            conversation.Image,
			Bio:              conversation.Bio,
			ConversationID:   &conversation.ID,
			ConversationType: conversation.ConversationType,
			HaveJoined:       true,
		}

		members, err := cs.memberRepository.FindByConversationId(ctx, conversation.ID)
		if err != nil {
			return isJoin, err
		}

		newConversationResponse := &dto.NewConversationResponse{
			ConversationSummary: conversationSummary,
			Members:             converter.MemberEntitiesToDto(members),
		}

		err = cs.hub.SendNewConversation(claims.Sub, newConversationResponse)
		if err != nil {
			return isJoin, err

		}

		memberJoinConversationResponse := &dto.MemberJoinConversationResponse{
			ConversationId: conversation.ID,
			Member:         converter.MemberEntityToDto(newMember),
		}

		cs.hub.SendMemberJoinConversation(conversation.ID, memberJoinConversationResponse)
	} else {
		leaveConversationResponse := &dto.LeaveConversationResponse{
			ConversationID: conversation.ID,
		}
		err := cs.hub.SendLeaveConversation(claims.Sub, leaveConversationResponse)
		if err != nil {
			return isJoin, err
		}

		memberLeaveConversationResponse := &dto.MemberLeaveConversationResponse{
			ConversationId: conversation.ID,
			MemberId:       conversation.Members[0].ID,
		}

		err = cs.hub.SendMemberLeaveConversation(memberLeaveConversationResponse)
		if err != nil {
			return isJoin, err
		}
	}

	return isJoin, err

}

func (cs *ConversationService) HandleUpdateGroupConversation(ctx context.Context, req *dto.UpdateGroupConversationRequest, claims *util.Claims) error {
	err := cs.validator.Struct(req)
	if err != nil {
		return err
	}

	if req.Image != nil {
		err = sage.Validate(req.Image)
		if err != nil {
			return err
		}

	}

	err = cs.conversationRepository.Transaction(ctx, func(tx *gorm.DB) error {

		memberTx := cs.memberRepository.WithTx(tx)
		conversationTx := cs.conversationRepository.WithTx(tx)

		conversation, err := conversationTx.TakeById(ctx, req.ConversationId)
		if err != nil {
			return err
		}

		isAdmin, err := memberTx.IsAdminConversationByConversationAndUserId(ctx, conversation.ID, claims.Sub)
		if err != nil {
			return err
		}

		if !isAdmin {
			return errors.New(message.ERR_UPDATE_CONVERSATION_ACCESS_DENIED)
		}

		conversation.Name = req.Name
		conversation.Bio = req.Bio

		err = conversationTx.Update(ctx, conversation)
		if err != nil {
			return err
		}

		if req.Image != nil {
			imageUrl, err := cs.storage.UploadFileConversationProfile(req.Image, conversation.ID)
			if err != nil {
				return err
			}

			conversation.Image = imageUrl

			err = conversationTx.Update(ctx, conversation)
			if err != nil {
				return err
			}
		}

		return nil

	})

	return err

}
