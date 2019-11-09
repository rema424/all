package user

import "context"

// Interactor ...
type Interactor struct {
	repo Repository
}

// NewInteractor ...
func NewInteractor(r Repository) *Interactor {
	return &Interactor{r}
}

// Register ...
func (i *Interactor) Register(ctx context.Context, user User) (User, error) {
	var u User

	txFn := func(ctx context.Context) (interface{}, error) {
		var err error

		u, err = i.repo.RegisterProfile(ctx, user)
		if err != nil {
			return nil, err
		}

		u, err = i.repo.RegisterFoods(ctx, u)
		if err != nil {
			return nil, err
		}

		return u, nil
	}

	v, err := i.repo.RunInTx(ctx, txFn)
	return v.(User), err
}
