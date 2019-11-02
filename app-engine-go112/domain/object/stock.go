package object

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
