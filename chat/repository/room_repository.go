package repository

import (
	"chat/infra"
	"chat/model"
)

// GetAllRooms ...
func GetAllRooms(db *infra.DB) ([]*model.Room, error) {
	q := `select id, name from room;`
	var rooms []*model.Room
	if err := db.Select(&rooms, q); err != nil {
		return nil, err
	}
	return rooms, nil
}
