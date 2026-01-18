package service

import (
	"context"
	"errors"

	"github.com/timac11/yp-gophermart/internal/common/util"
	"github.com/timac11/yp-gophermart/internal/model"
)

var (
	UserNotFound      = errors.New("User not found")
	UserAlreadyExist  = errors.New("User already exist")
	UserWrongPassword = errors.New("Password is not valid")
)

type ServiceConfig struct {
	Attempts         uint
	AttemptsInterval uint
}

type UserService struct {
	repository UserRepository
	config     ServiceConfig
}

type UserRepository interface {
	SaveUser(ctx context.Context, value model.UserDto) (*model.User, error)
	GetUserById(ctx context.Context, id string) (*model.User, error)
	GetUserByLogin(ctx context.Context, login string) (*model.User, error)
}

func (service *UserService) Register(ctx context.Context, value model.UserDto) (*model.User, error) {

	password, err := util.HashPassword(value.Password)

	if err != nil {
		return nil, err
	}

	value.Password = password

	userModel, err := service.repository.SaveUser(ctx, value)

	if err != nil {
		// TODO check err type and return definite type of error
		return nil, err
	}

	return userModel, nil

}

func (service *UserService) Login(ctx context.Context, value model.UserDto) (*model.User, error) {
	userModel, err := service.repository.GetUserByLogin(ctx, value.Login)

	if err != nil {
		return nil, err
	}

	passwordCorrect := util.CheckPasswordHash(value.Password, userModel.Password)

	if !passwordCorrect {
		return nil, UserWrongPassword
	}

	return userModel, nil
}

func NewUserService(repository UserRepository, config ServiceConfig) *UserService {
	return &UserService{
		repository: repository,
		config:     config,
	}
}
