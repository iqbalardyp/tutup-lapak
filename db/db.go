package db

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Postgres struct {
	Pool *pgxpool.Pool
}

var (
	pgInstance *Postgres
)

func Connect(ctx context.Context, dbUrl string) (*Postgres, error) {
	pool, err := pgxpool.New(ctx, dbUrl)
	if err != nil {
		return nil, err
	}

	if err := pool.Ping(ctx); err != nil {
		return nil, err
	}
	pgInstance = &Postgres{pool}

	return pgInstance, err
}
