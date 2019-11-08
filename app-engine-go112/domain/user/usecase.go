package user

import "context"

// Interactor ...
type Interactor struct {
	Repository Repository
}

// NewInteractor ...
func NewInteractor(r Repository) *Interactor {
	return &Interactor{r}
}

// Register ...
func (i *Interactor) Register(ctx context.Context, user User) (User, error) {
	var u User

	transactionFunc := func(ctx context.Context) error {
		var err error
		u, err = i.Repository.RegisterProfile(ctx, user)
		if err != nil {
			return err
		}
		u, err = i.Repository.RegisterProfile(ctx, u)
		return err
	}

	err := i.Repository.RunInTransaction(ctx, transactionFunc)
	return u, err
}
