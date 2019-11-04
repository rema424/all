package persistence

import (
	"myproject/infra/component"
)

// UserPersistence ...
type UserPersistence struct {
	db *component.DB
}

// RunInTransaction ...
func (u UserPersistence) RunInTransaction(fn func() error) error {
	// tx, err := component.GDBx.Beginx()
	// if err != nil {
	// 	return err
	// }

	// if err := fn(); err != nil {
	//   tx.
	// }
}
