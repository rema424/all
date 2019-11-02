package object

// OrderDetail ...
type OrderDetail struct {
	OrderID  int
	Item     Item
	Quantity int
}

// OrderDetails ...
type OrderDetails []OrderDetail
