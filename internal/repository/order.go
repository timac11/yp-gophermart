package repository

import (
	"context"
	"time"

	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/timac11/yp-gophermart/internal/auth"
	"github.com/timac11/yp-gophermart/internal/errors"
	"github.com/timac11/yp-gophermart/internal/model"
)

const (
	getProcessingAccruals = `
		SELECT order_num, status from "accrual"
		LEFT JOIN "order" ON "order".id = accrual.order_id
		WHERE (accrual.status = 'NEW' OR accrual.status = 'PROCESSING')
			  AND ($1 IS NULL OR accrual.created_at >= $1)
		ORDER BY accrual.created_at
		LIMIT $2;
	`
	updateAccrualStatus = `
		UPDATE accrual
		SET status = $1
		WHERE order_id = (SELECT id FROM "orders" WHERE "orders".order_num = $2);
	`
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
	payload, _ := auth.AuthPayloadFromContext(ctx)

	tx, err := client.pool.Begin(ctx)
	if err != nil {
		return nil, err
	}

	defer tx.Rollback(ctx)

	var orderModel model.OrderModel
	err = tx.
		QueryRow(ctx, createOrderQuery, orderNum, payload.UserID).
		Scan(&orderModel.ID, &orderModel.UserID, &orderModel.OrderNum, &orderModel.CreatedAt)

	if err != nil {
		if pgErr, ok := err.(*pgconn.PgError); ok && pgerrcode.IsIntegrityConstraintViolation(pgErr.Code) {
			return nil, errors.NewEntityError(err, errors.EntityAlreadyExists, orderNum)
		}
		return nil, err
	}

	_, err = tx.Exec(ctx, createAccrualQuery, orderModel.ID)
	if err != nil {
		return nil, err
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, err
	}

	return &orderModel, err
}

func (client *PgClient) GetOrders(ctx context.Context) ([]*model.OrderInfo, error) {
	payload, _ := auth.AuthPayloadFromContext(ctx)

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

func (client *PgClient) GetProcessingAccruals(ctx context.Context, date *time.Time, limit int) ([]*model.Accrual, error) {
	rows, err := client.pool.Query(ctx, getProcessingAccruals, date, limit)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	accruals, err := pgx.CollectRows(rows, pgx.RowToAddrOfStructByPos[model.Accrual])
	if err != nil {
		return nil, err
	}

	return accruals, nil
}

func (client *PgClient) UpdateAccrualStatus(ctx context.Context, accrual model.Accrual) error {
	_, err := client.pool.Exec(ctx, updateAccrualStatus, accrual.Status, accrual.Order)

	if err != nil {
		return err
	}

	return nil
}
