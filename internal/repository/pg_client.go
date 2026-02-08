package repository

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

type PgClient struct {
	pool *pgxpool.Pool
}

func NewPgClient(url string) (*PgClient, error) {
	pool, err := pgxpool.New(context.Background(), url)

	if err != nil {
		return nil, err
	}

	pgClient := PgClient{pool: pool}

	err = pgClient.applyMigration()

	if err != nil {
		return nil, err
	}

	return &pgClient, nil
}

func (client *PgClient) applyMigration() error {
	sqlDB := stdlib.OpenDBFromPool(client.pool)
	driver, err := postgres.WithInstance(sqlDB, &postgres.Config{})

	defer sqlDB.Close()

	if err != nil {
		return err
	}

	migrations, err := migrate.NewWithDatabaseInstance(
		"file://./migrations",
		"postgres",
		driver,
	)

	if err != nil {
		return err
	}

	err = migrations.Up()

	if err != nil && err != migrate.ErrNoChange {
		return err
	}

	return nil
}
