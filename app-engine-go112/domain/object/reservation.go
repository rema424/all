package object

import "time"

// Reservation ...
type Reservation struct {
	ID       int
	OrderID  int
	Stocks   Stocks
	ExpireAt time.Time
}
