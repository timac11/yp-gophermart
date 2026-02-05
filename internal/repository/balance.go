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
	decreaseBalanceQuery = `
		UPDATE "balance"
		SET value = value-$1, used_value = used_value+$1
		WHERE user_id = $2
	`
	increaseBalanceQuery = `
		UPDATE "balance"
		SET value = value+$1
		WHERE user_id = $2
	`
	createWithdrawalQuery = `
		INSERT INTO "withdrawal" (order_id, value)
		VALUES ($1, $2)
		RETURNING id, order_id, CAST(value AS double precision) / 100.0;
	`
	getBalanceQuery = `
		SELECT CAST(value AS double precision) / 100.0, CAST(used_value AS double precision) / 100.0
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
	payload, _ := auth.AuthPayloadFromContext(ctx)

	var balanceModel model.BalanceInfo
	err := client.pool.QueryRow(ctx, getBalanceQuery, payload.UserID).Scan(&balanceModel.Current, &balanceModel.Withdrawn)
	if err != nil {
		return nil, err
	}

	return &balanceModel, nil
}

func (client *PgClient) CreateWithdraw(ctx context.Context, withdraw *model.CreateWithdraw) (*model.WithdrawModel, error) {
	payload, _ := auth.AuthPayloadFromContext(ctx)

	tx, err := client.pool.Begin(ctx)
	if err != nil {
		return nil, err
	}

	defer tx.Rollback(ctx)

	value := int64(withdraw.Sum * 100)

	var orderModel model.OrderModel
	err = tx.
		QueryRow(ctx, createOrderQuery, withdraw.Order, payload.UserID).
		Scan(&orderModel.ID, &orderModel.UserID, &orderModel.OrderNum, &orderModel.CreatedAt)
	if err != nil {
		return nil, err
	}

	_, err = tx.Exec(ctx, decreaseBalanceQuery, value, payload.UserID)
	if err != nil {
		return nil, err
	}

	var withdrawModel model.WithdrawModel
	err = tx.
		QueryRow(ctx, createWithdrawalQuery, orderModel.ID, value).
		Scan(&withdrawModel.ID, &withdrawModel.OrderID, &withdrawModel.Value)
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
