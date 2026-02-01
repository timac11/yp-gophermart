package service

import (
	"context"

	"github.com/timac11/yp-gophermart/internal/common/util"
	"github.com/timac11/yp-gophermart/internal/errors"
	"github.com/timac11/yp-gophermart/internal/model"
)

type BalanceService struct {
	repository BalanceRepository
	config     *ServiceConfig
}

type BalanceRepository interface {
	GetWithdrawals(ctx context.Context) ([]*model.WithdrawHistory, error)
	CreateWithdraw(ctx context.Context, withdraw *model.CreateWithdraw) (*model.WithdrawModel, error)
	GetBalance(ctx context.Context) (*model.BalanceInfo, error)
}

func (service *BalanceService) GetWithdrawals(ctx context.Context) ([]*model.WithdrawHistory, error) {
	return service.repository.GetWithdrawals(ctx)
}

func (service *BalanceService) CreateWithdraw(ctx context.Context, withdraw *model.CreateWithdraw) (*model.WithdrawModel, error) {
	if !util.CheckOrderNum(withdraw.Order) {
		return nil, errors.NewInvalidOrderNumErr(withdraw.Order)
	}

	return service.repository.CreateWithdraw(ctx, withdraw)
}

func (service *BalanceService) GetBalance(ctx context.Context) (*model.BalanceInfo, error) {
	return service.repository.GetBalance(ctx)
}

func NewBalanceService(repository BalanceRepository, config *ServiceConfig) *BalanceService {
	return &BalanceService{
		repository: repository,
		config:     config,
	}
}
