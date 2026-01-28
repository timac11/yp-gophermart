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
	getBalanceQuery = `
		SELECT value, used_value
		FROM "balance"
		WHERE user_id=$1;
	`
	getWithdrawalsListQuery = `
		SELECT order_num, value, created_at
		FROM withdrawal
		RIGHT JOIN "order"
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
	return nil, nil
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
