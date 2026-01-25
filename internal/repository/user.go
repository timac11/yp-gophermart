package repository

import (
	"context"
	internalErrors "errors"

	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"

	"github.com/timac11/yp-gophermart/internal/errors"
	"github.com/timac11/yp-gophermart/internal/model"
)

const (
	createUserQuery = `
		INSERT INTO "user" (login, password)
		VALUES ($1, $2)
		RETURNING id, login, password'
	`
	getUserByLoginQuery = `
		SELECT id, login, password from "user"
		WHERE login=$1;
	`
	getUserQuery = `
		SELECT id, login, password from "user"
		WHERE id=$1;
	`
)

func (client *PgClient) CreateUser(ctx context.Context, user *model.UserLoginDto) (*model.User, error) {
	var userModel model.User
	err := client.pool.QueryRow(ctx, createUserQuery, user.Login, user.Password).Scan(
		&userModel.Id, &userModel.Login, &userModel.Password,
	)

	if err != nil {
		if pgErr, ok := err.(*pgconn.PgError); ok && pgerrcode.IsIntegrityConstraintViolation(pgErr.Code) {
			return nil, errors.NewEntityError(err, errors.EntityAlreadyExists, user)
		}
		return nil, err
	}

	return &userModel, nil
}

func (client *PgClient) GetUserByLogin(ctx context.Context, login string) (*model.User, error) {
	var user model.User

	err := client.pool.QueryRow(ctx, getUserByLoginQuery, login).Scan(&user.Id, &user.Login, &user.Password)

	if err != nil {
		if internalErrors.Is(err, pgx.ErrNoRows) {
			return nil, errors.NewEntityError(err, errors.EntityNotFound, login)
		}
		return nil, err
	}

	return &user, nil
}

func (client *PgClient) GetUserById(ctx context.Context, id string) (*model.User, error) {
	var user model.User

	err := client.pool.QueryRow(ctx, getUserQuery, id).Scan(&user.Id, &user.Login, &user.Password)

	if err != nil {
		return nil, err
	}

	return &user, nil
}
