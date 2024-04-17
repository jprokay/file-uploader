package repo

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Postgres struct {
	DB *pgxpool.Pool
}

func NewPool(c context.Context, url string) (Postgres, error) {
	dbpool, err := pgxpool.New(c, url)

	if err != nil {
		return Postgres{DB: dbpool}, err
	}

	return Postgres{DB: dbpool}, nil
}

func (p *Postgres) Close() {
	p.DB.Close()
}
