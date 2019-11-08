package user

import "context"

// Repository ...
type Repository interface {
	RunInTransaction(ctx context.Context, fn func(ctx context.Context) error) error
	RegisterProfile(ctx context.Context, user User) (User, error)
	RegisterIcon(ctx context.Context, user User) (User, error)
}
