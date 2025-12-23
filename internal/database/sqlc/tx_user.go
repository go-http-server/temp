package database

import (
	"context"
)

type CreateUserTXParams struct {
	CreateUserParams
	AfterCreate func(User) error
}

type CreateUserTXResult struct {
	User User
}

func (storage *SQLStore) CreateUserTX(ctx context.Context, args CreateUserTXParams) (CreateUserTXResult, error) {
	var result CreateUserTXResult

	err := storage.execTX(ctx, func(q *Queries) error {
		var err error
		result.User, err = q.CreateUser(ctx, args.CreateUserParams)
		if err != nil {
			return err
		}

		return args.AfterCreate(result.User)
	})

	return result, err
}
