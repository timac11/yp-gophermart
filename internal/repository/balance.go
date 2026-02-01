package repository

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/timac11/yp-gophermart/internal/auth"
	"github.com/timac11/yp-gophermart/internal/model"
)

const (
	createBalanceQuery = `
		INSERT INTO "balance" (value, user_id)
		VALUES ($1, $2)
		RETURNING id, user_id, created_at;
	`
	updateBalanceQuery = `
		UPDATE "balance"
		SET value = value-$1
		WHERE user_id = $2
	`
	createWithdrawalQuery = `
		INSERT INTO "withdrawal" (order_id, value)
		VALUES ($1, $2)
		RETURNING id, order_id, CAST(value AS double precision) / 100.0;
	`
	getBalanceQuery = `
		SELECT value, used_value
		FROM "balance"
		WHERE user_id=$1;
	`
	getWithdrawalsListQuery = `
		SELECT order_num, CAST(value AS double precision) / 100.0, "withdrawal".created_at
		FROM withdrawal
		LEFT JOIN "order"
		ON withdrawal.order_id = "order".id
		ORDER BY created_at DESC;
	`
)

func (client *PgClient) GetBalance(ctx context.Context) (*model.BalanceInfo, error) {
	payload, err := auth.AuthPayloadFromContext(ctx)

	if err != nil {
		return nil, err
	}

	var balanceModel model.BalanceInfo
	err = client.pool.QueryRow(ctx, getBalanceQuery, payload.UserID).Scan(&balanceModel.Current, &balanceModel.Withdrawn)
	if err != nil {
		return nil, err
	}

	return &balanceModel, nil
}

func (client *PgClient) CreateWithdraw(ctx context.Context, withdraw *model.CreateWithdraw) (*model.WithdrawModel, error) {
	payload, err := auth.AuthPayloadFromContext(ctx)
	if err != nil {
		return nil, err
	}

	tx, err := client.pool.Begin(ctx)
	if err != nil {
		return nil, err
	}

	defer tx.Rollback(ctx)

	value := int64(withdraw.Sum * 100)

	var orderModel model.OrderModel
	err = tx.
		QueryRow(ctx, createOrderQuery, withdraw.Order, payload.UserID).
		Scan(&orderModel.Id, &orderModel.UserId, &orderModel.OrderNum, &orderModel.CreatedAt)
	if err != nil {
		return nil, err
	}

	_, err = tx.Exec(ctx, updateBalanceQuery, value, payload.UserID)
	if err != nil {
		return nil, err
	}

	var withdrawModel model.WithdrawModel
	err = tx.
		QueryRow(ctx, createWithdrawalQuery, withdraw.Order, value).
		Scan(&withdrawModel.Id, &withdrawModel.OrderId, &withdrawModel.Value)
	if err != nil {
		return nil, err
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, err
	}

	return &withdrawModel, err
}

func (client *PgClient) GetWithdrawals(ctx context.Context) ([]*model.WithdrawHistory, error) {
	rows, err := client.pool.Query(ctx, getWithdrawalsListQuery)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	withdrawals, err := pgx.CollectRows(rows, pgx.RowToAddrOfStructByPos[model.WithdrawHistory])
	if err != nil {
		return nil, err
	}

	return withdrawals, nil
}
