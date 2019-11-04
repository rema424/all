package reservation

import "time"

// Reservation ...
type Reservation struct {
	ID       int
	OrderID  int
	Stocks   Stocks
	ExpireAt time.Time
}

// Stock ...
type Stock struct {
	ID     int
	Status StockStatus
	Item   Item
}

// Stocks ...
type Stocks []Stock

// StockStatus ...
type StockStatus int

const (
	// StockStatusOnSale ...
	StockStatusOnSale StockStatus = 1
	// StockStatusReserved ...
	StockStatusReserved StockStatus = 2
	// StockStatusSoldOut ...
	StockStatusSoldOut StockStatus = 3
)

// Order ...
type Order struct {
	ID           int
	OrderDetails OrderDetails
	TotalPrice   int
	Customer     Customer
	Employee     Employee
	CreatedAt    time.Time
}

// OrderDetail ...
type OrderDetail struct {
	OrderID       int
	Item          Item
	Quantity      int
	SubTotalPrice int
}

// OrderDetails ...
type OrderDetails []OrderDetail

// Customer ...
type Customer struct {
	ID          int
	Name        string
	PhoneNumber string
}

// Item ...
type Item struct {
	ID    int
	Name  string
	Price int
}

// Items ...
type Items []Item

// Employee ...
type Employee struct {
	ID   int
	Name string
}
