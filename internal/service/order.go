package service

import (
	"context"

	"github.com/timac11/yp-gophermart/internal/model"
)

type OrderService struct {
	repository OrderRepository
	config     *ServiceConfig
}

type OrderRepository interface {
	CreateOrder(ctx context.Context, value string) (*model.OrderModel, error)
	GetOrders(ctx context.Context) ([]*model.OrderInfo, error)
}

func (service *OrderService) CreateOrder(ctx context.Context, value string) (*model.OrderModel, error) {
	return service.repository.CreateOrder(ctx, value)
}

func (service *OrderService) GetOrders(ctx context.Context) ([]*model.OrderInfo, error) {
	return service.repository.GetOrders(ctx)
}

func NewOrderService(repository OrderRepository, config *ServiceConfig) *OrderService {
	return &OrderService{
		repository: repository,
		config:     config,
	}
}
