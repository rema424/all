package greeter

// MemoryDB ...
type MemoryDB struct{}

// NewMemoryDB ...
func NewMemoryDB() *MemoryDB {
	return &MemoryDB{}
}

// GetGreeter ...
func (db *MemoryDB) GetGreeter() {}
