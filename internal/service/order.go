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
	CreateOrder(ctx context.Context, value string) (*model.Order, error)
	GetOrders(ctx context.Context) (*[]model.Order, error)
}

func (service *OrderService) CreateOrder(ctx context.Context, value string) (*model.Order, error) {
	return service.CreateOrder(ctx, value)
}

func (service *OrderService) GetOrders(ctx context.Context) (*[]model.Order, error) {
	return service.GetOrders(ctx)
}

func NewOrderService(repository OrderRepository, config *ServiceConfig) *OrderService {
	return &OrderService{
		repository: repository,
		config:     config,
	}
}
