package database

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Store interface {
	Querier
	CreateUserTX(context.Context, CreateUserTXParams) (CreateUserTXResult, error)
}

// SQLStore provider all functions to execute SQL queries and transaction
type SQLStore struct {
	*Queries
	connPool *pgxpool.Pool
}

func NewStore(db *pgxpool.Pool) Store {
	return &SQLStore{
		Queries:  New(db),
		connPool: db,
	}
}

func (storage *SQLStore) execTX(ctx context.Context, fn func(*Queries) error) error {
	trx, err := storage.connPool.Begin(ctx)
	if err != nil {
		return err
	}

	defer trx.Rollback(ctx)

	queries := New(trx)
	if err := fn(queries); err != nil {
		return err
	}

	return trx.Commit(ctx)
}
