package gateway

import (
	"context"

	"myproject/domain/user"
	"myproject/infra/mysql"
)

// UserGateway ...
type UserGateway struct{}

// NewUserGateway ...
func NewUserGateway() *UserGateway {
	return &UserGateway{}
}

// RunInTx ...
func (ug *UserGateway) RunInTx(ctx context.Context, fn func(ctx context.Context) (interface{}, error)) (interface{}, error) {
	return mysql.RunInTx(ctx, fn)
}

// RegisterProfile ...
func (ug *UserGateway) RegisterProfile(ctx context.Context, user user.User) (user.User, error) {
	return user, nil
}

// RegisterFoods ...
func (ug *UserGateway) RegisterFoods(ctx context.Context, user user.User) (user.User, error) {
	return user, nil
}
