package reservation

import "context"

// Repository ...
type Repository interface {
	RunInTransaction() error
	Reserve(context.Context, Order) (Reservation, []OrderDetails, error)
}
