package repo

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
)

type User struct {
	ID string `db:"user_id"`
}

func (pg *Postgres) CreateUser(ctx context.Context) error {
	query := `INSERT INTO users DEFAULT VALUES`
	_, err := pg.DB.Exec(ctx, query)
	if err != nil {
		return fmt.Errorf("Failed to insert user: %w", err)
	}

	return nil
}

func (pg *Postgres) GetAllUsers(ctx context.Context) ([]User, error) {
	query := `SELECT user_id FROM users`

	rows, err := pg.DB.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("Failed to get users: %w", err)
	}
	defer rows.Close()

	return pgx.CollectRows(rows, pgx.RowToStructByName[User])
}
