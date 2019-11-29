package greeter

// MysqlDB ...
type MysqlDB struct{}

// NewMysqlDB ...
func NewMysqlDB() *MysqlDB {
	return &MysqlDB{}
}

// GetGreeter ...
func (db *MysqlDB) GetGreeter() {}
