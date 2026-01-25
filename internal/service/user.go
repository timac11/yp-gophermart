package service

import (
	"context"

	"github.com/timac11/yp-gophermart/internal/common/util"
	"github.com/timac11/yp-gophermart/internal/errors"
	"github.com/timac11/yp-gophermart/internal/model"
)

type UserService struct {
	repository UserRepository
	config     *ServiceConfig
}

type UserRepository interface {
	CreateUser(ctx context.Context, value *model.UserLoginDto) (*model.User, error)
	GetUserById(ctx context.Context, id string) (*model.User, error)
	GetUserByLogin(ctx context.Context, login string) (*model.User, error)
}

func (service *UserService) Register(ctx context.Context, value *model.UserLoginDto) (*model.User, error) {

	password, err := util.HashPassword(value.Password)

	if err != nil {
		return nil, err
	}

	userModel, err := service.repository.CreateUser(ctx, &model.UserLoginDto{Login: value.Login, Password: password})

	if err != nil {
		return nil, err
	}

	return userModel, nil

}

func (service *UserService) Login(ctx context.Context, value *model.UserLoginDto) (*model.User, error) {
	userModel, err := service.repository.GetUserByLogin(ctx, value.Login)

	if err != nil {
		return nil, err
	}

	passwordCorrect := util.CheckPasswordHash(value.Password, userModel.Password)

	if !passwordCorrect {
		return nil, errors.NewInvalidPasswordError(value.Password)
	}

	return userModel, nil
}

func NewUserService(repository UserRepository, config *ServiceConfig) *UserService {
	return &UserService{
		repository: repository,
		config:     config,
	}
}
