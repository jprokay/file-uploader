package repo

import (
	"context"
	"encoding/csv"
	"fmt"

	"github.com/jackc/pgx/v5"
)

type User struct {
	FirstName string `db:"user_first_name"`
	LastName  string `db:"user_last_name"`
	Email     string `db:"user_email"`
}

func (pg *Postgres) CreateUser(ctx context.Context, user User) error {
	query := `INSERT INTO users (user_first_name, user_last_name, user_email) VALUES (@firstName, @lastName, @email)`
	args := pgx.NamedArgs{
		"firstName": user.FirstName,
		"lastName":  user.LastName,
		"email":     user.Email,
	}

	_, err := pg.DB.Exec(ctx, query, args)
	if err != nil {
		return fmt.Errorf("Failed to insert row: %w", err)
	}

	return nil
}

func (pg *Postgres) GetAllUsers(ctx context.Context) ([]User, error) {
	query := `SELECT user_first_name, user_last_name, user_email FROM users`

	rows, err := pg.DB.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("Failed to get users: %w", err)
	}
	defer rows.Close()

	return pgx.CollectRows(rows, pgx.RowToStructByName[User])
}

func (pg *Postgres) UsersCopyFromCSV(ctx context.Context, csv *csv.Reader) (int64, error) {
	return pg.DB.CopyFrom(
		ctx,
		pgx.Identifier{"users"},
		[]string{"user_first_name", "user_last_name", "user_email"},
		pgx.CopyFromFunc(func() ([]any, error) {
			row, err := csv.Read()
			if err != nil {
				return []any{}, err
			}
			output := make([]any, len(row))

			for i, v := range row {
				output[i] = v
			}

			return output, nil
		}),
	)
}

func (pg *Postgres) UsersCopyFrom(ctx context.Context, users []User) (int64, error) {
	return pg.DB.CopyFrom(
		ctx,
		pgx.Identifier{"users"},
		[]string{"user_first_name", "user_last_name", "user_email"},
		pgx.CopyFromSlice(len(users), func(i int) ([]any, error) {
			return []any{users[i].FirstName, users[i].LastName, users[i].Email}, nil
		}),
	)
}
