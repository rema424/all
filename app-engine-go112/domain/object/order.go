package object

import "time"

// Order ...
type Order struct {
	ID           int
	OrderDetails OrderDetails
	TotalPrice   int
	Customer     Customer
	Employee     Employee
	CreatedAt    time.Time
}
