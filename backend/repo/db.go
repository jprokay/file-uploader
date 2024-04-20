package repo

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Postgres struct {
	DB *pgxpool.Pool
}

// enums present in the database schema
var enums = []string{"directory_statuses"}

func NewPool(c context.Context, url string) (Postgres, error) {
	dbpool, err := pgxpool.New(c, url)

	if err != nil {
		return Postgres{DB: dbpool}, err
	}

	dbpool.Config().AfterConnect = func(ctx context.Context, conn *pgx.Conn) error {
		// inspects the database for defined enums and registers them to the TypeMap
		for _, v := range enums {
			pgType, err := conn.LoadType(ctx, v)
			if err != nil {
				return err
			}
			conn.TypeMap().RegisterType(pgType)
		}
		return nil
	}

	return Postgres{DB: dbpool}, nil
}

func (p *Postgres) Close() {
	p.DB.Close()
}
