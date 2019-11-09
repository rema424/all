package gateway

import (
	"context"

	"myproject/domain/user"
	"myproject/infra/mysql"
)

// Gateway ...
type Gateway struct{}

// RunInTx ...
func (g *Gateway) RunInTx(ctx context.Context, fn func(ctx context.Context) (interface{}, error)) (interface{}, error) {
	return mysql.RunInTx(ctx, fn)
}

// RegisterProfile ...
func (g *Gateway) RegisterProfile(ctx context.Context, user user.User) (user.User, error) {
	return user, nil
}

// RegisterFoods ...
func (g *Gateway) RegisterFoods(ctx context.Context, user user.User) (user.User, error) {
	return user, nil
}
