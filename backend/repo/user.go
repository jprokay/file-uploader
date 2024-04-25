package repo

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
)

type User struct {
	ID string `db:"user_id"`
}

type UserRepo struct {
	Postgres
}

func (pg Postgres) NewUserRepo() UserRepo {
	return UserRepo{pg}
}

func (pg *UserRepo) Create(ctx context.Context) (User, error) {
	query := `INSERT INTO users DEFAULT VALUES RETURNING *`
	rows, err := pg.DB.Query(ctx, query)
	if err != nil {
		return User{}, fmt.Errorf("Failed to insert user: %w", err)
	}

	defer rows.Close()

	return pgx.CollectExactlyOneRow(rows, pgx.RowToStructByName[User])
}

func (pg *UserRepo) GetAll(ctx context.Context) ([]User, error) {
	query := `SELECT user_id FROM users`

	rows, err := pg.DB.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("Failed to get users: %w", err)
	}
	defer rows.Close()

	return pgx.CollectRows(rows, pgx.RowToStructByName[User])
}
