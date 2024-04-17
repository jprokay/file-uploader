package repo

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
)

type baseAnon struct {
	Id int `db:"anon_id"`
}

type Anon struct {
	baseAnon
}

func (pg *Postgres) CreateAnon(ctx context.Context) error {
	query := `INSERT INTO anons DEFAULT VALUES`
	_, err := pg.DB.Exec(ctx, query)
	if err != nil {
		return fmt.Errorf("Failed to insert anon: %w", err)
	}

	return nil
}

func (pg *Postgres) GetAllAnons(ctx context.Context) ([]Anon, error) {
	query := `SELECT anon_id FROM anons`

	rows, err := pg.DB.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("Failed to get anons: %w", err)
	}
	defer rows.Close()

	return pgx.CollectRows(rows, pgx.RowToStructByName[Anon])
}
