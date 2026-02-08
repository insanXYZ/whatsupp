package service

import (
	"context"
	"errors"
	"whatsupp-backend/dto"
	"whatsupp-backend/entity"
	"whatsupp-backend/repository"
	"whatsupp-backend/storage"
	"whatsupp-backend/util"

	"github.com/go-playground/validator/v10"
	"github.com/golang-jwt/jwt/v5"
	"gorm.io/gorm"
)

type UserService struct {
	validator      *validator.Validate
	userRepository *repository.UserRepository
}

func NewUserService(validator *validator.Validate, userRepository *repository.UserRepository) *UserService {
	return &UserService{
		userRepository: userRepository,
		validator:      validator,
	}
}

func (u *UserService) isUserExist(ctx context.Context, email string) bool {
	err := u.userRepository.TakeByEmail(ctx, email, &entity.User{})
	return !errors.Is(err, gorm.ErrRecordNotFound)
}

func (u *UserService) HandleLogin(ctx context.Context, req *dto.LoginRequest) (*entity.User, error) {
	err := u.validator.Struct(req)
	if err != nil {
		return nil, err
	}

	user := new(entity.User)

	err = u.userRepository.TakeByEmail(ctx, req.Email, user)
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
		Image:    storage.DEFAULT_PROFILE_PICTURE_URL,
	}

	return u.userRepository.Create(ctx, newUser)

}

func (u *UserService) HandleUpdateUser(ctx context.Context, req *dto.UpdateUserRequest, claims jwt.MapClaims) error {
	err := u.validator.Struct(req)
	if err != nil {
		return err
	}

	if req.Email != claims["email"].(string) && u.isUserExist(ctx, req.Email) {
		return errors.New("email has been used")
	}

	user := new(entity.User)
	err = u.userRepository.TakeByEmail(ctx, claims["email"].(string), user)
	if err != nil {
		return err
	}

	user.Email = req.Email
	user.Name = req.Name

	if req.Password != "" {
		pw, err := util.GenerateBcrypt(req.Password)
		if err != nil {
			return err
		}
		user.Password = pw
	}

	return u.userRepository.Update(ctx, user)
}

func (u *UserService) HandleMe(ctx context.Context, claims jwt.MapClaims) (*entity.User, error) {
	user := new(entity.User)
	err := u.userRepository.TakeById(ctx, user, claims["id"].(string))
	return user, err
}
