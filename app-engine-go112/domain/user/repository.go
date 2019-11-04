package user

import "context"

// Repository ...
type Repository interface {
	RunInTransaction(fn func() error) error
	Register(ctx context.Context, user User) (User, error)
}
