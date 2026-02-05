package service

import (
	"context"

	"github.com/timac11/yp-gophermart/internal/common/util"
	"github.com/timac11/yp-gophermart/internal/errors"
	"github.com/timac11/yp-gophermart/internal/model"
)

type OrderService struct {
	repository OrderRepository
}

type OrderRepository interface {
	CreateOrder(ctx context.Context, value string) (*model.OrderModel, error)
	GetOrders(ctx context.Context) ([]*model.OrderInfo, error)
}

func (service *OrderService) CreateOrder(ctx context.Context, orderNum string) (*model.OrderModel, error) {
	if !util.CheckOrderNum(orderNum) {
		return nil, errors.NewInvalidOrderNumErr(orderNum)
	}

	return service.repository.CreateOrder(ctx, orderNum)
}

func (service *OrderService) GetOrders(ctx context.Context) ([]*model.OrderInfo, error) {
	return service.repository.GetOrders(ctx)
}

func NewOrderService(repository OrderRepository) *OrderService {
	return &OrderService{
		repository: repository,
	}
}
