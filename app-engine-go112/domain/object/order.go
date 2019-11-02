package object

import "time"

// Order ...
type Order struct {
	ID           int
	OrderDetails OrderDetails
	Customer     Customer
	Employee     Employee
	CreatedAt    time.Time
}
