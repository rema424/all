package reservation

import "context"

// Service ...
type Service struct {
	r Repository
}

// MakeReservartion ...
func (s Service) MakeReservartion(ctx context.Context, order Order) {
}
