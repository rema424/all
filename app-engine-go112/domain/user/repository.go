package user

import "context"

// Repository ...
type Repository interface {
	RunInTx(context.Context, func(context.Context) (interface{}, error)) (interface{}, error)
	RegisterProfile(ctx context.Context, user User) (User, error)
	RegisterFoods(ctx context.Context, user User) (User, error)
}
