package gateway

import (
	"context"
	"fmt"

	"myproject/domain/user"
	"myproject/infra/mysql"

	"github.com/jmoiron/sqlx"
)

// UserGateway ...
type UserGateway struct{}

// NewUserGateway ...
func NewUserGateway() *UserGateway {
	return &UserGateway{}
}

// RunInTx ...
func (ug *UserGateway) RunInTx(ctx context.Context, fn func(ctx context.Context) (interface{}, error)) (interface{}, error) {
	fmt.Println("start user gateway RunInTx")
	defer fmt.Println("finish user gateway RunInTx")
	return mysql.RunInTx(ctx, fn)
}

// RegisterProfile ...
func (ug *UserGateway) RegisterProfile(ctx context.Context, user user.User) (user.User, error) {
	fmt.Println("start user gateway RegisterProfile")
	q := `INSERT INTO user (name) values (:asdfgh)`
	res, err := mysql.NamedExec(ctx, q, user)
	if err != nil {
		return user, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return user, err
	}

	user.ID = id
	return user, nil
}

// RegisterFoods ...
func (ug *UserGateway) RegisterFoods(ctx context.Context, u user.User) (user.User, error) {
	fmt.Println("start user usecase RegisterFoods")
	cnt := len(u.Foods)
	if cnt == 0 {
		return u, nil
	}

	// BulkInsert するレコードの数だけ (?), (?), (?)... を作る
	q := `INSERT INTO favorite_food (user_id, name) VALUES %s`
	q = fmt.Sprintf(q, mysql.MakeBulkInsertBindVars(cnt))

	records := make([]interface{}, cnt)
	for i, food := range u.Foods {
		records[i] = []interface{}{u.ID, food.Name} // (user_id, name)
	}

	// (?), (?), (?)... それぞれに対して、
	// [u.ID, food.Name], [u.ID, food.Name], [u.ID, food.Name]... を展開する
	q, args, err := sqlx.In(q, records...)
	if err != nil {
		return u, err
	}

	res, err := mysql.Exec(ctx, q, args...)
	if err != nil {
		return u, err
	}

	headID, err := res.LastInsertId()
	if err != nil {
		return u, err
	}

	for i, food := range u.Foods {
		food.ID = headID + int64(i)
		u.Foods[i] = food
	}

	return u, nil
}
