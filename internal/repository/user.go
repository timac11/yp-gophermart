package repository

import (
	"context"

	"github.com/timac11/yp-gophermart/internal/model"
)

func (client *PgClient) SaveUser(ctx context.Context, user model.UserDto) (*model.User, error) {
	statement := `
		INSERT INTO "user" (login, password)
		VALUES ($1, $2)
		RETURNING id, login, password
	`
	var userModel model.User
	err := client.pool.QueryRow(ctx, statement, user.Login, user.Password).Scan(
		&userModel.Id, &userModel.Login, &userModel.Password,
	)

	if err != nil {
		return nil, err
	}

	return &userModel, nil
}

func (client *PgClient) GetUserByLogin(ctx context.Context, login string) (*model.User, error) {
	statement := `
		SELECT id, login, password from "user"
		WHERE login=$1
	`

	var user model.User

	err := client.pool.QueryRow(ctx, statement, login).Scan(&user.Id, &user.Login, &user.Password)

	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (client *PgClient) GetUserById(ctx context.Context, id string) (*model.User, error) {
	statement := `
		SELECT id, login, password from "user"
		WHERE id=$1
	`

	var user model.User

	err := client.pool.QueryRow(ctx, statement, id).Scan(&user.Id, &user.Login, &user.Password)

	if err != nil {
		return nil, err
	}

	return &user, nil
}
