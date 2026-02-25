package service

import (
	"context"
	"errors"
	"whatsupp-backend/dto"
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
	memberRepsitory        *repository.MemberRepository
	storage                *storage.Storage
	hub                    *websocket.Hub
}

func NewConversationService(validator *validator.Validate, memberRepository *repository.MemberRepository, conversationRepository *repository.ConversationRepository, storage *storage.Storage, hub *websocket.Hub) *ConversationService {

	return &ConversationService{
		storage:                storage,
		validator:              validator,
		memberRepsitory:        memberRepository,
		conversationRepository: conversationRepository,
		hub:                    hub,
	}
}

func (cs *ConversationService) HandleFindConversations(ctx context.Context, claims *util.Claims, req *dto.SearchConversationRequest) ([]dto.SearchConversationResponse, error) {
	return cs.conversationRepository.SearchConversationWithNameAndUserId(ctx, claims.Sub, req.Name)
}

func (cs *ConversationService) HandleLoadRecentConversations(ctx context.Context, claims *util.Claims) ([]*entity.Conversation, error) {
	return cs.conversationRepository.FindConversationsByUserId(ctx, claims.Sub)
}

func (cs *ConversationService) HandleCreateGroupConversation(ctx context.Context, req *dto.CreateGroupConversationRequest, claims *util.Claims) (*entity.Conversation, error) {
	err := cs.validator.Struct(req)
	if err != nil {
		return nil, err
	}

	if req.Image != nil {
		err = sage.Validate(req.Image)
		if err != nil {
			return nil, err
		}

	}

	var newConversation *entity.Conversation

	err = cs.conversationRepository.Transaction(ctx, func(tx *gorm.DB) error {

		conversationTx := cs.conversationRepository.WithTx(tx)

		imageConversation := storage.DEFAULT_CONVERSATION_PROFILE_PICTURE

		newConversation = &entity.Conversation{
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

		newConversationResponse := &dto.NewConversationResponse{
			ID:               newConversation.ID,
			Name:             newConversation.Name,
			Image:            newConversation.Image,
			Bio:              newConversation.Bio,
			ConversationType: newConversation.ConversationType,
			ConversationID:   &newConversation.ID,
			HaveJoined:       true,
		}

		cs.hub.CreateConversation(newConversation.ID, []int{claims.Sub})

		err = cs.hub.SendNewConversation(claims.Sub, newConversationResponse)
		if err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return newConversation, err

}

func (cs *ConversationService) HandleJoinGroupConversation(ctx context.Context, req *dto.JoinGroupConversationRequest, claims *util.Claims) (bool, error) {

	isJoin := true

	err := cs.conversationRepository.Transaction(ctx, func(tx *gorm.DB) error {

		conversationTx := cs.conversationRepository.WithTx(tx)
		memberTx := cs.memberRepsitory.WithTx(tx)

		conversation, err := conversationTx.TakeById(ctx, req.ConversationID)
		if err != nil {
			return err
		}

		conversationWithMember, err := conversationTx.TakeGroupConversationByUserAndConversationId(ctx, claims.Sub, conversation.ID)
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return err
		}

		if err != nil {
			isJoin = true

			newMember := &entity.Member{
				UserID:         claims.Sub,
				Role:           entity.MEMBER_TYPE_MEMBER,
				ConversationID: conversation.ID,
			}

			err = memberTx.Create(ctx, newMember)
		} else {
			isJoin = false
			err = memberTx.DeleteById(ctx, conversationWithMember.Members[0].ID)
		}

		if err != nil {
			return err
		}

		if isJoin {
			newConversationResponse := &dto.NewConversationResponse{
				ID:               conversation.ID,
				Name:             conversation.Name,
				Image:            conversation.Image,
				Bio:              conversation.Bio,
				ConversationID:   &conversation.ID,
				ConversationType: conversation.ConversationType,
				HaveJoined:       true,
			}
			err = cs.hub.SendNewConversation(claims.Sub, newConversationResponse)
		} else {
			leaveConversationResponse := &dto.LeaveConversationResponse{
				ConversationID: conversation.ID,
			}

			err = cs.hub.SendLeaveConversation(claims.Sub, leaveConversationResponse)
		}

		if !cs.hub.IsExistConversation(conversation.ID) {
			memberIds, err := memberTx.GetUserIdsWithConversationId(ctx, conversation.ID)
			if err != nil {
				return err
			}

			cs.hub.CreateConversation(conversation.ID, memberIds)
		} else {
			cs.hub.DeleteClientConversation(conversation.ID, claims.Sub)
		}

		return err
	})

	return isJoin, err

}

func (cs *ConversationService) HandleListMembersConversation(ctx context.Context, req *dto.ListMembersConversationRequest, claims *util.Claims) ([]*entity.Member, error) {
	conversation, err := cs.conversationRepository.TakeConversationByConversationAndUserId(ctx, req.ConversationID, claims.Sub)
	if err != nil {
		return nil, err
	}

	return cs.memberRepsitory.FindMembersWithConversationId(ctx, conversation.ID)

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

		memberTx := cs.memberRepsitory.WithTx(tx)
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
