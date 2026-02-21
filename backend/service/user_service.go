package service

import (
	"context"
	"errors"
	"whatsupp-backend/dto"
	"whatsupp-backend/entity"
	"whatsupp-backend/repository"
	"whatsupp-backend/storage"
	"whatsupp-backend/util"
	"whatsupp-backend/websocket"

	"github.com/go-playground/validator/v10"
	"gorm.io/gorm"
)

type UserService struct {
	hub            *websocket.Hub
	validator      *validator.Validate
	userRepository *repository.UserRepository
}

func NewUserService(validator *validator.Validate, userRepository *repository.UserRepository, hub *websocket.Hub) *UserService {
	return &UserService{
		hub:            hub,
		userRepository: userRepository,
		validator:      validator,
	}
}

func (u *UserService) isUserExist(ctx context.Context, email string) bool {
	_, err := u.userRepository.TakeByEmail(ctx, email)
	return !errors.Is(err, gorm.ErrRecordNotFound)
}

func (u *UserService) HandleLogin(ctx context.Context, req *dto.LoginRequest) (*entity.User, error) {
	err := u.validator.Struct(req)
	if err != nil {
		return nil, err
	}

	user, err := u.userRepository.TakeByEmail(ctx, req.Email)
	if err != nil {
		return nil, err
	}

	if util.ComparePassword(req.Password, user.Password) != nil {
		return nil, errors.New("email or password wrong")
	}

	return user, nil
}

func (u *UserService) HandleRegister(ctx context.Context, req *dto.RegisterRequest) error {
	err := u.validator.Struct(req)
	if err != nil {
		return err
	}

	if u.isUserExist(ctx, req.Email) {
		return errors.New("user has been used")
	}

	password, err := util.GenerateBcrypt(req.Password)
	if err != nil {
		return err
	}

	newUser := &entity.User{
		Name:     req.Name,
		Email:    req.Email,
		Password: password,
		Bio:      "~",
		Image:    storage.DEFAULT_PROFILE_PICTURE_URL,
	}

	return u.userRepository.Create(ctx, newUser)

}

func (u *UserService) HandleUpdateUser(ctx context.Context, req *dto.UpdateUserRequest, claims *util.Claims) (*entity.User, error) {
	err := u.validator.Struct(req)
	if err != nil {
		return nil, err
	}

	if req.Email != claims.Email && u.isUserExist(ctx, req.Email) {
		return nil, errors.New("email has been used")
	}

	user, err := u.userRepository.TakeByEmail(ctx, claims.Email)
	if err != nil {
		return nil, err
	}

	user.Email = req.Email
	user.Name = req.Name

	if req.Password != "" {
		pw, err := util.GenerateBcrypt(req.Password)
		if err != nil {
			return nil, err
		}
		user.Password = pw
	}

	err = u.userRepository.Update(ctx, user)
	if err != nil {
		return nil, err
	}

	u.hub.UpdateClient(user.ID, user)

	return user, err
}

func (u *UserService) HandleMe(ctx context.Context, claims *util.Claims) (*entity.User, error) {
	return u.userRepository.TakeById(ctx, claims.Sub)
}
