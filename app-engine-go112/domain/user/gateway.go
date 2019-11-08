package user

import (
	"context"

	"myproject/infra/mymysql"
)

// Gateway ...
type Gateway struct{}

// RunInTransaction ...
func RunInTransaction(ctx context.Context, fn func(ctx context.Context) error) error {
	return mymysql.RunInTransaction(ctx, fn)
}

// RegisterProfile ...
func RegisterProfile(ctx context.Context, user User) (User, error) {
	return user, nil
}

// RegisterIcon ...
func RegisterIcon(ctx context.Context, user User) (User, error) {
	return user, nil
}
