package repository

import (
	"context"

	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/timac11/yp-gophermart/internal/auth"
	"github.com/timac11/yp-gophermart/internal/errors"
	"github.com/timac11/yp-gophermart/internal/model"
)

const (
	createOrderQuery = `
		INSERT INTO "order" (order_num, user_id)
		VALUES ($1, $2)
		RETURNING id, user_id, order_num, created_at;
	`
	createAccrualQuery = `
		INSERT INTO "accrual" (order_id, status)
		VALUES ($1, 'NEW')
		RETURNING id, created_at;
	`
	getOrdersQuery = `
		SELECT order_num, status, CAST(value AS double precision) / 100.0 as value, accrual.created_at
		FROM "accrual"
		LEFT JOIN "order" ON "order".id = accrual.order_id
		WHERE user_id = $1
		ORDER BY created_at DESC;
	`
)

func (client *PgClient) CreateOrder(ctx context.Context, orderNum string) (*model.OrderModel, error) {
	payload, err := auth.AuthPayloadFromContext(ctx)
	if err != nil {
		return nil, err
	}

	tx, err := client.pool.Begin(ctx)
	if err != nil {
		return nil, err
	}

	defer tx.Rollback(ctx)

	var orderModel model.OrderModel
	err = tx.
		QueryRow(ctx, createOrderQuery, orderNum, payload.UserID).
		Scan(&orderModel.Id, &orderModel.UserId, &orderModel.OrderNum, &orderModel.CreatedAt)

	if err != nil {
		if pgErr, ok := err.(*pgconn.PgError); ok && pgerrcode.IsIntegrityConstraintViolation(pgErr.Code) {
			return nil, errors.NewEntityError(err, errors.EntityAlreadyExists, orderNum)
		}
		return nil, err
	}

	_, err = tx.Exec(ctx, createAccrualQuery, orderModel.Id)
	if err != nil {
		return nil, err
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, err
	}

	return &orderModel, err
}

func (client *PgClient) GetOrders(ctx context.Context) ([]*model.OrderInfo, error) {
	payload, err := auth.AuthPayloadFromContext(ctx)
	rows, err := client.pool.Query(ctx, getOrdersQuery, payload.UserID)
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
