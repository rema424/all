package greeter

// Database provides thread-safe access to a database of greets.
type Database interface {
	GetGreeter()
}
