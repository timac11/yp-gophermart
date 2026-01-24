package repository

import (
	"context"

	"github.com/timac11/yp-gophermart/internal/model"
)

func (client *PgClient) CreateOrder(ctx context.Context, value string) (*model.Order, error) {
	return nil, nil
}

func (client *PgClient) GetOrders(ctx context.Context) (*[]model.Order, error) {
	return nil, nil
}
