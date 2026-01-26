package repository

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/timac11/yp-gophermart/internal/auth"
	"github.com/timac11/yp-gophermart/internal/model"
)

const (
	createOrderQuery = `
		INSERT INTO "order" (order_num, user_id)
		VALUES ($1, $2)
		RETURNING id, user_id, order_num, created_at;
	`
	getOrdersQuery = `
		SELECT order_num, status, accrual, created_at
		FROM "order"
		ORDER BY created_at DESC;
	`
)

func (client *PgClient) CreateOrder(ctx context.Context, orderNum string) (*model.OrderModel, error) {
	payload, err := auth.AuthPayloadFromContext(ctx)
	if err != nil {
		return nil, err
	}

	var orderModel model.OrderModel
	err = client.pool.QueryRow(ctx, createOrderQuery, orderNum, payload.UserID).Scan(orderModel.UserId, orderModel.OrderNum, orderModel.CreatedAt)
	if err != nil {
		return nil, err
	}

	return &orderModel, err
}

func (client *PgClient) GetOrders(ctx context.Context) ([]*model.OrderInfo, error) {
	rows, err := client.pool.Query(ctx, getOrdersQuery)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	orders, err := pgx.CollectRows(rows, pgx.RowToAddrOfStructByPos[model.OrderInfo])
	if err != nil {
		return nil, err
	}

	return orders, nil
}
